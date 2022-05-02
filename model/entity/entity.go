package entity

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"project/pkg/json"
)

const (
	StatusEnabled  = 1
	StatusDisabled = 2
)

var ErrInvalidDataType = errors.New("invalid data type")

type JsonSlice[T any] []T

func (v JsonSlice[T]) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'[', ']'}, nil
	}
	return json.Marshal(v)
}

func (v *JsonSlice[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return ErrInvalidDataType
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

type JsonMap[K comparable, V any] map[K]V

func (v JsonMap[K, V]) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'{', '}'}, nil
	}
	return json.Marshal(v)
}

func (v *JsonMap[K, V]) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return ErrInvalidDataType
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}
