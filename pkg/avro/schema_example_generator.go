package avro

import (
	"errors"
	"fmt"
)

type ExampleGenerator struct {
}

func NewExampleGenerator() *ExampleGenerator {
	return &ExampleGenerator{}
}

func (g *ExampleGenerator) Generate(schema map[string]interface{}) (map[string]interface{}, error) {
	fields, ok := schema["fields"].([]interface{})
	if !ok {
		return nil, errors.New("invalid schema structure")
	}

	fs := sliceOfInterfacesToSliceOfMaps(fields)
	v, err := readObjectFields(fs)
	if err != nil {
		return nil, fmt.Errorf("reading object fields: %w", err)
	}

	return v, nil
}

func sliceOfInterfacesToSliceOfMaps(slice []interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, val := range slice {
		v, _ := val.(map[string]interface{})
		res = append(res, v)
	}
	return res
}

// read map with: name, type, (for type="array", additional "items")
func readField(field map[string]interface{}) (fieldName string, fieldValue interface{}, err error) {
	fieldName, _ = field["name"].(string)

	// fieldType can be: "string", "array", "long", "int", "record", ["null", "type"]
	tp, _ := field["type"]

	typeName, nullable, err := parseType(tp)
	if err != nil {
		err = fmt.Errorf("parsing type: %w", err)
		return
	}

	if nullable {
		fieldValue = nil
	} else {
		switch typeName {
		case "string", "long", "int", "boolean":
			fieldValue, err = getPlainDefaultValue(typeName)
			if err != nil {
				err = fmt.Errorf("getting default value: %w", err)
				return
			}
		case "array":
			v, _ := tp.(map[string]interface{})
			fieldValue, err = readArrayType(v)
			if err != nil {
				err = fmt.Errorf("reading array type: %w", err)
				return
			}
		case "record":
			fieldValue, err = readRecordType(field)
			if err != nil {
				err = fmt.Errorf("reading record type: %w", err)
				return
			}
		default:
			err = fmt.Errorf("unsupported type %s", typeName)
		}
	}

	return
}

func parseType(tp interface{}) (fieldName string, nullable bool, err error) {
	switch val := tp.(type) {
	case string:
		fieldName = val
	case []interface{}:
		// this is union type, for example ["null", "string"]
		for _, v := range val {
			if v == "null" {
				nullable = true
			} else {
				fieldName = v.(string)
			}
		}
	case map[string]interface{}:
		fieldName, _ = val["type"].(string)
	default:
		err = errors.New(fmt.Sprintf("field has unknown type: %+v", tp))
	}
	return
}

func readArrayType(array map[string]interface{}) (value interface{}, err error) {
	items := array["items"].(map[string]interface{})
	tp, _ := items["type"].(string)
	switch tp {
	case "record":
		v, e := readRecordType(items)
		if e != nil {
			err = fmt.Errorf("reading record type: %w", err)
			return
		}
		value = []interface{}{v}
	default:
		err = errors.New(fmt.Sprintf("unsupported type %s", tp))
	}
	return
}

func readRecordType(record map[string]interface{}) (map[string]interface{}, error) {
	fields, _ := record["fields"].([]interface{})
	fs := sliceOfInterfacesToSliceOfMaps(fields)

	val, err := readObjectFields(fs)
	if err != nil {
		return nil, fmt.Errorf("reading object fields: %w", err)
	}

	return val, nil
}

func getPlainDefaultValue(tp string) (interface{}, error) {
	switch tp {
	case "string":
		return "", nil
	case "long", "int":
		return 0, nil
	case "boolean":
		return false, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported type %s", tp))
	}
}

func readObjectFields(fields []map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		name, value, err := readField(field)
		if err != nil {
			return nil, fmt.Errorf("reading field: %w", err)
		}
		res[name] = value
	}
	return res, nil
}
