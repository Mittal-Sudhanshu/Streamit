package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"streamit/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ConfigGoogle to set config of oauth
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		}, // you can use other scopes to get more data
		Endpoint: google.Endpoint,
	}
	return conf
}

// GetEmail of user
func GetUserDetails(token string) models.GoogleResponse {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		panic(err)
	}
	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken},
		},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data models.GoogleResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	log.Print(data)
	return data
}
