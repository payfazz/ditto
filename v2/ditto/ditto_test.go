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
  person1: [{
	id: x,
    name: "Alice",
    welcome: "Hello " + self.name + "!",
  } for x in std.range(0,5)],
  person2: self.person1[0] { name: "Bob" },
}
`