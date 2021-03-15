package util

import (
	simpleJson "github.com/bitly/go-simplejson"
	"github.com/valyala/fasthttp"
	"regexp"
)

func SendRequest(method string, url string, params *simpleJson.Json, header map[string]string) *simpleJson.Json {
	byteData, _ := params.Encode()
	request := fasthttp.AcquireRequest()
	request.SetBody(byteData)
	request.Header.SetMethod(method)
	request.Header.SetContentType("application/json")

	for key, val := range(header) {
		request.Header.Add(key, val)
	}

	request.SetRequestURI(url)

	response := fasthttp.AcquireResponse()

	if err := fasthttp.Do(request, response); err != nil {
		panic("Error while sending request!")
	}

	fasthttp.ReleaseRequest(request)

	defer fasthttp.ReleaseResponse(response)

	jsonData, _ := simpleJson.NewJson(response.Body())
	return jsonData
}

func Standardlize(input string) string {
	doubleQuoteRegex := regexp.MustCompile("[“”]")
	endingRegex := regexp.MustCompile("\\r\\n")
	output := doubleQuoteRegex.ReplaceAllString(input, "\"")
	output = endingRegex.ReplaceAllString(output, "\n")
	return output
}
