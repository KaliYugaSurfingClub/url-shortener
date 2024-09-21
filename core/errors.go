package core

import "errors"

//todo maybe it is worth returning only a custom error so that users don't get sqlite3 errors.

var ErrLinkNotFound = errors.New("original url not found")
var ErrLastInsertId = errors.New("getting last insert Id failed")
var ErrAliasExists = errors.New("alias exists")
