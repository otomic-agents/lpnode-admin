package utils

import (
	"admin-panel/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
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
	request := gorequest.New().Timeout(time.Duration(requestOption.Timeout) * time.Millisecond)

	err = nil
	if requestOption.JsonStruct != nil {
		jsonStr, marshalErr := json.Marshal(requestOption.JsonStruct)
		if marshalErr != nil {
			err = marshalErr
			return
		}
		requestOption.JsonStr = string(jsonStr)
	}
	superAgent := request.Post(requestOption.Url).
		Send(requestOption.JsonStr)
	for k, v := range requestOption.Header {
		logger.System.Debug("add Header", k, v)
		superAgent.Set(k, v)
	}
	resp, body, errs := superAgent.End()

	if len(errs) > 0 {
		err = errors.WithMessage(errs[0], fmt.Sprintf("请求%s发生了错误", requestOption.Url))
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("%s状态码不正确:%d", requestOption.Url, resp.StatusCode))
		return
	}

	if requestOption.TestOKFun != nil {
		ok = requestOption.TestOKFun(body)
	} else {
		ok = false
	}
	return
}
