package ditto

func extractArrayMap(child []interface{}) []map[string]interface{} {
	var arrMap []map[string]interface{}
	var tmpMap map[string]interface{}
	for _, v := range child {
		tmpMap = v.(map[string]interface{})
		arrMap = append(arrMap, tmpMap)
	}
	return arrMap
}
