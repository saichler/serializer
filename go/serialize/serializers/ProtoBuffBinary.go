package serializers

import (
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

type ProtoBuffBinary struct{}

func (s *ProtoBuffBinary) Mode() ifs.SerializerMode {
	return ifs.BINARY
}

func (s *ProtoBuffBinary) Marshal(any interface{}, resources ifs.IResources) ([]byte, error) {
	obj := object.NewEncode()
	obj.Add(any)
	return obj.Data(), nil
}

func (s *ProtoBuffBinary) Unmarshal(data []byte, resources ifs.IResources) (interface{}, error) {
	obj := object.NewDecode(data, 0, resources.Registry())
	return obj.Get()
}
