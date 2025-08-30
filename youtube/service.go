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
	"github.com/dsnet/try"
	"github.com/go-json-experiment/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func NewService(credentialsPath, tokenPath string) (_ *youtube.Service, err error) {
	defer try.Handle(&err)

	config := try.E1(getConfig(credentialsPath))
	token := try.E1(getToken(config, tokenPath))
	service := try.E1(getService(token))
	return service, nil
}

func getConfig(path string) (config *oauth2.Config, err error) {
	defer try.Handle(&err)

	data := try.E1(os.ReadFile(path))
	config = try.E1(google.ConfigFromJSON(
		data,
		youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
		youtube.YoutubepartnerScope,
	))
	config.RedirectURL = "http://localhost:8080"
	return
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

func getTokenFromFile(path string) (t *oauth2.Token, err error) {
	defer try.Handle(&err)

	f := try.E1(os.Open(path))
	defer f.Close() //nolint:errcheck
	try.E(json.UnmarshalRead(f, &t))
	return
}

func getTokenFromWeb(config *oauth2.Config) (t *oauth2.Token, err error) {
	defer try.Handle(&err)

	authURL := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	try.E(browser.OpenURL(authURL))
	code := try.E1(getCode())
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

func saveToken(path string, token *oauth2.Token) (err error) {
	defer try.Handle(&err)

	f := try.E1(os.Create(path))
	defer f.Close() //nolint:errcheck
	return json.MarshalWrite(f, token)
}
