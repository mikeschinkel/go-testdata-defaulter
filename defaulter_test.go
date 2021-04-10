package testdatadefaulter

import (
	"github.com/mikeschinkel/go-only"
	"testing"
)

var defaultData = testdata{
	Name: "default-data",
	Child: &child{
		Key: "default",
	},
	Points: 100,
}

type child struct {
	Key string
}
type testdata struct {
	Name   string
	Child  *child
	Points int
}

var testData = map[string]*testdata{
	"test1": {
		Name:   "foo",
		Points: ZeroInt,
	},
	"test2": {
		Child: &child{Key: "bar"},
	},
}

var defaultedResult = "{\"test1\":{\"Name\":\"foo\",\"Child\":{\"Key\":\"default\"},\"Points\":0},\"test2\":{\"Name\":\"default-data\",\"Child\":{\"Key\":\"bar\"},\"Points\":100}}"

func TestDefaulter(t *testing.T) {
	for range only.Once {
		d := New()
		err := d.ApplyDefaults(&testData, defaultData)
		t.Run("ApplyDefaults ran without error", func(t *testing.T) {
			if err != nil {
				t.Error(err)
			}
		})
		if err != nil {
			break
		}
		t.Run("Results matched expected results", func(t *testing.T) {
			result := toJson(testData)
			if result != defaultedResult {
				t.Errorf("defaulted result differs from expected result:\n\tResult:   %s\n\tExpected: %s",
					result,
					defaultedResult)
				return
			}
		})
		//p(testData)
	}
}
