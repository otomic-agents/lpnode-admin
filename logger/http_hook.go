package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type Msg struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Call    bool   `json:"call"`
	At      string `json:"at"`
}

type hooker struct {
}

func NewHooker() (*hooker, error) {

	return &hooker{}, nil
}

func (h *hooker) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields)
	data["Level"] = entry.Level.String()
	if data["Level"] != "error" {
		return nil
	}
	data["Time"] = entry.Time
	data["Message"] = entry.Message

	for k, v := range entry.Data {
		if errData, isError := v.(error); logrus.ErrorKey == k && v != nil && isError {
			data[k] = errData.Error()
		} else {
			data[k] = v
		}
	}
	strData := fmt.Sprint(data)
	go h.ReportMsgToHTTP(strData) // 异步发送出去，把错误
	return nil
}
func (h *hooker) ReportMsgToHTTP(msg string) {
	client := &http.Client{Timeout: 5 * time.Second}
	//生成要访问的url
	//提交请求
	message := &Msg{"2RkVROtJOe1", msg, false, ""}
	jsonData, _ := json.Marshal(message)
	request, err := http.NewRequest("POST", "http://47.240.72.34:7001/sendTemplateMsg", bytes.NewBuffer(jsonData))
	//增加header选项
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", "this is a token")

	if err != nil {
		fmt.Errorf("发送产生了错误[%s]", err)
	}
	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		log.Println("report 发生了错误", err)
		return
	}

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	return
}

func (h *hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}
