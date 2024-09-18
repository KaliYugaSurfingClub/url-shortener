package storage

import "errors"

var (
	ErrURLNotFound          = errors.New("url not found")
	NothingToChange         = errors.New("nothing to change")
	ErrAliasExists          = errors.New("alias already exists")
	NotEnoughTimeToGenerate = errors.New("failed to generate in the allotted time")
)
