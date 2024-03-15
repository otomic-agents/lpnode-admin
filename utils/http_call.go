package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type HttpCallRequestOption struct {
	Url        string
	Timeout    int64
	Header     map[string]string
	JsonStr    string
	TestOKFun  func(bodyStr string) bool
	JsonStruct interface{}
}
type HttpCall struct {
}

func NewHttpCall() *HttpCall {
	return &HttpCall{}
}

func (*HttpCall) PostJsonCall(requestOption *HttpCallRequestOption) (body string, ok bool, err error) {
	// request := gorequest.New().Timeout(time.Duration(requestOption.Timeout) * time.Millisecond)

	err = nil
	if requestOption.JsonStruct != nil {
		jsonStr, marshalErr := json.Marshal(requestOption.JsonStruct)
		if marshalErr != nil {
			err = marshalErr
			return
		}
		requestOption.JsonStr = string(jsonStr)
	}
	if requestOption.Header == nil {
		requestOption.Header = make(map[string]string)
	}

	if requestOption.Timeout == 0 {
		requestOption.Timeout = 2000
	}
	client := resty.New()
	client.SetTimeout(time.Duration(requestOption.Timeout) * time.Millisecond)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(requestOption.Header).
		SetBody(requestOption.JsonStr).
		Post(requestOption.Url)

	log.Println("resp:", resp.StatusCode())
	if resp.StatusCode() != 200 {
		err = errors.New(fmt.Sprintf("%s:incorrect status code :%d", requestOption.Url, resp.StatusCode()))
		return
	}
	body = string(resp.Body())
	if requestOption.TestOKFun != nil {
		ok = requestOption.TestOKFun(body)
	} else {
		ok = false
	}
	return
}
