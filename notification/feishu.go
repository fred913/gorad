package notification

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yiffyi/gorad/radhttp"
)

var feishuDefaultClient = http.Client{
	Timeout: time.Second * 10,
}

type FeishuBot struct {
	WebhookURL string
	Secret     string
}

type feishuResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func (bot *FeishuBot) genSign(timestamp int64) string {

	if bot.Secret == "" {
		return ""
	}

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, bot.Secret)

	h := hmac.New(sha256.New, []byte(stringToSign))
	h.Write([]byte{})

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (bot *FeishuBot) SendText(msg string) error {
	content := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": msg,
		},
	}

	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		content["timestamp"] = strconv.FormatInt(timestamp, 10)
		content["sign"] = bot.genSign(timestamp)
	}

	return bot.SendMessage(content)
}

func (bot *FeishuBot) SendRichText(title string, content [][]map[string]interface{}) error {
	message := map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   title,
					"content": content,
				},
			},
		},
	}

	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		message["timestamp"] = strconv.FormatInt(timestamp, 10)
		message["sign"] = bot.genSign(timestamp)
	}

	return bot.SendMessage(message)
}

func (bot *FeishuBot) SendImage(imageKey string) error {
	message := map[string]interface{}{
		"msg_type": "image",
		"content": map[string]interface{}{
			"image_key": imageKey,
		},
	}

	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		message["timestamp"] = strconv.FormatInt(timestamp, 10)
		message["sign"] = bot.genSign(timestamp)
	}

	return bot.SendMessage(message)
}

func (bot *FeishuBot) SendShareChat(chatID string) error {
	message := map[string]interface{}{
		"msg_type": "share_chat",
		"content": map[string]interface{}{
			"share_chat_id": chatID,
		},
	}

	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		message["timestamp"] = strconv.FormatInt(timestamp, 10)
		message["sign"] = bot.genSign(timestamp)
	}

	return bot.SendMessage(message)
}

func (bot *FeishuBot) SendCard(card map[string]interface{}) error {
	message := map[string]interface{}{
		"msg_type": "interactive",
		"card":     card,
	}

	if bot.Secret != "" {
		timestamp := time.Now().Unix()
		message["timestamp"] = strconv.FormatInt(timestamp, 10)
		message["sign"] = bot.genSign(timestamp)
	}

	return bot.SendMessage(message)
}

func (bot *FeishuBot) SendMessage(msg map[string]interface{}) error {
	req, err := radhttp.NewJSONPostRequest(bot.WebhookURL, msg)
	if err != nil {
		return err
	}

	var r feishuResponse
	resp, b, err := radhttp.DoAsJSON(&feishuDefaultClient, req, &r)
	if err != nil {
		return err
	}

	if err = radhttp.EnsureSuccessful(resp); err != nil {
		return fmt.Errorf("%w\nbody: %s", err, string(b))
	}

	if r.Code != 0 {
		return errors.New(r.Msg)
	}
	return nil
}

func AtUser(userID, userName string) string {
	return fmt.Sprintf("<at user_id=\"%s\">%s</at>", userID, userName)
}

func AtAll() string {
	return "<at user_id=\"all\">所有人</at>"
}

func NewRichTextAt(userID, userName string) map[string]interface{} {
	return map[string]interface{}{
		"tag":       "at",
		"user_id":   userID,
		"user_name": userName,
	}
}

func NewRichTextAtAll(displayAs *string) map[string]interface{} {
	return map[string]interface{}{
		"tag":       "at",
		"user_id":   "all",
		"user_name": *displayAs,
	}
}

func NewRichTextText(text string) map[string]interface{} {
	return map[string]interface{}{
		"tag":  "text",
		"text": text,
	}
}

func NewRichTextLink(text, href string) map[string]interface{} {
	return map[string]interface{}{
		"tag":  "a",
		"text": text,
		"href": href,
	}
}

func NewRichTextImg(imageKey string) map[string]interface{} {
	return map[string]interface{}{
		"tag":       "img",
		"image_key": imageKey,
	}
}
