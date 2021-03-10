package ksyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func SchemaSetToInstanceMap(s interface{}, prefix string, input *map[string]interface{}) {
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)[prefix+"."+strconv.Itoa(count)] = v
	}
}

func SchemaSetToFilterMap(s interface{}, prefix string, index int, input *map[string]interface{}) {
	(*input)["Filter."+strconv.Itoa(index)+".Name"] = prefix
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(count)] = v
	}
}

func SchemaSetsToFilterMap(d *schema.ResourceData, filters []string, req *map[string]interface{}) *map[string]interface{} {
	index := 0
	for _, v := range filters {
		var idsString []string
		if ids, ok := d.GetOk(v); ok {
			idsString = SchemaSetToStringSlice(ids)
		}
		if len(idsString) > 0 {
			index++
			(*req)[fmt.Sprintf("Filter.%v.Name", index)] = strings.Replace(v, "_", "-", -1)
		}
		for k1, v1 := range idsString {
			if v1 == "" {
				continue
			}
			(*req)[fmt.Sprintf("Filter.%v.Value.%d", index, k1+1)] = v1
		}
	}
	return req
}
func hashStringArray(arr []string) string {
	var buf bytes.Buffer

	for _, s := range arr {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func writeToFile(filePath string, data interface{}) error {
	absPath, err := getAbsPath(filePath)
	if err != nil {
		return err
	}
	os.Remove(absPath)
	var bs []byte
	switch data := data.(type) {
	case string:
		bs = []byte(data)
	default:
		bs, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}

	return ioutil.WriteFile(absPath, bs, 0422)
}

func getAbsPath(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("get current user got an error: %#v", err)
		}

		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}
	return filepath.Abs(filePath)
}

func merageResultDirect(result *[]map[string]interface{}, source []interface{}) {
	for _, v := range source {
		*result = append(*result, v.(map[string]interface{}))
	}
}

// schemaSetToStringSlice used for converting terraform schema set to a string slice
func SchemaSetToStringSlice(s interface{}) []string {
	vL := []string{}

	for _, v := range s.(*schema.Set).List() {
		vL = append(vL, v.(string))
	}

	return vL
}

func getSdkValue(keyPattern string, obj interface{}) (interface{}, error) {
	keys := strings.Split(keyPattern, ".")
	root := obj
	for index, k := range keys {
		if reflect.ValueOf(root).Kind() == reflect.Map {
			root = root.(map[string]interface{})[k]
			if root == nil {
				return root, nil
			}

		} else if reflect.ValueOf(root).Kind() == reflect.Slice {
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("keyPattern %s index %d must number", keyPattern, index)
			}
			if len(root.([]interface{})) < i {
				return nil, nil
			}
			root = root.([]interface{})[i]
		}
	}
	return root, nil
}

type SliceMappingFunc func(map[string]interface{}) map[string]interface{}

type IdMappingFunc func(string, map[string]interface{}) string

type SdkSliceData struct {
	IdField          string
	IdMappingFunc    IdMappingFunc
	SliceMappingFunc SliceMappingFunc
	TargetName       string
}

func sliceMapping(ids []string, data []map[string]interface{}, sdkSliceData SdkSliceData, item interface{}) ([]string, []map[string]interface{}) {
	if mm, ok := item.(map[string]interface{}); ok {
		if sdkSliceData.IdMappingFunc != nil && sdkSliceData.IdField != "" {
			ids = append(ids, sdkSliceData.IdMappingFunc(sdkSliceData.IdField, mm))
		}
		if sdkSliceData.SliceMappingFunc != nil {
			data = append(data, sdkSliceData.SliceMappingFunc(mm))
		}
	}
	return ids, data
}

func mapMapping(sdkSliceData SdkSliceData, item interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	if mm, ok := item.(map[string]interface{}); ok {
		if sdkSliceData.SliceMappingFunc != nil {
			data = sdkSliceData.SliceMappingFunc(mm)
		}
	}
	return data
}

func getSchemeElem(resource *schema.Resource, keys []string) *schema.Resource {
	r := resource
	if r == nil {
		return nil
	}
	for _, v := range keys {
		if elem, o := r.Schema[v].Elem.(*schema.Resource); o {
			r = elem
		}
	}
	return r
}

func SdkRequestAutoMapping(d *schema.ResourceData, resource *schema.Resource, isUpdate bool, only []string, extraMapping map[string]SdkRequestMapping) (map[string]interface{}, error) {
	var req map[string]interface{}
	var err error
	req = make(map[string]interface{})
	if only != nil {
		for _, k := range only {
			if isUpdate {
				err = requestUpdateMapping(d, k, extraMapping, &req)
			} else {
				err = requestCreateMapping(d, k, extraMapping, &req)
			}
		}
	}
	for k, _ := range resource.Schema {
		if isUpdate {
			err = requestUpdateMapping(d, k, extraMapping, &req)
		} else {
			err = requestCreateMapping(d, k, extraMapping, &req)
		}
	}

	return req, err
}

func requestCreateMapping(d *schema.ResourceData, k string, extraMapping map[string]SdkRequestMapping, req *map[string]interface{}) error {
	var err error
	if v, ok := d.GetOk(k); ok {
		if _, ok := extraMapping[k]; !ok {
			(*req)[Downline2Hump(k)] = v
		} else {
			m := extraMapping[k]
			if m.FieldReqFunc == nil {
				(*req)[m.Field] = v
			} else {
				err = m.FieldReqFunc(v, m.Field, req)
			}
		}
	}
	return err
}

