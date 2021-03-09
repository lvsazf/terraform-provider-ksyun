package ksyun

type SdkResponseData struct {
	Next map[string]SdkResponseData
}

type FieldRespFunc func(interface{}) interface{}

type FieldReqFunc func(interface{}, string, *map[string]interface{}) error

type SdkResponseMapping struct {
	Field         string
	FieldRespFunc FieldRespFunc
}

type SdkRequestMapping struct {
	Field        string
	FieldReqFunc FieldReqFunc
}
