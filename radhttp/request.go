package radhttp

import (
	"bytes"
	"encoding/json"
	"errors"
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

func NewGetRequest(baseURL string, query url.Values) (*http.Request, error) {
	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		panic(err)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	return req, nil
}

func AsJSON(resp *http.Response, v interface{}) ([]byte, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return b, fmt.Errorf("unable to decode body as JSON: %w", err)
	}
	return b, nil
}

// first, check if resp == nil
// here does not check statusCode
func DoAsJSON(client *http.Client, req *http.Request, v interface{}) (*http.Response, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to perform request: %w", err)
	}
	defer resp.Body.Close()

	b, err := AsJSON(resp, v)
	return resp, b, err
}

func IsSuccessful(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func EnsureSuccessful(resp *http.Response) error {
	if resp == nil {
		return errors.New("response is nil")
	}
	if !IsSuccessful(resp) {
		return fmt.Errorf("bad statusCode: %s", resp.Status)
	}
	return nil
}
