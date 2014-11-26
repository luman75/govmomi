/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package govmomi

import (
	"path"

	"github.com/luman75/govmomi/vim25/types"
)

type Network struct {
	types.ManagedObjectReference

	InventoryPath string

	c *Client
}

func NewNetwork(c *Client, ref types.ManagedObjectReference) *Network {
	return &Network{
		ManagedObjectReference: ref,
		c: c,
	}
}

func (n Network) Reference() types.ManagedObjectReference {
	return n.ManagedObjectReference
}

func (n Network) Name() string {
	return path.Base(n.InventoryPath)
}

// EthernetCardBackingInfo returns the VirtualDeviceBackingInfo for this Network
func (n Network) EthernetCardBackingInfo() (types.BaseVirtualDeviceBackingInfo, error) {
	name := n.Name()

	backing := &types.VirtualEthernetCardNetworkBackingInfo{
		VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
			DeviceName: name,
		},
	}

	return backing, nil
}
