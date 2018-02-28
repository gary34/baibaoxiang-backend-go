package protocol

import "encoding/json"

const (
	ErrCodeUnkown = 1
)

type Err struct {
	Code int
	Msg  string
}

func (err Err) Error() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func NewError(code int, msg string) *Err {
	return &Err{
		code, msg,
	}
}

func NewUnkownErr(err error) *Err {
	return NewError(ErrCodeUnkown, err.Error())
}
