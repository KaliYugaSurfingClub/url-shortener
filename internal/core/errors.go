package core

import "errors"

var ErrLinkNotFound = errors.New("url not found")
var ErrAliasExists = errors.New("alias name already exists")
var ErrCustomNameExists = errors.New("custom name already exists")
var ErrCantGenerateInTries = errors.New("can't generate in tries")
var ErrLinkDoesNotBelongsUser = errors.New("link does not belong to user")
var ErrTryToCompleteUnexactingClick = errors.New("try to complete existing click")
