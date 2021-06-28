package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleDrive struct {
	CredentialPath, TokenPath string
	Service                   *drive.Service
}

func (g *GoogleDrive) Init() (err error) {

	if b, err := ioutil.ReadFile(g.CredentialPath); err != nil {
		return err
	} else {

		config, err := google.ConfigFromJSON(b, drive.DriveScope)
		if err != nil {
			return err
		}
		client := g.getClient(config)
		g.Service, err = drive.New(client)

		return err
	}
}

func (g *GoogleDrive) getClient(config *oauth2.Config) *http.Client {
	tok, err := g.tokenFromFile()
	if err != nil {
		tok = g.getTokenFromWeb(config)
		g.saveToken(tok)
	}
	return config.Client(context.Background(), tok)
}

func (g *GoogleDrive) tokenFromFile() (*oauth2.Token, error) {
	if f, err := os.Open(g.TokenPath); err != nil {
		return nil, err
	} else {
		defer f.Close()
		tok := &oauth2.Token{}
		return tok, json.NewDecoder(f).Decode(tok)
	}
}

func (g *GoogleDrive) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		ErrorLogger.Printf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		ErrorLogger.Printf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func (g *GoogleDrive) saveToken(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", g.TokenPath)
	f, err := os.OpenFile(g.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		ErrorLogger.Printf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
