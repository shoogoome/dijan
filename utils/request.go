package utils

import (
	"io"
	"net/http"
)

func Requests(method string, url string, body io.Reader) (*http.Response, error){
	client := &http.Client{}
	request, _ := http.NewRequest(method, url, body)

	request.Header.Add("token", GlobalSystemConfig.Server.SystemApiToken)
	return client.Do(request)
}
