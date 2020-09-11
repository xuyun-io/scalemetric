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

package types

import (
	"time"
)

var (

	// provider supported list
	providerMaps = map[Provider]bool{
		ProviderKubeconfig: true,
		ProviderAWS:        true,
		ProviderGCP:        true,
		ProviderALiCloud:   true,
		ProviderHuawei:     true,
		ProviderTencent:    true,
	}
)

// Provider define provider
type Provider string

// provider list
const (
	ProviderAWS        Provider = "ProviderAWS"
	ProviderGCP        Provider = "ProviderGCP"
	ProviderALiCloud   Provider = "ProviderALiCloud"
	ProviderKubeconfig Provider = "ProviderKubeconfig"
	ProviderHuawei     Provider = "ProviderHuawei"
	ProviderTencent    Provider = "ProviderTencent"
)

// SupportProviders return supported provider list
func SupportProviders() map[Provider]bool {
	return providerMaps
}

// DescribeCluster define describe cluster struct
type DescribeCluster struct {
	ClusterMeta
}

// ClusterMeta define cluster metadata
type ClusterMeta struct {
	UseUnique string `json:"useUnique"`
	ZoneID    string `json:"zoneId"`
	RegionID  string `json:"regionId"`
	Name      string `json:"name"`
	State     string `json:"state"`
	ClusterID string `json:"clusterId"`
}

// ClusterList define cluster list struct
type ClusterList struct {
	Item []DescribeCluster `json:"item,omitempty"`
}

// Zone define zone, zone include id and name
type Zone struct {
	ID   string
	Name string
}

// ZoneList define zone list
type ZoneList struct {
	RegionID string `json:"regionId,omitempty"`
	Item     []Zone `json:"item,omitempty"`
}

// Region define region
type Region struct {
	RegionID       string   `json:"regionId,omitempty"`
	Name           string   `json:"name,omitempty"`
	LocalName      string   `json:"localName"`
	RegionEndpoint string   `json:"regionEndpoint,omitempty"`
	Zones          ZoneList `json:"zones,omitempty"`
}

// Area define area
type Area struct {
	Name    string   `json:"name,omitempty"`
	Regions []Region `json:"regions,omitempty"`
}

// RegionList define region list
type RegionList struct {
	Item []Area `json:"item,omitempty"`
}

// Repository define repository details
type Repository struct {
	// The date and time, in JavaScript date format, when the repository was created.
	CreatedAt      *time.Time `locationName:"createdAt" json:"createdAt"`
	RegistryID     *string    `locationName:"registryId" json:"registryId"`
	RepositoryName *string    `locationName:"repositoryName" json:"repositoryName"`
}

// ImageDetail define image detail
type ImageDetail struct {
	// The sha256 digest of the image manifest.
	ImageDigest      *string    `locationName:"imageDigest" json:"imageDigest"`
	ImagePushedAt    *time.Time `locationName:"imagePushedAt" json:"imagePushedAt"`
	ImageSizeInBytes *int64     `locationName:"imageSizeInBytes" json:"imageSizeInBytes"`
	// The list of tags associated with this image.
	ImageTags []*string `locationName:"imageTags" type:"list" json:"imageTags"`
	// The AWS account ID associated with the registry to which this image belongs.
	RegistryID *string `locationName:"registryId" json:"registryId"`
	// The name of the repository to which this image belongs.
	RepositoryName *string `locationName:"repositoryName" json:"repositoryName"`
}

// NewRegionList return region list
func NewRegionList() *RegionList {
	return &RegionList{
		Item: []Area{},
	}
}

// AddRegion add region
func (rl *RegionList) AddRegion(areaName string, r Region) {
	if len(rl.Item) <= 0 {
		rl.Item = []Area{
			{
				Name: areaName,
				Regions: []Region{
					r,
				},
			},
		}
		return
	}
	for i := range rl.Item {
		if areaName == rl.Item[i].Name {
			if len(rl.Item[i].Regions) <= 0 {
				rl.Item[i].Regions = []Region{r}
				return
			}
			for j := range rl.Item[i].Regions {
				if rl.Item[i].Regions[j].RegionID == r.RegionID {
					rl.Item[i].Regions[j] = r
					return
				}
			}
			rl.Item[i].Regions = append(rl.Item[i].Regions, r)
			return
		}
	}
	rl.Item = append(rl.Item, Area{
		Name:    areaName,
		Regions: []Region{r},
	})
}

// AddZone add zone
func (rl *RegionList) AddZone(regionID string, zone Zone) {
	if len(rl.Item) <= 0 {
		return
	}
	for i := range rl.Item {
		if len(rl.Item[i].Regions) <= 0 {
			return
		}
		for j := range rl.Item[i].Regions {
			if regionID == rl.Item[i].Regions[j].RegionID {
				if len(rl.Item[i].Regions[j].Zones.Item) <= 0 {
					rl.Item[i].Regions[j].Zones.Item = []Zone{zone}
					return
				}
				for k := range rl.Item[i].Regions[j].Zones.Item {
					if rl.Item[i].Regions[j].Zones.Item[k].ID == zone.ID {
						rl.Item[i].Regions[j].Zones.Item[k] = zone
						return
					}
				}
				rl.Item[i].Regions[j].Zones.Item = append(rl.Item[i].Regions[j].Zones.Item, zone)
				return
			}
		}
	}
}

// SliceToRepositories return repositories
func SliceToRepositories(repos []string) ([]*Repository, error) {
	if len(repos) <= 0 {
		return make([]*Repository, 0), nil
	}
	rs := make([]*Repository, 0)
	for i := range repos {
		rs = append(rs, &Repository{RepositoryName: &repos[i]})
	}
	return rs, nil
}

// SliceToImageDetails return image tag list
func SliceToImageDetails(repo string, tags []string) ([]*ImageDetail, error) {
	imageDetails := make([]*ImageDetail, 0)
	if len(tags) <= 0 {
		return imageDetails, nil
	}
	for i := range tags {
		imageDetails = append(imageDetails, &ImageDetail{ImageTags: []*string{&tags[i]}, RepositoryName: &repo})
	}
	return imageDetails, nil
}
