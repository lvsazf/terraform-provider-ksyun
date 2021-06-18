package ksyun

type SdkReqTransform struct {
	forceUpdateParam bool
	Type             TransformType
	mapping          string
	mappings         map[string]string
	FieldReqFunc     FieldReqFunc
}

type SdkResponseData struct {
	Next map[string]SdkResponseData
}

type FieldRespFunc func(interface{}) interface{}

type FieldReqFunc func(interface{}, string, map[string]string, int, string, *map[string]interface{}) (int, error)

type FieldReqSingleFunc func(interface{}, string, string, *map[string]interface{}) error

type SdkResponseMapping struct {
	Field         string
	FieldRespFunc FieldRespFunc
	KeepAuto      bool
}

type SdkRequestMapping struct {
	Field        string
	FieldReqFunc FieldReqSingleFunc
}
