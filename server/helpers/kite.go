package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"sync"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var (
	kiteClient   *kiteconnect.Client
	KiteInitOnce *sync.Once
)

func init() {
	KiteInitOnce = new(sync.Once)
}

func GetKiteClient() *kiteconnect.Client {

	KiteInitOnce.Do(func() {
		apiKey := config.GetConfig().Secret.ApiKey
		kiteClient = kiteconnect.New(apiKey)
	})

	if kiteClient == nil {
		GetKiteClient()
	}

	return kiteClient
}

func GetRequestTokenUsingBrowser() (string, error) {
	body := make(map[string]any)
	kc := GetKiteClient()

	body["url"] = kc.GetLoginURL()
	body["username"] = config.GetConfig().Secret.UserName
	body["password"] = config.GetConfig().Secret.Password

	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(body)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "http://headless-app:7878/zerodha-auth", buffer)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get auth token: %s", resp.Status)
	}

	var respBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	value, ok := respBody["request_token"].(string)

	if !ok {
		return "", fmt.Errorf("failed to get request token")
	}

	return value, nil
}