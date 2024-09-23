package sqlite

import "url_shortener/core/model"

func OrderToStr(order model.Order) string {
	if order == model.Desc {
		return "DESC"
	}

	return "ASC"
}
