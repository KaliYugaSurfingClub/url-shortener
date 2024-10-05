package request

//func ParsePagination(query url.Values, defaultPageSize int64) model.Pagination {
//	var ok bool
//
//	params := model.Pagination{}
//
//	if params.Page, ok = positiveIntFromUrl(query, "page", 1); !ok {
//		return params, fmt.Errorf("invalid query type") //todo validation error
//	}
//
//	if params.Size, ok = positiveIntFromUrl(query, "size", defaultPageSize); !ok {
//		return params, fmt.Errorf("invalid query type")
//	}
//
//	return params, nil
//}
//
//func positiveIntFromUrl(query url.Values, varName string, defaultValue int64) (int64, bool) {
//	str := query.Get(varName)
//	if str == "" {
//		return defaultValue, true
//	}
//
//	num, err := strconv.Atoi(str)
//	if err != nil {
//		return 0, false
//	}
//
//	if num <= 0 {
//		return 0, false
//	}
//
//	return int64(num), true
//}
