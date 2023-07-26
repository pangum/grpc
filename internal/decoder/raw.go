package decoder

import (
	"fmt"
	"io"
	"reflect"

	"github.com/goexl/exc"
	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Raw struct {
	*runtime.JSONPb

	bytesType reflect.Type
}

func NewRaw(pb *runtime.JSONPb) *Raw {
	return &Raw{
		JSONPb: pb,

		bytesType: reflect.TypeOf([]byte(nil)),
	}
}

func (r *Raw) NewDecoder(reader io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(to any) (err error) {
		value := reflect.ValueOf(to)
		if data, re := io.ReadAll(reader); nil != re {
			err = re
		} else if value.Kind() != reflect.Ptr {
			err = exc.NewField("必须是指针类型", field.New("field", fmt.Sprintf("%T", to)))
		} else if value = value.Elem(); value.Type() != r.bytesType {
			err = exc.NewField("必须是二进制数组", field.New("field", fmt.Sprintf("%T", to)))
		} else {
			value.Set(reflect.ValueOf(data))
		}

		return
	})
}
