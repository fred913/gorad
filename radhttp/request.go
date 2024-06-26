package radhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewJSONPostRequest(url string, body interface{}) (*http.Request, error) {
	var b, err = json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func NewURLEncodedFormRequest(url string, body url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func JSON(resp *http.Response, v interface{}) ([]byte, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return b, fmt.Errorf("unable to decode body: %w", err)
	}
	return b, nil
}

// first, check resp == nil
// second, check resp.StatusCode
// third, check err
func JSONDo(client *http.Client, req *http.Request, v interface{}) (*http.Response, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to perform request: %w", err)
	}
	defer resp.Body.Close()

	b, err := JSON(resp, v)
	return resp, b, err
}
