package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"reflect"
)

type Object struct {
	data     []byte
	location int
	typeName string
	registry common.IRegistry
}

type Primitive interface {
	add(interface{}) ([]byte, int)
	get([]byte, int) (interface{}, int)
}

type Complex interface {
	add(interface{}) ([]byte, int, error)
	get([]byte, int, string, common.IRegistry) (interface{}, int, error)
}

var primitives = make(map[reflect.Kind]Primitive)
var complex = make(map[reflect.Kind]Complex)

var sizeObjectType = &Int32{}
var stringObjectType = &String{}

func init() {
	primitives[reflect.Int] = &Int{}
	primitives[reflect.Uint32] = &UInt32{}
	primitives[reflect.Uint64] = &UInt64{}
	primitives[reflect.Int32] = &Int32{}
	primitives[reflect.Int64] = &Int64{}
	primitives[reflect.Float32] = &Float32{}
	primitives[reflect.Float64] = &Float64{}
	primitives[reflect.String] = &String{}
	primitives[reflect.Bool] = &Bool{}

	complex[reflect.Ptr] = &Struct{}
	complex[reflect.Slice] = &Slice{}
	complex[reflect.Map] = &Map{}
}

func NewEncode(data []byte, location int) *Object {
	return NewDecode(data, location, "", nil)
}

func NewDecode(data []byte, location int, typeName string, registry common.IRegistry) *Object {
	obj := &Object{}
	obj.data = data
	obj.location = location
	obj.registry = registry
	obj.typeName = typeName
	return obj
}

func (obj *Object) Data() []byte {
	return obj.data
}

func (obj *Object) Location() int {
	return obj.location
}

func (obj *Object) Add(any interface{}) error {
	kind := reflect.ValueOf(any).Kind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return errors.New("Did not find any Object for kind " + kind.String())
	}

	obj.addKind(kind)
	var b []byte
	var l int
	var e error

	if pOK {
		b, l = p.add(any)
	} else {
		b, l, e = c.add(any)
	}
	obj.location += l
	obj.data = append(obj.data, b...)
	return e
}

func (obj *Object) Get() (interface{}, error) {
	kind := obj.getKind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return nil, errors.New("Did not find any Object for kind " + kind.String())
	}

	var d interface{}
	var l int
	var e error

	if pOK {
		d, l = p.get(obj.data, obj.location)
	} else {
		d, l, e = c.get(obj.data, obj.location, obj.typeName, obj.registry)
	}

	obj.location += l
	return d, e
}

func (obj *Object) addKind(kind reflect.Kind) {
	b, l := sizeObjectType.add(int32(kind))
	obj.location += l
	obj.data = append(obj.data, b...)
}

func (obj *Object) getKind() reflect.Kind {
	i, l := sizeObjectType.get(obj.data, obj.location)
	obj.location += l
	return reflect.Kind(i.(int32))
}
