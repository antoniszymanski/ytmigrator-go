// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package youtube

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/antoniszymanski/stacktrace-go"
	"github.com/cli/browser"
	"github.com/go-json-experiment/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func NewService(credentialsPath, tokenPath string) (*youtube.Service, error) {
	config, err := getConfig(credentialsPath)
	if err != nil {
		return nil, err
	}
	token, err := getToken(config, tokenPath)
	if err != nil {
		return nil, err
	}
	return getService(token)
}

func getConfig(path string) (*oauth2.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(
		data,
		youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
		youtube.YoutubepartnerScope,
	)
	if err != nil {
		return nil, err
	}
	config.RedirectURL = "http://localhost:8080"
	return config, nil
}

func getService(t *oauth2.Token) (*youtube.Service, error) {
	return youtube.NewService(
		context.Background(),
		option.WithTokenSource(oauth2.StaticTokenSource(t)),
	)
}

func getToken(config *oauth2.Config, path string) (*oauth2.Token, error) {
	t, err := getTokenFromFile(path)
	if err == nil && t.Valid() {
		return t, nil
	}

	_ = os.Remove(path)
	t, err = getTokenFromWeb(config)
	if err == nil {
		_ = saveToken(path, t)
		return t, nil
	}

	return nil, err
}

func getTokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck
	var t *oauth2.Token
	if err = json.UnmarshalRead(f, &t); err != nil {
		return nil, err
	}
	return t, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	if err := browser.OpenURL(authURL); err != nil {
		return nil, err
	}
	code, err := getCode()
	if err != nil {
		return nil, err
	}
	return config.Exchange(context.Background(), code)
}

func getCode() (string, error) {
	srv := http.Server{Addr: ":8080", ReadHeaderTimeout: 10 * time.Second}
	mux := http.NewServeMux()
	var code string
	mux.HandleFunc("/{$}",
		//nolint:errcheck
		func(w http.ResponseWriter, r *http.Request) {
			code = r.URL.Query().Get("code")
			if code == "" {
				w.Write([]byte("Error: 'code' parameter is missing from the URL"))
				return
			}
			w.Write([]byte("Success! Now you can close this page"))
			stacktrace.Go(func() {
				<-r.Context().Done()
				srv.Close()
			}, nil, nil)
		},
	)
	srv.Handler = mux

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return "", err
	}
	return code, nil
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	return json.MarshalWrite(f, token)
}
