package testdatadefaulter

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"runtime"
)

const EmptyString = "~"
const EmptyInt = math.MaxInt32

// TestDataDefaulter it used to apply defaults in a struct to a map of pointers to same struct type as defaults.
type TestDataDefaulter interface {
	ApplyDefaults(interface{}, interface{}) error
	SetEmptyString(string)
	SetEmptyInt(int)
}

// default defaulter type returned by New()
type testDataDefaulter struct {
	emptyString string
	emptyInt    int
}

// New returns an instantiated an object that implements the TestDataDefaulter interface
func New() TestDataDefaulter {
	return &testDataDefaulter{
		emptyString: EmptyString,
		emptyInt:    int(EmptyInt),
	}
}

// SetEmptyString sets a string value to be used to indicate
// an "empty" string for purposes of keeping a string property
// as "" instead of overwriting it by the property's default
func (tdd testDataDefaulter) SetEmptyString(empty string) {
	tdd.emptyString = empty
}

// SetEmptyInt sets an int value to be used to indicate
// an "empty" int for purposes of keeping an int property
// as 0 instead of overwriting it by the property's default
func (tdd testDataDefaulter) SetEmptyInt(empty int) {
	tdd.emptyInt = empty
}

//ApplyDefaults applies the default properties to testdata when its properties are empty.
func (tdd testDataDefaulter) ApplyDefaults(testdata interface{}, defaults interface{}) (err error) {
	var v, dv reflect.Value

	for range only.Once {
		v, err = tdd.validateTestData(testdata)
		if err != nil {
			break
		}

		dv, err = tdd.validateDefaults(defaults)
		if err != nil {
			break
		}

		var keys []interface{}
		keys, err = tdd.getMapKeys(testdata, defaults)
		if err != nil {
			break
		}

		for _, k := range keys {
			tdd.defaultTestCase(v, dv, k)
		}

	}
	if err != nil {
		err = fmt.Errorf("%s() %s",
			caller().Function,
			err.Error())
	}
	return err
}

var firstParamErrorFunc = func(kind string) error {
	return errors.New(
		fmt.Sprintf("first param must be pointer to a map of structs; got %s", kind))
}

func (tdd testDataDefaulter) getMapKeys(testdata interface{}, defaults interface{}) (keys []interface{}, err error) {
	var v reflect.Value
	for range only.Once {
		v = reflect.ValueOf(testdata)
		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
		if v.Kind() != reflect.Map {
			break
		}
		mk := v.MapKeys()
		lmk := len(mk)
		if lmk == 0 {
			err = errors.New("expects a first parameter with at least one element")
			break
		}
		f := v.MapIndex(mk[0])
		if f.Kind() != reflect.Ptr {
			ft := f.Type().String()
			fk := f.Kind().String()
			if ft == fk {
				ft = fmt.Sprintf("%s %s", ft, fk)
			}
			err = firstParamErrorFunc(fmt.Sprintf("map of '%s' instead", ft))
			break
		}
		keys = make([]interface{}, lmk)
		for ks, e := range mk {
			keys[ks] = e
		}
	}
	if keys == nil && err == nil {
		err = firstParamErrorFunc(fmt.Sprintf("'%s'",
			v.Kind().String()))
	}
	return keys, err
}

func (tdd testDataDefaulter) validateTestData(testdata interface{}) (v reflect.Value, err error) {
	for range only.Once {
		v = reflect.ValueOf(testdata)
		if v.Kind() != reflect.Ptr {
			err = firstParamErrorFunc(v.Kind().String())
			break
		}
		v = v.Elem()
		if v.Kind() != reflect.Map {
			err = firstParamErrorFunc(v.Kind().String())
			break
		}
	}
	return v, err
}

func (tdd testDataDefaulter) validateDefaults(defaults interface{}) (dv reflect.Value, err error) {
	for range only.Once {
		dv = reflect.ValueOf(defaults)
		if dv.Kind() == reflect.Ptr {
			err = errors.New("second param must not be a pointer")
			break
		}
		if dv.Kind() != reflect.Struct {
			err = errors.New(
				fmt.Sprintf("second param must be a struct; got %s",
					dv.Kind().String()))
			break
		}
	}
	return dv, err
}

func (tdd testDataDefaulter) defaultTestCase(v, d reflect.Value, mi interface{}) {
	fe := v.MapIndex(mi.(reflect.Value)).Elem()
	for fi := 0; fi < fe.NumField(); fi++ {
		tdd.defaultFieldValue(fe.Field(fi), d.Field(fi))
	}
}

// defaultFieldValue defaults to field f's value to the value in the dv parameter
func (tdd testDataDefaulter) defaultFieldValue(f, dv reflect.Value) {
	switch f.Kind() {
	case reflect.String:
		tdd.defaultFieldString(f,dv.String())
	case reflect.Int:
		tdd.defaultFieldInt(f,dv.Int())
	case reflect.Ptr:
		tdd.defaultFieldPtr(f,dv)
	}
}

// defaultFieldPtr defaults to field f's pointer value to the value in the dv parameter
func (tdd testDataDefaulter) defaultFieldPtr(f, dp reflect.Value) {
	if f.IsNil() {
		f.Set(dp)
	}
}

// defaultFieldString defaults to field f's string value to the value in the dv parameter
func (tdd testDataDefaulter) defaultFieldString(f reflect.Value, ds string) {
	fs := f.String()
	for range only.Once {
		if fs == tdd.emptyString {
			ds = ""
		} else if fs != "" {
			break
		}
		f.SetString(ds)
	}
}

// defaultFieldInt defaults to field f's int value to the value in the dv parameter
func (tdd testDataDefaulter) defaultFieldInt(f reflect.Value, di int64) {
	fi := int(f.Int())
	for range only.Once {
		if fi == tdd.emptyInt {
			di = 0
		} else if fi != 0 {
			break
		}
		f.SetInt(int64(di))
	}
}

// caller returns a Frame with the File/Line/Function that called it.
func caller() runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}
