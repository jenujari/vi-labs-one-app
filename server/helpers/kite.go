package helpers

import (
	"server/config"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var (
	kiteClient *kiteconnect.Client
)

func GetKiteClient() *kiteconnect.Client {
	if kiteClient == nil {
		apiKey := config.GetConfig().Secret.ApiKey
		kiteClient = kiteconnect.New(apiKey)
	}
	return kiteClient
}
