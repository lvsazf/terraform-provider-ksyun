package ksyun

type TransformType int

const (
	TransformDefault TransformType = iota
	TransformWithN
	TransformWithFilter
	TransformListFilter
	TransformListUnique
	TransformListN
	TransformSingleN
)
