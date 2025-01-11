package domain

import (
	"log"
	"reflect"
)

type columnDataContainer[T interface{}] struct {
	Str T
	AsArray []string
}

// all field in T must hold string value
func NewColumnDataContainer[T interface{}](names T) columnDataContainer[T] {
	return columnDataContainer[T]{
		Str: names,
		AsArray: getNameArr(names),
	}
}

func getNameArr[T interface{}](colNames T)[]string {
	r := []string{}
	t := reflect.TypeFor[T]()
	for i := 0; i < t.NumField(); i++ {
		fname := t.Field(i).Name
		name := reflect.Indirect(reflect.ValueOf(colNames)).FieldByName(fname)
		if name.Kind() != reflect.String {
			log.Fatalf("NewColumnDataContainer() :: All Field In Column Data Container's struct must contain string")
		}
		r = append(r, name.String())
	}
	return r
}
