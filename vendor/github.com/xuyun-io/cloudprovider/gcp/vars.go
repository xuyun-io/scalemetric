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
	"github.com/xuyun-io/cloudprovider/types"
)

var (
	regionList = regions()
)

var (
	// NorthAmerica return all area of north america
	NorthAmerica = types.Area{
		Name: "North America",
		Regions: []types.Region{
			{
				RegionID: "us-east1",
				Name:     "US East (South Carolina)",
			},
			{
				RegionID: "us-central1",
				Name:     "US Central (Iowa)",
			},
			{
				RegionID: "northamerica-northeast1",
				Name:     "Canada (Montréal)",
			},
			{
				RegionID: "us-east4",
				Name:     "US East (Northern Virginia)",
			},
			{
				RegionID: "us-west1",
				Name:     "US West (Oregon)",
			},
			{
				RegionID: "us-west2",
				Name:     "US West (Los Angeles)",
			},
		},
	}
	// Europe return all zones of europe
	Europe = types.Area{
		Name: "Europe",
		Regions: []types.Region{
			{
				RegionID: "europe-north1",
				Name:     "EU (Finland)",
			},
			{
				RegionID: "europe-west1",
				Name:     "EU (Belgium)",
			},
			{
				RegionID: "europe-west2",
				Name:     "EU (London)",
			},
			{
				RegionID: "europe-west3",
				Name:     "EU (Frankfurt)",
			},
			{
				RegionID: "europe-west4",
				Name:     "EU (Netherlands)",
			},
			{
				RegionID: "europe-west6",
				Name:     "EU (Zurich)",
			},
		},
	}

	//Australia return all zones of australia.
	Australia = types.Area{
		Name: "Australia",
		Regions: []types.Region{
			{
				RegionID: "australia-southeast1",
				Name:     "Asia Pacific (Sydney)",
			},
		},
	}
	// Asia return all zones of asia
	Asia = types.Area{
		Name: "Asia",
		Regions: []types.Region{
			{
				RegionID: "asia-northeast1",
				Name:     "Asia Pacific (Tokyo)",
			},
			{
				RegionID: "asia-northeast2",
				Name:     "Asia Pacific (Osaka)",
			},
			{
				RegionID: "asia-southeast1",
				Name:     "Asia Pacific (Singapore)",
			},
			{
				RegionID: "asia-east1",
				Name:     "Asia Pacific (Taiwan)",
			},
			{
				RegionID: "asia-south1",
				Name:     "Asia Pacific (Mumbai)",
			},
			{
				RegionID: "asia-east2",
				Name:     "Asia Pacific (Hong Kong)",
			},
		},
	}

	// SouthAmerica return all zones of south america
	SouthAmerica = types.Area{
		Name: "South America",
		Regions: []types.Region{
			{
				RegionID: "southamerica-east1",
				Name:     "South America (São Paulo)",
			},
		},
	}
)

