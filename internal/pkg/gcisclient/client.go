package gcisclient

import (
	"net/url"
	"os"

	"github.com/minchao/go-gcis/gcis"
)

var (
	client *gcis.Client
)

func New() *gcis.Client {
	if client == nil {
		client = gcis.NewClient()
		if baseURL := os.Getenv("GCIS_BASE_URL"); baseURL != "" {
			client.BaseURL, _ = url.Parse(baseURL)
		}
	}
	return client
}
