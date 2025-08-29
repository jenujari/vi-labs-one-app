package helpers

import (
	"server/config"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"

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

func SetAuthTokenUsingBrowser() string {
	return ""
}

func GetTOTP() (string, error) {
	secret := config.GetConfig().Secret.Secret
	return totp.GenerateCode(secret, time.Now())
}
