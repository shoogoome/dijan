package ui

import (
	uiException "dijan/exception/ui"
	"github.com/kataras/iris"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var contentTypeMap = map[string]string {
	"html": "text/html",
	"css": "text/css",
	"js": "application/javascript",
}

// 静态资源获取
func StaticResource(ctx iris.Context, fileType string, fileName string) {

	_, name := path.Split(fileName)
	nameList := strings.Split(name, ".")
	suffix := nameList[len(nameList) - 1]

	file, err := os.Open(path.Join("/static", fileType, fileName))
	if err != nil {
		panic(uiException.UIFileNoFoundException)
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		panic(uiException.UIFileNoFoundException)
	}

	ctx.ResponseWriter().Header().Set("Via", "dijan")
	ctx.ResponseWriter().Header().Set("Content-type", contentTypeMap[suffix])
	ctx.ResponseWriter().Header().Set("Content-Disposition", "inline")
	ctx.ResponseWriter().Write(fileByte)
	defer func() {
		if err := recover(); err != nil {
			ctx.Text("找无该资源")
		}
	}()
}
