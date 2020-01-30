package ditto_test

import (
	"github.com/payfazz/ditto/v2/ditto"
	"testing"
)

func TestJsonnetToJSON(t *testing.T) {
	d := ditto.New("test1", net)
	res, err := d.JSON(nil)
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res)
}

func TestExt(t *testing.T) {
	d := ditto.New("test-ext", netExt)

	vm := ditto.PrepareVM(map[string]interface{}{
		"prefix": "Happy Hour",
		"brunch": true,
	}, nil)

	res, err := d.JSON(vm)
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res)
}

func TestTLA(t *testing.T) {
	d := ditto.New("test2", netTla)

	vm := ditto.PrepareVM(nil, map[string]interface{}{
		"prefix": "Happy Hour",
		//"brunch": true,
	})

	res, err := d.JSON(vm)
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res)
}

var net = `
{
  person1: [{
	id: x,
    name: "Alice",
    welcome: "Hello " + self.name + "!",
  } for x in std.range(0,2)],
  person2: self.person1[0] { name: "Bob" },
}
`

var netExt = `
{
  [std.extVar('prefix') + 'Pina Colada']: {
    ingredients: [
      { kind: 'Rum', qty: 3 },
      { kind: 'Pineapple Juice', qty: 6 },
      { kind: 'Coconut Cream', qty: 2 },
      { kind: 'Ice', qty: 12 },
    ],
    garnish: 'Pineapple slice',
    served: 'Frozen',
  },

  [if std.extVar('brunch') then
    std.extVar('prefix') + 'Bloody Mary'
  ]: {
    ingredients: [
      { kind: 'Vodka', qty: 1.5 },
      { kind: 'Tomato Juice', qty: 3 },
      { kind: 'Lemon Juice', qty: 1.5 },
      { kind: 'Worcestershire', qty: 0.25 },
      { kind: 'Tobasco Sauce', qty: 0.15 },
    ],
    garnish: 'Celery salt & pepper',
    served: 'Tall',
  }
}
`

var netTla = `
function(prefix, brunch=false) {
  [prefix + 'Pina Colada']: {
    ingredients: [
      { kind: 'Rum', qty: 3 },
      { kind: 'Pineapple Juice', qty: 6 },
      { kind: 'Coconut Cream', qty: 2 },
      { kind: 'Ice', qty: 12 },
    ],
    garnish: 'Pineapple slice',
    served: 'Frozen',
  },

  [if brunch then prefix + 'Bloody Mary']: {
    ingredients: [
      { kind: 'Vodka', qty: 1.5 },
      { kind: 'Tomato Juice', qty: 3 },
      { kind: 'Lemon Juice', qty: 1.5 },
      { kind: 'Worcestershire', qty: 0.25 },
      { kind: 'Tobasco Sauce', qty: 0.15 },
    ],
    garnish: 'Celery salt & pepper',
    served: 'Tall',
  }
}
`
