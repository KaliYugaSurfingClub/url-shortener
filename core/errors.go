package core

import "errors"

//todo maybe it is worth returning only a custom error so that users don't get sqlite3 errors.

var ErrLinkNotFound = errors.New("original url not found")
var ErrAliasExists = errors.New("alias exists")
var ErrExpiredLink = errors.New("link is expired")
var ErrCantGenerateInTries = errors.New("can't generate in tries")
var ErrInternalError = errors.New("internal error")
