package model

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("error not found")
	ErrBadRequest         = errors.New("error bad request")
	ErrServiceUnavailable = errors.New("error service unavailable")
)

type ServiceError struct {
	Err error
	Key string
}

func (e *ServiceError) Error() string {
	return e.Err.Error()
}

const (
	UndefinedError = "error"

	// validation
	AgentGroupNameIsTaken = "agent_group_name_taken"
)

func NewError(err error) *ServiceError {
	return &ServiceError{
		Err: err,
		Key: UndefinedError,
	}
}

func ErrorAgentGroupNameIsTaken() *ServiceError {
	return &ServiceError{
		Err: fmt.Errorf("agent group name is taken: %w", ErrBadRequest),
		Key: AgentGroupNameIsTaken,
	}
}
