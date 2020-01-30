package ditto_test

import (
	"github.com/payfazz/ditto/v2/ditto"
	"testing"
)

func TestJsonnetToJSON(t *testing.T) {
	res, err := ditto.JsonnetToJSON(json)
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res)
}

var json = `
{
  person1: {
    name: "Alice",
    welcome: "Hello " + self.name + "!",
  },
  person2: self.person1 { name: "Bob" },
}
`