package ditto

import "github.com/google/go-jsonnet"

func JsonnetToJSON(net string) (string, error) {
	return jsonnet.MakeVM().EvaluateSnippet("file", net)
}