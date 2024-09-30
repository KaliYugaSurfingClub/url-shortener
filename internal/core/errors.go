package core

type LogicErr struct {
	msg string
}

func (e LogicErr) Error() string {
	return e.msg
}

var ErrLinkNotFound = LogicErr{msg: "original url not found"}
var ErrAliasExists = LogicErr{msg: "alias name already exists"}
var ErrCustomNameExists = LogicErr{msg: "custom name already exists"}
var ErrCantGenerateInTries = LogicErr{msg: "can't generate in tries"}
