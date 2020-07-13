package main

func GenerateAvroJSONExample(schema map[string]interface{}) map[string]interface{} {
	example := make(map[string]interface{})
	fields, ok := schema["fields"].([]interface{})
	if !ok {
		return nil
	}
	for _, field := range fields {
		f, ok := field.(map[string]interface{})
		if ok {
			k, v := readField(f)
			example[k] = v
		}
	}
	return example
}

// read map with: name, type
func readField(f map[string]interface{}) (key string, value interface{}) {
	name, _ := f["name"].(string)
	tp, _ := f["type"]

	key = name

	detectTypeValue := func(tp interface{}) interface{} {
		value = nil
		switch v := tp.(type) {
		case string:
			switch tp {
			case "string":
				value = ""
			case "long":
				value = 0
			default:
				value = nil
			}
		case []interface{}:
			value = v[0].(string)
		default:

		}
		return value
	}

	value = detectTypeValue(tp)

	return
}
