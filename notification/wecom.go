package notification

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/yiffyi/gorad/radhttp"
)

var weComDefaultClient = http.Client{
	Timeout: time.Second * 10,
}

type WeComBot struct {
	Key string
}

type weComResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (bot *WeComBot) SendText(msg string) error {
	return bot.SendMessage(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": msg,
		},
	})
}

func (bot *WeComBot) SendMarkdown(md string) error {
	return bot.SendMessage(map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": md,
		},
	})
}

func (bot *WeComBot) SendMessage(msg map[string]interface{}) error {
	req, err := radhttp.NewJSONPostRequest("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+bot.Key, msg)
	if err != nil {
		return err
	}

	var r weComResponse
	resp, b, err := radhttp.JSONDo(&weComDefaultClient, req, &r)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP Status: %s\n\t%s", resp.Status, string(b))
	}

	if r.ErrCode != 0 {
		return errors.New(r.ErrMsg)
	}
	return nil
}
