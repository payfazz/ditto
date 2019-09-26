package structure_test

import (
	"fmt"
	"github.com/payfazz/ditto/structure"
	"testing"
)

func TestRegisterFieldFail(t *testing.T) {
	typ := &structure.Type{}
	err := structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &structure.Type{
		Type: "test",
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}

	typ = &structure.Type{
		Value: "test",
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestRegisterGroupAndField(t *testing.T) {
	g := &structure.Group{
		Name:          "test",
		ValidInfoKeys: nil,
	}
	err := structure.RegisterGroup(g)
	if err != nil {
		t.Fatal(err)
	}

	err = structure.RegisterGroup(g)
	if err == nil {
		t.Fatal("error expected")
	}

	typ := &structure.Type{
		Type:          "test",
		Value:         "empty",
		Group:         g,
		ValidInfoKeys: nil,
	}
	err = structure.RegisterType(typ)
	if err != nil {
		t.Fatal(err)
	}

	err = structure.RegisterType(typ)
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestMerge(t *testing.T) {
	a := []structure.Info{
		{
			Key: "a",
			Child: []structure.Info{
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
			Child: []structure.Info{
				{
					Key: "c2",
				},
				{
					Key: "c1",
					Child: []structure.Info{
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

	b := []structure.Info{
		{
			Key: "A",
			Child: []structure.Info{
				{
					Key: "A1",
				},
				{
					Key: "A2",
					Child: []structure.Info{
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
			Child: []structure.Info{
				{
					Key: "c3",
				},
				{
					Key: "c2",
					Child: []structure.Info{
						{
							Key: "c21",
							Child: []structure.Info{
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

	infos := make(map[string]structure.Info)
	for _, inf := range a {
		infos[inf.Key] = inf
	}

	result := structure.MergeInfoKey(infos, b)
	print(result, 0)
}

func print(infos []structure.Info, level int) {
	for _, inf := range infos {
		for i := 0; i < level; i++ {
			fmt.Printf(" ")
		}
		fmt.Println(inf.Key)
		print(inf.Child, level + 1)
	}
}