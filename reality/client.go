package reality

import (
	"net/http"
	"path/filepath"

	"github.com/docker/engine-api/client"
	"github.com/docker/go-connections/tlsconfig"
)

// NewClient initializes a new API client based on the following parameters:
// Use `host` to set the url to the docker server.
// Use `dockerCertPath` to load the TLS certificates from.
// Use `tlsVerify` to enable or disable TLS verification
func NewClient(host, dockerCertPath string, tlsVerify bool) (client.APIClient, error) {
	options := tlsconfig.Options{
		CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
		CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
		KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
		InsecureSkipVerify: !tlsVerify,
	}
	tlsc, err := tlsconfig.Client(options)
	if err != nil {
		return nil, err
	}

	httpCli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}

	cli, err := client.NewClient(host, client.DefaultVersion, httpCli, nil)
	if err != nil {
		return cli, err
	}
	return cli, nil
}
