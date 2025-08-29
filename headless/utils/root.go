// Package utils provides utility functions for the application
package utils

import (
	"headless/config"
	"time"

	"github.com/pquerna/otp/totp"
)

func GetTOTP() (string, error) {
	secret := config.GetConfig().Secret.Secret
	return totp.GenerateCode(secret, time.Now())
}
