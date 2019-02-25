package collections

import (
	"github.com/ONSdigital/log.go/log"
)

type Error struct {
	Data        log.Data
	Message     string
	OriginalErr error
}

type MovePlan struct {
	Collection    *Collection
	MovingFromAbs string
	MovingFromRel string
	MovingToRel   string
	MasterDir     string
}

func (e Error) Error() string {
	return e.Message
}

func NewErr(message string, err error, data log.Data) Error {
	return Error{
		Message:     message,
		Data:        data,
		OriginalErr: err,
	}
}
