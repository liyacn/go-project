package errcode

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type ErrCode struct {
	int
	string
}

type export struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e ErrCode) Code() int { return e.int }

func (e ErrCode) String() string { return strconv.Itoa(e.int) + " : " + e.string }

func (e ErrCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(export{
		Code: e.int,
		Msg:  e.string,
	})
}

func (e ErrCode) Response() (int, ErrCode) {
	return e.int / 100, e
}

func (e ErrCode) WithMsg(msg string) (int, ErrCode) {
	return e.int / 100, ErrCode{
		e.int,
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
