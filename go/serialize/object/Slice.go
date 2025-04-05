package object

import (
	"github.com/saichler/types/go/common"
	"reflect"
)

type Slice struct{}

func (this *Slice) add(any interface{}, data *[]byte, location *int) error {
	if any == nil {
		sizeObjectType.add(int32(-1), data, location)
		return nil
	}
	slice := reflect.ValueOf(any)
	if slice.Len() == 0 {
		sizeObjectType.add(int32(-1), data, location)
		return nil
	}

	sizeObjectType.add(int32(slice.Len()), data, location)
	dataByte, ok := any.([]byte)
	if ok {
		(*data)[*location] = 1
		*location++
		copy((*data)[*location:*location+len(dataByte)], dataByte)
		*location += len(dataByte)
	} else {
		(*data)[*location] = 0
		*location++
		obj := NewDecode(data, location, nil)
		for i := 0; i < slice.Len(); i++ {
			element := slice.Index(i).Interface()
			obj.Add(element)
		}
	}
	return nil
}

func (this *Slice) get(data *[]byte, location *int, registry common.IRegistry) (interface{}, error) {
	l := sizeObjectType.get(data, location)
	size := int(l.(int32))
	if size == -1 || size == 0 {
		return nil, nil
	}

	if (*data)[*location] == 1 {
		*location++
		result := make([]byte, size)
		copy(result, (*data)[*location:*location+size])
		*location += size
		return result, nil
	}

	elems := make([]interface{}, 0)
	var sliceType reflect.Type

	obj := NewDecode(data, location, registry)

	for i := 0; i < size; i++ {
		element, _ := obj.Get()
		if i == 0 {
			sliceType = reflect.SliceOf(reflect.ValueOf(element).Type())
		}
		elems = append(elems, element)
	}

	newSlice := reflect.MakeSlice(sliceType, len(elems), len(elems))
	for i := 0; i < int(size); i++ {
		if elems[i] != nil {
			newSlice.Index(i).Set(reflect.ValueOf(elems[i]))
		}
	}

	return newSlice.Interface(), nil
}
