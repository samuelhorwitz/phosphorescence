package session

import (
	"fmt"

	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
)

func refreshIfNeeded(token *oauth2.Token) (*oauth2.Token, bool, error) {
	newToken, err := spotifyclient.GetToken(token)
	if err != nil {
		return nil, false, fmt.Errorf("Could not refresh session token: %s", err)
	}
	return newToken, newToken.AccessToken != token.AccessToken, nil
}