var (
	zoneList = []types.ZoneList{
		// us
		{
			RegionID: "us-central1",
			Item: []types.Zone{
				{ID: "us-central1-a"}, {ID: "us-central1-b"}, {ID: "us-central1-c"}, {ID: "us-central1-f"},
			},
		},
		{
			RegionID: "us-east1",
			Item: []types.Zone{
				{ID: "us-east1-b"}, {ID: "us-east1-c"}, {ID: "us-east1-d"},
			},
		},

		{
			RegionID: "us-east4",
			Item: []types.Zone{
				{ID: "us-east4-a"}, {ID: "us-east4-b"}, {ID: "us-east4-c"},
			},
		},
		{
			RegionID: "us-west1",
			Item: []types.Zone{
				{ID: "us-west1-a"}, {ID: "us-west1-b"}, {ID: "us-west1-c"},
			},
		},
		{
			RegionID: "us-west2",
			Item: []types.Zone{
				{ID: "us-west2-a"}, {ID: "us-west2-b"}, {ID: "us-west2-c"},
			},
		},
		// north merica
		{
			RegionID: "northamerica-northeast1",
			Item: []types.Zone{
				{ID: "northamerica-northeast1-a"}, {ID: "northamerica-northeast1-b"}, {ID: "northamerica-northeast1-c"},
			},
		},
		// eu
		{
			RegionID: "europe-north1",
			Item: []types.Zone{
				{ID: "europe-north1-a"},
				{ID: "europe-north1-b"},
				{ID: "europe-north1-c"},
			},
		},
		{
			RegionID: "europe-west1",
			Item: []types.Zone{
				{ID: "europe-west1-b"},
				{ID: "europe-west1-c"},
				{ID: "europe-west1-d"},
			},
		},
		{
			RegionID: "europe-west2",
			Item: []types.Zone{
				{ID: "europe-west2-a"},
				{ID: "europe-west2-b"},
				{ID: "europe-west2-c"},
			},
		},
		{
			RegionID: "europe-west3",
			Item: []types.Zone{
				{ID: "europe-west3-a"},
				{ID: "europe-west3-b"},
				{ID: "europe-west3-c"},
			},
		},
		{
			RegionID: "europe-west4",
			Item: []types.Zone{
				{ID: "europe-west4-a"},
				{ID: "europe-west4-b"},
				{ID: "europe-west4-c"},
			},
		},
		{
			RegionID: "europe-west6",
			Item: []types.Zone{
				{ID: "europe-west6-a"},
				{ID: "europe-west6-b"},
				{ID: "europe-west6-c"},
			},
		},
		// asia
		{
			RegionID: "asia-northeast1",
			Item: []types.Zone{
				{ID: "asia-northeast1-a"},
				{ID: "asia-northeast1-b"},
				{ID: "asia-northeast1-c"},
			},
		},
		{
			RegionID: "asia-northeast2",
			Item: []types.Zone{
				{ID: "asia-northeast2-a"},
				{ID: "asia-northeast2-b"},
				{ID: "asia-northeast2-c"},
			},
		},
		{
			RegionID: "asia-southeast1",
			Item: []types.Zone{
				{ID: "asia-southeast1-a"},
				{ID: "asia-southeast1-b"},
				{ID: "asia-southeast1-c"},
			},
		},
		{
			RegionID: "asia-south1",
			Item: []types.Zone{
				{ID: "asia-south1-a"},
				{ID: "asia-south1-b"},
				{ID: "asia-south1-c"},
			},
		},
		{
			RegionID: "asia-east1",
			Item: []types.Zone{
				{ID: "asia-east1-a"},
				{ID: "asia-east1-b"},
				{ID: "asia-east1-c"},
			},
		},
		{
			RegionID: "asia-east2",
			Item: []types.Zone{
				{ID: "asia-east2-a"},
				{ID: "asia-east2-b"},
				{ID: "asia-east2-c"},
			},
		},
		{
			RegionID: "southamerica-east1",
			Item: []types.Zone{
				{ID: "southamerica-east1-a"},
				{ID: "southamerica-east1-b"},
				{ID: "southamerica-east1-c"}},
		},
		// australia
		{
			RegionID: "australia-southeast1",
			Item: []types.Zone{
				{ID: "australia-southeast1-a"},
				{ID: "australia-southeast1-b"},
				{ID: "australia-southeast1-c"},
			},
		},
	}
)

func regions() *types.RegionList {
	gcpRegionList := &types.RegionList{
		Item: []types.Area{
			Asia,
			NorthAmerica,
			SouthAmerica,
			Europe,
			Australia,
		},
	}
	for i := range gcpRegionList.Item {
		if len(gcpRegionList.Item[i].Regions) <= 0 {
			continue
		}
		for j := range gcpRegionList.Item[i].Regions {
			gcpRegionList.Item[i].Regions[j].Zones = Zones(gcpRegionList.Item[i].Regions[j].RegionID)
		}
	}
	return gcpRegionList
}

// Zones return all zones by regionID
func Zones(regionID string) types.ZoneList {
	for i := range zoneList {
		if regionID == zoneList[i].RegionID {
			return zoneList[i]
		}
	}
	return types.ZoneList{}
}

/*
link: https://cloud.google.com/compute/docs/regions-zones/?hl=zh-cn#available
*/
