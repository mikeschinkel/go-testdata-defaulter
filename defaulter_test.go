package testdatadefaulter

import (
	"github.com/mikeschinkel/go-only"
	"testing"
)

var defaultChild = child{
	Key: "default",
}
var defaultInfo = info{
	Data:  "default-data",
	Child: &defaultChild,
}

type child struct {
	Key string
}
type info struct {
	Data  string
	Child *child
}

var testInfo = map[string]*info{
	"test1": {
		Data: "foo",
	},
	"test2": {
		Child: &child{Key: "bar"},
	},
}

var defaultedResult = "{\"test1\":{\"Data\":\"foo\",\"Child\":{\"Key\":\"default\"}},\"test2\":{\"Data\":\"default-data\",\"Child\":{\"Key\":\"bar\"}}}"

func TestDefaulter(t *testing.T) {
	for range only.Once {
		d := New()
		err := d.ApplyDefaults(&testInfo, defaultInfo)
		t.Run("ApplyDefaults ran without error", func(t *testing.T) {
			if err != nil {
				t.Error(err)
			}
		})
		if err != nil {
			break
		}
		t.Run("Results matched expected results", func(t *testing.T) {
			result := toJson(testInfo)
			if result != defaultedResult {
				t.Errorf("defaulted result differs from expected result:\n\tResult:   %s\n\tExpected: %s",
					result,
					defaultedResult)
				return
			}
		})
	}
}
