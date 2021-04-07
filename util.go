package testdatadefaulter

import (
	"encoding/json"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"reflect"
	"strings"
)

var counter = 0

func p(i interface{}) {
	counter++
	v := reflect.ValueOf(i)
	if v2, k := i.(reflect.Value); k {
		v = v2
	}
	t := v.Kind().String()
	k := v.Kind().String()
	var nl string
	if t == k {
		k = ""
		nl = "\n\n"
	} else {
		k = fmt.Sprintf("%sKind:  %s\n",
			indent,
			t)
	}

	fmt.Printf("%d.) Type:  %s\n%s%sValue: %#v%s",
		counter,
		t,
		k,
		indent,
		i,
		nl)
	x := byType(i)
	if x != nil {
		fmt.Printf("%s%s\n\n",
			indent,
			toJson(x))
	}
}

func byType(i interface{}) (x interface{}) {
	for range only.Once {

		if v, k := i.([]reflect.Value); k {
			var is []interface{}
			for _, s := range v {
				is = append(is, s.Interface())
			}
			x = is
			break
		}

		if v, k := i.(reflect.Value); k {
			i = v.Interface()
			if v.Kind() == reflect.String {
				break
			}
			x = i
			break
		}
		x = byType(reflect.ValueOf(i))

	}
	return x
}

var indent = strings.Repeat(" ", 4)

type marshaler func(interface{}) ([]byte, error)

func toIndentedJson(i interface{}) string {
	return _toJson(i, func(i interface{}) ([]byte, error) {
		return json.MarshalIndent(i, indent, "   ")
	})
}

func toJson(i interface{}) string {
	return _toJson(i, func(i interface{}) ([]byte, error) {
		return json.Marshal(i)
	})
}

func _toJson(i interface{}, fn marshaler) string {
	j, err := fn(i)
	if err != nil {
		println(err.Error())
	}
	return string(j)
}
