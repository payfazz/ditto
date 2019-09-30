package ditto_test

import (
	"fmt"
	"github.com/payfazz/ditto"
	"testing"
)

func TestRegisterFieldFail(t *testing.T) {
	typ := &ditto.Type{}
	err := ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &ditto.Type{
		Type: "test",
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &ditto.Type{
		Value: "test",
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestRegisterGroupAndField(t *testing.T) {
	g := &ditto.Group{
		Name:          "test",
		ValidInfoKeys: nil,
	}
	err := ditto.RegisterGroup(g)
	if err != nil {
		t.Fatal(err)
	}

	err = ditto.RegisterGroup(g)
	if err == nil {
		t.Fatal("error expected")
	}

	typ := &ditto.Type{
		Type:          "test",
		Value:         "empty",
		Group:         g,
		ValidInfoKeys: nil,
	}
	err = ditto.RegisterType(typ)
	if err != nil {
		t.Fatal(err)
	}

	err = ditto.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestMerge(t *testing.T) {
	a := []ditto.Info{
		{
			Key: "a",
			Child: []ditto.Info{
				{
					Key: "a1",
				},
				{
					Key: "a2",
				},
			},
		},
		{
			Key: "b",
		},
		{
			Key: "c",
			Child: []ditto.Info{
				{
					Key: "c2",
				},
				{
					Key: "c1",
					Child: []ditto.Info{
						{
							Key: "c11",
						},
						{
							Key: "c12",
						},
					},
				},
			},
		},
	}

	b := []ditto.Info{
		{
			Key: "A",
			Child: []ditto.Info{
				{
					Key: "A1",
				},
				{
					Key: "A2",
					Child: []ditto.Info{
						{
							Key: "A21",
						},
					},
				},
			},
		},
		{
			Key:        "b",
			IsOptional: true,
		},
		{
			Key: "c",
			Child: []ditto.Info{
				{
					Key: "c3",
				},
				{
					Key: "c2",
					Child: []ditto.Info{
						{
							Key: "c21",
							Child: []ditto.Info{
								{
									Key: "c211",
								},
							},
						},
						{
							Key: "c12",
						},
					},
				},
			},
		},
	}

	infos := make(map[string]ditto.Info)
	for _, inf := range a {
		infos[inf.Key] = inf
	}

	result := ditto.MergeInfoKey(infos, b)
	print(result, 0)
}

func print(infos []ditto.Info, level int) {
	for _, inf := range infos {
		for i := 0; i < level; i++ {
			fmt.Printf(" ")
		}
		fmt.Println(inf.Key)
		print(inf.Child, level+1)
	}
}
