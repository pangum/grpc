package decoder

import (
	"fmt"
	"io"
	"reflect"

	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Raw struct {
	*runtime.JSONPb

	bytesType reflect.Type
}

func NewRaw() *Raw {
	return &Raw{
		JSONPb: new(runtime.JSONPb),

		bytesType: reflect.TypeOf([]byte(nil)),
	}
}

func (r *Raw) NewDecoder(reader io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(to any) (err error) {
		value := reflect.ValueOf(to)
		if data, re := io.ReadAll(reader); nil != re {
			err = re
		} else if value.Kind() != reflect.Ptr {
			err = exception.New().Message("必须是指针类型").Field(field.New("field", fmt.Sprintf("%T", to))).Build()
		} else if value = value.Elem(); value.Type() != r.bytesType {
			err = exception.New().Message("必须是二进制数组").Field(field.New("field", fmt.Sprintf("%T", to))).Build()
		} else {
			value.Set(reflect.ValueOf(data))
		}

		return
	})
}
