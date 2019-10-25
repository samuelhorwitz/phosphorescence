package mail

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	mailgunAPIKey string
	mailgunClient *http.Client
)

type Config struct {
	MailgunAPIKey string
}

func Initialize(cfg *Config) {
	mailgunAPIKey = cfg.MailgunAPIKey
	mailgunClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func SendAuthentication(name string, email string, magicLink string, utcOffsetMinutes int64) error {
	form := url.Values{}
	form.Add("from", `"Samuel @ Phosphorescence" <noreply@phosphor.me>`)
	form.Add("to", fmt.Sprintf(`"%s" <%s>`, strings.Replace(name, `"`, "", -1), strings.Replace(email, ">", "", -1)))
	form.Add("subject", fmt.Sprintf("Phosphorescence Log In Magic Link (Attempt at %s)", time.Now().UTC().Add(time.Duration(utcOffsetMinutes)*time.Minute).Format("3:04 PM on Monday, January _2")))
	form.Add("template", "magiclinklogin")
	form.Add("h:X-Mailgun-Variables", fmt.Sprintf(`{"token": "%s"}`, magicLink))
	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/phosphor.me/messages", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Could not build Mailgun profile request: %s", err)
	}
	req.SetBasicAuth("api", mailgunAPIKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := mailgunClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Mailgun send mail request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Mailgun request responded with %d", res.StatusCode)
	}
	return nil
}
