package paramsUtils

import "reflect"

func TuoFeng(name string) string {
	index := 0
	newName := make([]byte, len(name) * 2 - 1)
	for i := 0; i < len(name); i++ {
		if name[i] >= 'A' && name[i] <= 'Z' {
			if i == 0 {
				newName[index] = name[i] + 32
				index += 1
			} else {
				newName[index] = '_'
				newName[index + 1] = name[i] + 32
				index += 2
			}
		} else {
			newName[index] = name[i]
			index += 1
		}
	}
	return string(newName[:index])
}

func FTuoFeng(name string) string {
	index := 0
	newName := make([]byte, len(name))
	for i := 0; i < len(name); i++ {
		if name[i] >= 'a' && name[i] <= 'z' {
			if i == 0 {
				newName[index] = name[i] - 32
				index += 1
			} else{
				newName[index] = name[i]
				index += 1
			}
		} else {
			newName[index] = name[i + 1] - 32
			i += 1
			index += 1
		}
	}
	return string(newName[:index])
}

// 结构体转字典
func StructToDict(obj interface{}) map[string]interface{} {

	m := make(map[string]interface{})
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[TuoFeng(relType.Field(i).Name)] = elem.Field(i).Interface()
	}
	return m
}