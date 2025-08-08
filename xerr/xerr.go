package xerr

import "fmt"

type BizError struct {
	Code string
	Msg  string
}

func NewBizErr(code string, msg string) *BizError {
	return &BizError{code, msg}
}

func (e *BizError) Error() string {
	return fmt.Sprintf(`{"code":"%s","msg":"%s"}`, e.Code, e.Msg)
}

func NewCommBizErr(msg string) *BizError {
	return NewBizErr("1", msg)
}
