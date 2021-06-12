package util

import (
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/txfs19260817/scopelens/server/config"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type reCAPTCHAResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname,omitempty"`
	ErrorCodes  []string  `json:"error-codes,omitempty"`
}

func ReCaptcha(token string) error {
	// send request to google
	secret := config.Jwt.ReCaptchaSecret
	response, err := http.Get("https://recaptcha.net/recaptcha/api/siteverify?secret=" + secret + "&response=" + token)
	if err != nil {
		return err
	}

	// get response
	defer response.Body.Close()
	var data reCAPTCHAResponse
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return err
	}

	// validate response
	if !data.Success {
		return fmt.Errorf("reCAPTCHA verification failed. ")
	}
	return nil
}
