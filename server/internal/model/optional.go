package model

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Optional[T any] struct {
	Set   bool
	Null  bool
	Value T
}

func (o *Optional[T]) Unmarshal(b []byte) error {
	o.Set = true
	if string(b) == "null" {
		o.Null = true
		return nil
	}
	return json.Unmarshal(b, &o.Value)
}

func optionalValuer[T any](v reflect.Value) any {
	o := v.Interface().(Optional[T])
	if !o.Set || o.Null {
		return nil // nil → omitempty пропустит проверку
	}
	return o.Value
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterCustomTypeFunc(optionalValuer[string], Optional[string]{})
		v.RegisterCustomTypeFunc(optionalValuer[int], Optional[int]{})
	}
}
