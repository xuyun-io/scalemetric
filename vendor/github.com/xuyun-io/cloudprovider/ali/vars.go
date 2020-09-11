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

package ali

import "github.com/xuyun-io/cloudprovider/types"

// define ali provider region list
var (
	regionLists = &types.RegionList{}
	UAE         = types.Area{
		Name: "UAE",
		Regions: []types.Region{
			{
				RegionID: "me-east-1",
				Name:     "UAE (Dubai)",
			},
		},
	}
	US = types.Area{
		Name: "US",
		Regions: []types.Region{
			{
				RegionID: "us-west-1",
				Name:     "US (Silicon Valley)",
			},
			{
				RegionID: "us-east-1",
				Name:     "US (Virginia)",
			},
		},
	}
	UK = types.Area{
		Name: "UK",
		Regions: []types.Region{
			{
				Name:     "UK (London)",
				RegionID: "eu-west-1",
			},
		},
	}
	Germany = types.Area{
		Name: "Germany",
		Regions: []types.Region{
			{
				Name:     "Germany (Frankfurt)",
				RegionID: "eu-central-1",
			},
		},
	}
	India = types.Area{
		Name: "India",
		Regions: []types.Region{
			{
				Name:     "India (Mumbai)",
				RegionID: "ap-south-1",
			},
		},
	}
	Japan = types.Area{
		Name: "Japan",
		Regions: []types.Region{
			{
				RegionID: "ap-northeast-1",
				Name:     "Japan (Tokyo)",
			},
		},
	}
	Indonesia = types.Area{
		Name: "Indonesia",
		Regions: []types.Region{
			{
				RegionID: "ap-southeast-5",
				Name:     "Indonesia (Jakarta)",
			},
		},
	}
	Malaysia = types.Area{
		Name: "Malaysia",
		Regions: []types.Region{
			{
				RegionID: "ap-southeast-3",
				Name:     "Malaysia (Kuala Lumpur)",
			},
		},
	}
	Australia = types.Area{
		Name: "Australia",
		Regions: []types.Region{
			{
				RegionID: "ap-southeast-2",
				Name:     "Australia (Sydney)",
			},
		},
	}
	Singapore = types.Area{
		Name: "Singapore",
		Regions: []types.Region{
			{
				RegionID: "ap-southeast-1",
				Name:     "Singapore",
			},
		},
	}
	China = types.Area{
		Name: "China",
		Regions: []types.Region{
			{
				RegionID: "cn-hangzhou",
				Name:     "China (Shanghai)",
			},
			{

				RegionID: "cn-shanghai",
				Name:     "China (Shanghai)",
			},
			{

				RegionID: "cn-qingdao",
				Name:     "China (Qingdao)",
			},
			{
				Name:     "China (Beijing)",
				RegionID: "cn-beijing",
			},
			{
				Name:     "China (Zhangjiakou)",
				RegionID: "cn-zhangjiakou",
			},
			{
				Name:     "China (Hohhot)",
				RegionID: "cn-huhehaote",
			},
			{
				Name:     "China (Shenzhen)",
				RegionID: "cn-shenzhen",
			},
			{
				Name:     "China (Chengdu)",
				RegionID: "cn-chengdu",
			},
			{
				Name:     "China (Hong Kong)",
				RegionID: "cn-hongkong",
			},
		},
	}
)

func regions() *types.RegionList {
	rl := &types.RegionList{
		Item: []types.Area{
			China,
			Singapore,
			Australia,
			Malaysia,
			Indonesia,
			Japan,
			India,
			Germany,
			UK,
			US,
			UAE,
		},
	}
	for i := range rl.Item {
		if len(rl.Item) <= 0 {
			continue
		}
		for j := range rl.Item[i].Regions {
			zl := types.ZoneList{
				RegionID: rl.Item[i].Regions[j].RegionID,
				Item: []types.Zone{
					{
						ID: rl.Item[i].Regions[j].RegionID,
					},
				},
			}
			rl.Item[i].Regions[j].Zones = zl
		}
	}
	return rl
}
