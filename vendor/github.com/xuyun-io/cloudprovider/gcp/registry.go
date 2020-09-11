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

package gcp

import (
	"context"
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/oauth2/google"
	container "google.golang.org/api/container/v1"

	"github.com/xuyun-io/cloudprovider/registry"
	"github.com/xuyun-io/cloudprovider/types"
)

// RepoClient define gcp repository client
type RepoClient struct {
	registry.Client
}

var (
	_ types.RepoInterface = &RepoClient{}
)

// NewRepoClient return gcp repository client
func NewRepoClient(ctx context.Context, host string, cfg Config) types.RepoInterface {
	cfgByts, err := json.Marshal(cfg)
	if err != nil {
		return &RepoClient{Client: registry.Client{Error: err}}
	}
	creds, err := google.CredentialsFromJSON(ctx, cfgByts, container.CloudPlatformScope)
	if err != nil {
		return &RepoClient{Client: registry.Client{Error: err}}
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return &RepoClient{Client: registry.Client{Error: err}}
	}
	opt := registry.NewBearerTokenAuth(token.AccessToken)
	return &RepoClient{Client: *newRepoClient(host, opt)}
}

func newRepoClient(host string, options ...remote.Option) *registry.Client {
	if len(options) <= 0 {
		return &registry.Client{Host: host}
	}
	return &registry.Client{Host: host, Options: options}
}
