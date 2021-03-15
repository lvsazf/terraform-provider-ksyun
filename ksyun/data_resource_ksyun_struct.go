package ksyun

type SdkReqTransform struct {
	Type     TransformType
	mapping  string
	mappings map[string]string
}

type SdkResponseData struct {
	Next map[string]SdkResponseData
}

type FieldRespFunc func(interface{}) interface{}

type FieldReqFunc func(interface{}, string, string, *map[string]interface{}) error

type SdkResponseMapping struct {
	Field         string
	FieldRespFunc FieldRespFunc
}

type SdkRequestMapping struct {
	Field        string
	FieldReqFunc FieldReqFunc
}
