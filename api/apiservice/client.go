package apiservice

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/etix/stoppableListener"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

// getAnalyticsClientConfig returns the oAuth client configuration for the Google Analytics API.
func getAnalyticsClientConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics",
			"https://www.googleapis.com/auth/analytics.edit",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/analytics.manage.users",
			"https://www.googleapis.com/auth/analytics.manage.users.readonly",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// receiveAuthorizationCode returns the authorization code.
func receiveAuthorizationCode(conf *oauth2.Config, listenAddress, route string) (string, error) {

	authorizationCode := make(chan string, 1)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	fmt.Println()

	// open the URL in the user's default browser
	open.Run(url)

	go func() {

		listener, err := net.Listen("tcp", listenAddress)
		if err != nil {
			log.Fatal(err)
		}

		handler := stoppableListener.Handle(listener)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if len(code) > 0 {
				fmt.Fprintf(w, "%s", code)
				authorizationCode <- code

				handler.Stop <- true
				return
			}

			fmt.Fprintf(w, "No code received")
		})

		http.Serve(handler, nil)
	}()

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	code := <-authorizationCode
	if len(code) == 0 {
		return "", fmt.Errorf("No authorization code received.")
	}

	return code, nil
}

// getAnalyticsClient returns a Google Analytics client instance.
func getAnalyticsClient(store apicredentials.TokenStorer, oAuthClientConfig *oauth2.Config, listenAddress, route string) (*http.Client, error) {

	// fetch token from store
	exchangeToken, tokenStoreError := store.GetToken()
	if tokenStoreError != nil {

		code, err := receiveAuthorizationCode(oAuthClientConfig, listenAddress, route)
		if err != nil {
			return nil, err
		}

		// request a new token
		newToken, requestTokenError := oAuthClientConfig.Exchange(oauth2.NoContext, code)
		if requestTokenError != nil {
			return nil, err
		}

		// save token to store
		store.SaveToken(newToken)

		exchangeToken = newToken
	}

	client := oAuthClientConfig.Client(oauth2.NoContext, exchangeToken)
	return client, nil
}
