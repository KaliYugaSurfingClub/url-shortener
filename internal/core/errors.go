package core

import "errors"

var ErrLinkNotFound = errors.New("original url not found")
var ErrAliasExists = errors.New("alias already exists")
var ErrCustomNameExists = errors.New("custom name already exists")
var ErrCantGenerateInTries = errors.New("can't generate in tries")
var ErrInternal = errors.New("internal error")
