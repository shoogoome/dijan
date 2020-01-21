package paramsUtils

import (
	paramsException "dijan/exception/params"
	reflectUtils "dijan/utils/reflect"
	"reflect"
	"strings"
)

// 接口定义
type ParamsParser interface {
	Str(key string, desc string, defaultValue ...interface{}) string
	Int(key string, desc string, defaultValue ...interface{}) int
	Float(key string, desc string, defaultValue ...interface{}) float64
	Time(key string, desc string, defaultValue ...interface{}) int64
	Has(key string) bool
	Bool(key string, desc string, defaultValue ...interface{}) bool
	List(key string, desc string, defaultValue ...interface{}) []interface{}
	Map(key string, desc string, defaultValue ...interface{}) map[string]interface{}
	getRow(key string, transformation ...bool) reflect.Value
	Diff(diffObject interface{})
	DisDiff()
}

// 接口实体定义
type params struct {
	object map[string]interface{}
	diffObject interface{}
	isDiff bool
}

// 创建params类
func NewParamsParser(object map[string]interface{}) ParamsParser {
	return &params{
		object: object,
		isDiff: false,
	}
}

// 获取原始数据
func (p *params) getRow(key string, transformation ...bool) reflect.Value {
	// 获取值
	var v reflect.Value
	v = reflect.ValueOf(p.diffObject)
	if len(transformation) > 0 && transformation[0] {
		return v.FieldByName(FTuoFeng(key))
	} else {
		return v.FieldByName(key)
	}
}

// 判断数据是否存在
func (p *params) Has(key string) bool {
	v, ok := p.object[key]

	if !ok {
		return false
	}
	switch k := v.(type) {
	case string:
		return len(k) > 0
	}
	return true
}

// 获取int类型数据
func (p *params) Int(key string, desc string, defaultValue ...interface{}) int {

	value, ok := p.object[key]; if ok {
		if valueInt, ok := value.(float64); ok == true {
			return int(valueInt)
		} else {
			panic(paramsException.ParamsIsNotStandard(key, "int"))
		}
	}
	if p.isDiff {
		rowValue := p.getRow(key, true)
		// 判断类型是否正确
		if reflectUtils.IsExist(rowValue) && strings.Contains(rowValue.Type().String(), "int") {
			return int(rowValue.Int())
		} else {
			return 0
		}
	}

	if len(defaultValue) > 0{
		return defaultValue[0].(int)
	}
	panic(paramsException.LackParams(desc))
}

// 获取float类型数据
func (p *params) Float(key string, desc string, defaultValue ...interface{}) float64 {
	value, ok := p.object[key]; if ok {
		if valueFloat, ok := value.(float64); ok == true {
			return valueFloat
		} else {
			panic(paramsException.ParamsIsNotStandard(key, "float"))
		}
	}
	if p.isDiff {
		rowValue := p.getRow(key, true)
		// 判断类型是否正确
		if reflectUtils.IsExist(rowValue) && strings.Contains(rowValue.Type().String(), "float") {
			return rowValue.Float()
		} else {
			return 0
		}
	}

	if len(defaultValue) > 0{
		return defaultValue[0].(float64)
	}
	panic(paramsException.LackParams(desc))
}

// 获取bool类型数据
func (p *params) Bool(key string, desc string, defaultValue ...interface{}) bool {
	value, ok := p.object[key]; if ok {
		if valueBool, ok := value.(bool); ok == true {
			return valueBool
		} else {
			panic(paramsException.ParamsIsNotStandard(key, "bool"))
		}
	}
	if p.isDiff {
		rowValue := p.getRow(key, true)
		// 判断类型是否正确
		if reflectUtils.IsExist(rowValue) && strings.Contains(rowValue.Type().String(), "bool") {
			return rowValue.Bool()
		} else {
			return false
		}
	}

	if len(defaultValue) > 0{
		return defaultValue[0].(bool)
	}
	panic(paramsException.LackParams(desc))
}

// 获取string类型值
func (p *params) Str(key string, desc string, defaultValue ...interface{}) string {

	value, ok := p.object[key]; if ok {
		if valueString, ok := value.(string); ok == true {
			return valueString
		} else {
			panic(paramsException.ParamsIsNotStandard(key, "string"))
		}
	}
	if p.isDiff {
		rowValue := p.getRow(key, true)
		// 判断类型是否正确
		if reflectUtils.IsExist(rowValue) && strings.Contains(rowValue.Type().String(), "string") {
			return rowValue.String()
		} else {
			return ""
		}
	}

	if len(defaultValue) > 0{
		return defaultValue[0].(string)
	}
	panic(paramsException.LackParams(desc))
}

// 获取string类型值
func (p *params) Time(key string, desc string, defaultValue ...interface{}) int64 {
	return int64(p.Int(key, desc, defaultValue...))
}

// 获取list类型
func (p *params) List(key string, desc string, defaultValue ...interface{}) []interface{} {
	value, ok := p.object[key]; if ok {
		return value.([]interface{})
	}
	if len(defaultValue) > 0 {
		return defaultValue[0].([]interface{})
	}
	panic(paramsException.LackParams(desc))
}

// 获取list类型
func (p *params) Map(key string, desc string, defaultValue ...interface{}) map[string]interface{} {
	value, ok := p.object[key]; if ok {
		return value.(map[string]interface{})
	}
	if len(defaultValue) > 0 {
		return defaultValue[0].(map[string]interface{})
	}
	panic(paramsException.LackParams(desc))
}

func (p *params) Diff(diffObject interface{}) {
	p.diffObject = diffObject
	p.isDiff = true
}

func (p *params) DisDiff() {
	p.diffObject = nil
	p.isDiff = false
}


