package serializers

import (
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/types/go/common"
)

type ProtoBuffBinary struct{}

func (s *ProtoBuffBinary) Mode() common.SerializerMode {
	return common.BINARY
}

func (s *ProtoBuffBinary) Marshal(any interface{}, registry common.IRegistry) ([]byte, error) {
	obj := object.NewEncode()
	obj.Add(any)
	return obj.Data(), nil
}

func (s *ProtoBuffBinary) Unmarshal(data []byte, registry common.IRegistry) (interface{}, error) {
	obj := object.NewDecode(data, 0, registry)
	return obj.Get()
}
