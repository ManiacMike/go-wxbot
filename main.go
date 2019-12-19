package main

import (
	"flag"
	"fmt"
	// "log"
	// "net/http"
)

//ServiceError 自定义错误
type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

//NewServiceError 自定义错误func
func NewServiceError(msg string) error {
	return &ServiceError{msg}
}

const (
	//Version 当前版本
	Version = "0.1.1"
)

var debug = flag.String("d", "off", "if on debug mode")

func main() {

	flag.Parse()

	fmt.Printf("debug mode %s\n", *debug)

	wx := wxweb{}
	wx.start()
}
