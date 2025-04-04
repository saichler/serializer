package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Struct struct {
}

func (this *Struct) add(any interface{}) ([]byte, int, error) {
	if any == nil {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4, nil
	}

	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			sizeBytes, _ := sizeObjectType.add(int32(-1))
			return sizeBytes, 4, nil
		}
		val = val.Elem()
	}

	typ := val.Type()
	typeName := typ.Name()
	var pbData []byte

	if typeName == "Transaction" {
		pbData, _ = TransactionSerializer.Marshal(any, nil)
	} else {
		pb := any.(proto.Message)
		pbd, err := proto.Marshal(pb)
		if err != nil {
			return []byte{}, 0, errors.New("Failed To marshal proto " + typeName + " in protobuf object:" + err.Error())
		}
		pbData = pbd
	}

	obj := NewEncode()
	obj.appendBytes(stringObjectType.add(typeName))
	obj.appendBytes(sizeObjectType.add(int32(len(pbData))))
	obj.appendBytes(pbData, len(pbData))

	return obj.Data(), obj.Location(), nil
}

func (this *Struct) get(data []byte, location int, registry common.IRegistry) (interface{}, int, error) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	if size == -1 || size == 0 {
		return nil, 4, nil
	}

	typeN, typeSize := stringObjectType.get(data, location)
	typeName := typeN.(string)
	var info common.IInfo
	var err error
	var pb interface{}
	isTransaction := typeName == "Transaction"
	if !isTransaction {
		info, err = registry.Info(typeName)
		if err != nil {
			return []byte{}, 0, errors.New("Unknown proto name " + typeName + " in registry, please register it.")
		}

		pb, err = info.NewInstance()
		if err != nil {
			return []byte{}, 0, errors.New("Unknown proto name " + typeName + " in registry, please register it.")
		}
	}

	location += typeSize
	s, _ := sizeObjectType.get(data, location)
	size = s.(int32)
	location += 4
	protoData := data[location : location+int(size)]

	if isTransaction {
		pb, _ = TransactionSerializer.Unmarshal(protoData, nil)
	} else {
		err = proto.Unmarshal(protoData, pb.(proto.Message))
		if err != nil {
			return []byte{}, 0, errors.New("Failed To unmarshal proto " + typeName + ":" + err.Error())
		}
	}

	return pb, typeSize + 4 + int(size), nil
}
