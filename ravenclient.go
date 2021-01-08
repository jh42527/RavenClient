package ravenclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// EmailRequest struct
type EmailRequest struct {
	Recipient                string
	Subject                  string
	Body                     string
	From                     string
	DisableDefaultTemplate   bool
	DisableMessageEncryption bool
}

// SmsRequest struct
type SmsRequest struct {
	Phone string
	Text  string
}

// RavenResponse struct
type RavenResponse struct {
	MessageID string `json:"MessageId"`
}

// RavenClient struct
type RavenClient struct {
	URL   string
	Token string
}

// NewClient function
func NewClient(url string, token string) *RavenClient {
	return &RavenClient{
		URL:   url,
		Token: token,
	}
}

// SendMail Method
func (client *RavenClient) SendMail(recipient string, from string, subject string, message string) (string, error) {
	email := &EmailRequest{
		Recipient:                recipient,
		Subject:                  subject,
		Body:                     message,
		From:                     from,
		DisableDefaultTemplate:   true,
		DisableMessageEncryption: true,
	}

	msgJSON, err := json.Marshal(email)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", client.URL+"messaging/email", bytes.NewBuffer(msgJSON))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+client.Token)
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case
		200,
		201:
		var ravenResponse RavenResponse
		err := json.Unmarshal(body, &ravenResponse)
		if err != nil {
			return "", err
		}
		return ravenResponse.MessageID, nil
	}

	return "", fmt.Errorf("%s", body)
}

// SendSMS Method
func (client *RavenClient) SendSMS(phone string, text string) (string, error) {
	sms := &SmsRequest{
		Phone: phone,
		Text:  text,
	}

	msgJSON, err := json.Marshal(sms)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", client.URL+"messaging/sms", bytes.NewBuffer(msgJSON))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+client.Token)
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case
		200,
		201:
		var ravenResponse RavenResponse
		err := json.Unmarshal(body, &ravenResponse)
		if err != nil {
			return "", err
		}
		return ravenResponse.MessageID, nil
	}

	return "", fmt.Errorf("%s", body)
}
