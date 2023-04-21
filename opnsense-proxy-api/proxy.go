package opnsense_proxy_api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type syncAliasesRequest struct {
	Host    string   `json:"host"`
	Aliases []string `json:"aliases"`
}

func Sync(proxyAddress string, hosts []string, hostname string) error {
	client := resty.New()
	endpoint := fmt.Sprintf("%v/sync", proxyAddress)
	body := syncAliasesRequest{
		Host:    hostname,
		Aliases: hosts,
	}
	_, err := client.R().SetBody(body).Post(endpoint)
	return err
}
