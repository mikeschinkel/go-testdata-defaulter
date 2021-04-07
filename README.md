# Test Data Defaulter

Simple Go package for defaulting tabular test data in a map of struct pointers
to structs to the values in a default struct, e.g. `map[string]*yourTestDataStruct` 
where `yourTestDataStruct` is whatever struct you need to create for your tests.

Here is an example struct `info` and it's related structure `child`:

```go
type testdata struct {
    Name  string
    Child *child
}

type child struct {
    Key string
}
```
Now using that `testdata` struct, let's create some actual test data:

```go
var testData = map[string]*testdata{
	"test1": {
		Name: "foo",
	},
	"test2": {
		Child: &child{Key: "bar"},
	},
}
```

And let's define our defaults, itself a single instance of `testdata`:

```go
var defaultData = testdata{
	Name: "default-data",
	Child: &child{
		Key: "default",
	},
}
```

To use this package you simply create a new instance and then run `ApplyDefaults`:

```go
defaulter := testdatadefaulter.New()
err := defaulter.ApplyDefaults(&testData, defaultData)
if err != nil {
    panic(err)
}
```
After running the above code you will see `testData` having the following values, show in JSON format:

```json
{
   "test1": {
      "Name": "foo",
      "Child": {
         "Key": "default"
      }
   },
   "test2": {
      "Name": "default-data",
      "Child": {
         "Key": "bar"
      }
   }
}
```

## Explicitly setting "empty" values.
Sometimes you want a default value, but you want to explictly set your test data to an 'empty' value. As we know Go does not support empty values except for the data types that support `nil`, so it is not possible to tell if `string`s and `int`s are explicitly set to their empty value, or just initialized to their empty value.

To address this logical conundrum `testdatadefaulter` defines `EmptyString` and `ZeroInt` constants you can use instead of `""` and `0` respectively in your test data, assuming the values we chose for these contstants do not clash with your test data needs.

Imagine our `testdata` struct from above also had a `Points` property of type `int` that we wanted to default to `100`, but we wanted `test2` to be zero after the defaulting. Here are the changes to the code above to see the merged value for `testdata["test1"].Points==100` and for `testdata["test2"].Points==0`:

```go
var defaultData = testdata{
    Name: "default-data",
    Child: &child{
        Key: "default",
    },
    Points: 100,
}

type testdata struct {
    Name  string
    Child *child
    Points int
}

var testData = map[string]*testdata{
    "test1": {
        Name: "foo"
    },
    "test2": {
        Child: &child{Key: "bar"},
        Points: ZeroInt,
    },
}
```

### SetEmptyString() and SetZeroInt()

If the values you chose do clash with your test data needs then `testdatadefaulter` defines `SetEmptyString()` and `SetZeroInt()` methods to allow you to assign your own sentinel values to use instead. You might use it like so:

```go
const myEmptyString = "~~~"
var testData = map[string]*testdata{
    "test3": {
        Name: myEmptyString,
    },
}
func main() {
    defaulter := testdatadefaulter.New()
    defaulter.SetEmptyString(myEmptyString)
    err := defaulter.ApplyDefaults(&testData, defaultData)
    if err != nil {
        panic(err)
    }
}
```

### String and Int not sufficient?
Need more than `string` or `int` values to be defaulted to an empty value?  Submit a pull request or even just an issue asking for the enhancement and it's like I will be able to add it quickly.

