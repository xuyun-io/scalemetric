/*
Provider™ is a cloud native CaaS platform.
Copyright (C) 2019  Xuyun Authors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package registry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/xuyun-io/cloudprovider/types"
)

// Client define docker register client
type Client struct {
	Host    string
	Options []remote.Option
	Error   error
}

var _ types.RepoInterface = &Client{}

// NewBearerTokenAuth return bearer token auth style
func NewBearerTokenAuth(bearerToken string) remote.Option {
	return remote.WithAuth(&authn.Bearer{Token: bearerToken})
}

// NewUsernamePasswordAuth return auth  for username and password
func NewUsernamePasswordAuth(username, password string) remote.Option {
	basic := authn.Basic{
		Username: username,
		Password: password,
	}
	return remote.WithAuth(&basic)
}

// NewCAAuth return auth for ca.
func NewCAAuth(caByts []byte) remote.Option {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caByts)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}
	return remote.WithTransport(tr)

}

// NewClient return docker registry client
func NewClient(host string, options ...remote.Option) types.RepoInterface {
	if len(options) <= 0 {
		return &Client{Host: host}
	}
	return &Client{Host: host, Options: options}
}

// Repositories return repo list
func (dr *Client) Repositories(ctx context.Context) ([]*types.Repository, error) {
	if dr.Error != nil {
		return nil, dr.Error
	}
	rg, err := name.NewRegistry(dr.Host)
	if err != nil {
		return nil, err
	}
	repos, err := remote.Catalog(ctx, rg, dr.Options...)
	if err != nil {
		return nil, err
	}
	return types.SliceToRepositories(repos)
}

// TagList return repository tag list
func (dr *Client) TagList(ctx context.Context, repoName string) ([]*types.ImageDetail, error) {
	if dr.Error != nil {
		return nil, dr.Error
	}
	repo, err := name.NewRepository(repoName)
	if err != nil {
		return nil, err
	}
	tags, err := remote.ListWithContext(ctx, repo, dr.Options...)
	if err != nil {
		return nil, err
	}
	return types.SliceToImageDetails(repoName, tags)
}

// ImageDescribe return image details
func (dr *Client) ImageDescribe(ctx context.Context, imageName string) (interface{}, error) {
	if dr.Error != nil {
		return nil, dr.Error
	}
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, err
	}
	return remote.Image(ref, dr.Options...)
}
