package model

type Order int8

const (
	Asc Order = iota
	Desc
)

type Pagination struct {
	Page int64
	Size int64
}
