package model

type SortBy int8

const (
	SortByCreatedAt SortBy = iota

	SortLinksByCustomName
	SortLinksByClicksCount
	SortLinksByLastAccess

	SortClickByAccessTime
)

type Order int8

const (
	Asc Order = iota
	Desc
)

type Pagination struct {
	Page int64
	Size int64
}

type Sort struct {
	By    SortBy
	Order Order
}