func requestUpdateMapping(d *schema.ResourceData, k string, extraMapping map[string]SdkRequestMapping, req *map[string]interface{}) error {
	var err error
	if d.HasChange(k) && !d.IsNewResource() {
		err = requestCreateMapping(d, k, extraMapping, req)
	}
	return err
}

func SdkResponseAutoResourceData(d *schema.ResourceData, resource *schema.Resource, item interface{}, extra map[string]SdkResponseMapping, start ...bool) interface{} {
	setFlag := false
	if start == nil || (len(start) > 0 && start[0]) {
		setFlag = true
	}
	if reflect.ValueOf(item).Kind() == reflect.Map {
		result := make(map[string]interface{})
		root := item.(map[string]interface{})
		for k, v := range root {
			var value interface{}
			var err error
			m := SdkResponseMapping{}
			target := Hump2Downline(k)
			if _, ok := extra[k]; ok {
				m = extra[k]
				target = m.Field
			}
			if r, ok := resource.Schema[target]; ok {
				if r.Elem != nil {
					if elem, ok := r.Elem.(*schema.Resource); ok {
						if m.FieldRespFunc != nil {
							value = m.FieldRespFunc(v)
						} else {
							value = SdkResponseAutoResourceData(d, elem, v, extra, false)
						}
					} else if _, ok := r.Elem.(*schema.Schema); ok {
						value = v
					}
				} else {
					value = v
				}
			} else {
				continue
			}
			if setFlag {
				err = d.Set(target, value)
				if err != nil {
					log.Println(err.Error())
					panic("ERROR: " + err.Error())
				}

			} else {
				result[target] = value
			}
		}
		if len(result) > 0 {
			return result
		}
	} else if reflect.ValueOf(item).Kind() == reflect.Slice {
		var result []interface{}
		result = []interface{}{}
		root := item.([]interface{})
		for _, v := range root {
			value := SdkResponseAutoResourceData(d, resource, v, extra, false)
			result = append(result, value)
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func SdkResponseAutoMapping(resource *schema.Resource, collectField string, item map[string]interface{}, computeItem map[string]interface{},
	extra map[string]interface{}, extraMapping map[string]SdkResponseMapping) map[string]interface{} {
	var result map[string]interface{}
	result = make(map[string]interface{})
	keys := strings.Split(collectField, ".")
	if len(keys) == 0 {
		return result
	}

	if computeItem != nil {
		for k, v := range computeItem {
			item[k] = v
		}
	}

	if _, ok := resource.Schema[keys[0]]; ok {
		elem := getSchemeElem(resource, keys)
		for k, v := range item {
			needExtraMapping := false
			target := Hump2Downline(k)
			m := SdkResponseMapping{}
			if extraMapping != nil {
				if _, ok := extraMapping[k]; ok {
					m = extraMapping[k]
					target = m.Field
					needExtraMapping = true
				}
			}
			if _, ok := elem.Schema[target]; !ok {
				continue
			}
			needDefaultMapping := false
			if extra == nil {
				needDefaultMapping = true
			} else {
				if _, ok := extra[target]; !ok {
					needDefaultMapping = true
				}
			}
			if needDefaultMapping {
				if needExtraMapping {
					if m.FieldRespFunc == nil {
						result[m.Field] = v
					} else {
						result[m.Field] = m.FieldRespFunc(v)
					}
				} else {
					result[target] = v
				}
			} else {
				result[target] = extra[target]
			}
		}
	}
	return result
}

func SdkSliceMapping(d *schema.ResourceData, result interface{}, sdkSliceData SdkSliceData) ([]string, []map[string]interface{}, error) {
	var err error
	var ids []string
	ids = []string{}
	var data []map[string]interface{}
	data = []map[string]interface{}{}

	if reflect.TypeOf(result).Kind() == reflect.Slice {
		var length = 0
		if v, ok := result.([]map[string]interface{}); ok {
			length = len(v)
			for _, v1 := range v {
				ids, data = sliceMapping(ids, data, sdkSliceData, v1)
			}
		} else {
			root := result.([]interface{})
			length = len(root)
			for _, v2 := range root {
				ids, data = sliceMapping(ids, data, sdkSliceData, v2)
			}
		}

		if d != nil && sdkSliceData.TargetName != "" {
			d.SetId(hashStringArray(ids))
			err = d.Set("total_count", length)
			if err != nil {
				return nil, nil, err
			}
			err = d.Set(sdkSliceData.TargetName, data)
			if err != nil {
				return nil, nil, err
			}
			if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
				err = writeToFile(outputFile.(string), data)
				if err != nil {
					return nil, nil, err
				}
			}
		}

	} else if reflect.TypeOf(result).Kind() == reflect.Map {
		if v, ok := result.(map[string]interface{}); ok {
			data = append(data, mapMapping(sdkSliceData, v))
		}
	}
	return ids, data, nil
}

func GetSdkParam(d *schema.ResourceData, params []string) map[string]interface{} {
	sdkParam := make(map[string]interface{})
	for _, v := range params {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			sdkParam[vv] = fmt.Sprintf("%v", v1)
		}
	}
	return sdkParam
}
