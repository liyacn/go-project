package entity

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type JsonObject[T any] struct {
	Object T
	Valid  bool
}

func (n JsonObject[T]) Value() (driver.Value, error) {
	return n.MarshalJSON()
}

func (n *JsonObject[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return ErrInvalidDataType
	}
	return n.UnmarshalJSON(b)
}

func (n JsonObject[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Object)
}

func (n *JsonObject[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	err := dec.Decode(&n.Object)
	n.Valid = err == nil
	return err
}

func NewJsonObject[T any](v T) JsonObject[T] {
	return JsonObject[T]{Object: v, Valid: true}
}
