package xerr

import "fmt"

type BizError struct {
	Code int    `json:"code"`
	Msg  string `jsonL:"msg"`
}

func (e BizError) Error() string {
	return fmt.Sprintf("{code:%d,msg:%s}", e.Code, e.Msg)
}

func NewCommBizErr(msg string) BizError {
	return BizError{
		Code: 1,
		Msg:  msg,
	}
}

func NewBizErr(code int, msg string) BizError {
	return BizError{
		Code: code,
		Msg:  msg,
	}
}
