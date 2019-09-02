package misc

import (
	"fmt"
	"reflect"
)

type AsData struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64

	Float32 float32
	Float64 float64

	Bool bool

	String string
}

func As(data interface{}) (*AsData, error) {
	switch data.(type) {
	case int:
		v := reflect.ValueOf(data).Int()
		return &AsData{
			Int: int(v),
		}, nil

	case int8:
		v := reflect.ValueOf(data).Int()
		return &AsData{
			Int8: int8(v),
		}, nil

	case int16:
		v := reflect.ValueOf(data).Int()
		return &AsData{
			Int16: int16(v),
		}, nil

	case int32:
		v := reflect.ValueOf(data).Int()
		return &AsData{
			Int32: int32(v),
		}, nil

	case int64:
		v := reflect.ValueOf(data).Int()
		return &AsData{
			Int64: int64(v),
		}, nil

	case uint:
		v := reflect.ValueOf(data).Uint()
		return &AsData{
			Uint: uint(v),
		}, nil

	case uint8:
		v := reflect.ValueOf(data).Uint()
		return &AsData{
			Uint8: uint8(v),
		}, nil

	case uint16:
		v := reflect.ValueOf(data).Uint()
		return &AsData{
			Uint16: uint16(v),
		}, nil

	case uint32:
		v := reflect.ValueOf(data).Uint()
		return &AsData{
			Uint32: uint32(v),
		}, nil

	case uint64:
		v := reflect.ValueOf(data).Uint()
		return &AsData{
			Uint64: uint64(v),
		}, nil

	case bool:
		v := reflect.ValueOf(data).Bool()
		return &AsData{
			Bool: bool(v),
		}, nil

	case float32:
		v := reflect.ValueOf(data).Float()
		return &AsData{
			Float32: float32(v),
		}, nil

	case float64:
		v := reflect.ValueOf(data).Float()
		return &AsData{
			Float64: float64(v),
		}, nil

	case string:
		v := reflect.ValueOf(data).String()
		return &AsData{
			String: string(v),
		}, nil

	}
	return nil, fmt.Errorf("Type not supported")
}
