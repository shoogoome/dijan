package paramsUtils

import (
	paramsException "dijan/exception/params"
	"encoding/json"
	"github.com/kataras/iris"
	"io/ioutil"
)

func RequestJSON(ctx iris.Context, object interface{}) {
	if err := ctx.ReadJSON(object); err != nil {
		panic(paramsException.UnmarshalBodyJsonFail())
	}
}

func RequestJsonInterface(ctx iris.Context) map[string]interface{} {

	var data map[string]interface{}

	rawData, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		panic(paramsException.UnmarshalBodyJsonFail())
	}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		panic(paramsException.UnmarshalBodyJsonFail())
	}
	return data
}

func TextToJson(text string) map[string]interface{} {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(text), &data); err != nil {
		panic(paramsException.UnmarshalTextJsonFail())
	}
	return data
}
