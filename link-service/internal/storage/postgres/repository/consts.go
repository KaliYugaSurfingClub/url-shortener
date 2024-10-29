package repository

import "shortener/internal/core/model"

var SortLinksBy = map[model.SortBy]string{
	model.SortByCreatedAt:        " created_at ",
	model.SortLinksByCustomName:  " custom_name ",
	model.SortLinksByClicksCount: " clicks_count ",
	model.SortLinksByLastAccess:  " last_access_time ",
}

var SortClicksBy = map[model.SortBy]string{
	model.SortClickByAccessTime: " access_time ",
}
