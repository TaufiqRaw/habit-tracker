package setmap

import "reflect"

type setMap struct {
	_map map[string]interface{}
}

// arg is optional, but arg's values must not be a pointer or any nullable value, 
// if it is a pointer use SetIfNotNil or SetIfNotNilMap method
func NewSetMap(
	arg map[string]interface{},
) setMap {
	if arg != nil {
		return setMap{arg}
	}
	return setMap{
		_map: make(map[string]interface{}),
	}
}

// arg value must be pointer.
// will be panic, if otherwise.
func (s *setMap) SetIfNotNilMap(arg map[string]interface{}) {
	for k, v := range arg {
		if v != nil {
			s._map[k] = reflect.ValueOf(v).Elem()
		}
	}
}

// pointer value must be pointer.
// will be panic, if otherwise.
func(s *setMap) SetIfNotNil(key string, pointer interface{}){
	if pointer != nil {
		s._map[key] = reflect.ValueOf(pointer).Elem()
	}
}

func (s *setMap) GetMap() map[string]interface{} {
	return s._map
}