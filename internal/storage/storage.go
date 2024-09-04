package storage

import "errors"

var (
	ErrURLNotFound          = errors.New("url not found")
	ErrAliasExists          = errors.New("alias already exists")
	NotEnoughTimeToGenerate = errors.New("failed to generate in the allotted time")
)
