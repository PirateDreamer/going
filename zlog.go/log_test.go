package zlog

import (
	"context"
	"fmt"
	"testing"
)

func TestZlog(t *testing.T) {
	InitZlog()
	c := context.WithValue(context.Background(), "req_id", "123")
	err := doSomething()
	Log(c).Error(err.Error())
}

func doSomething() error {
	return fmt.Errorf("something went wrong")
}
