package errcode

import (
	"reflect"
	"strconv"
)

type ErrCode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e ErrCode) String() string { return strconv.Itoa(e.Code) + " : " + e.Msg }

func (e ErrCode) Response() (int, ErrCode) {
	return e.Code / 100, e
}

func (e ErrCode) WithMsg(msg string) (int, ErrCode) {
	return e.Code / 100, ErrCode{
		e.Code,
		msg,
	}
}

func FromError(err error) (int, ErrCode) {
	s := reflect.TypeOf(err).String()
	ec := CommonServerError
	switch s {
	case "validator.ValidationErrors":
		ec = ParamsInvalid
	case "proto.RedisError":
		ec = RedisServerError
	case "nsq.ErrProtocol":
		ec = NsqServerError
	case "*url.Error":
		ec = ResponseTimeout
	case "*json.SyntaxError", "*json.UnmarshalTypeError":
		ec = ResponseWrong
	}
	return ec.Response()
}
