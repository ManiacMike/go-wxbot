package gowxbot

import "fmt"

// ServiceError 自定义错误
type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

// NewServiceError 自定义错误func
func NewServiceError(msg string) error {
	return &ServiceError{msg}
}

const (
	//Version 当前版本
	Version = "0.1.1"
)
