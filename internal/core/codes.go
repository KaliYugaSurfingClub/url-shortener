package core

import (
	"shortener/errs"
)

var AliasExistsCode = errs.Code("alias already exists")
var CustomNameExistsCode = errs.Code("custom name already exists")
