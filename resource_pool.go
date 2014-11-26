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
	"github.com/luman75/govmomi/vim25/methods"
	"github.com/luman75/govmomi/vim25/types"
)

type ResourcePool struct {
	types.ManagedObjectReference

	c *Client
}

func NewResourcePool(c *Client, ref types.ManagedObjectReference) *ResourcePool {
	return &ResourcePool{
		ManagedObjectReference: ref,
		c: c,
	}
}

func (p ResourcePool) Reference() types.ManagedObjectReference {
	return p.ManagedObjectReference
}

func (p ResourcePool) ImportVApp(spec types.BaseImportSpec, folder *Folder, host *HostSystem) (*HttpNfcLease, error) {
	req := types.ImportVApp{
		This: p.Reference(),
		Spec: spec,
	}

	if folder != nil {
		ref := folder.Reference()
		req.Folder = &ref
	}

	if host != nil {
		ref := host.Reference()
		req.Host = &ref
	}

	res, err := methods.ImportVApp(p.c, &req)
	if err != nil {
		return nil, err
	}

	return NewHttpNfcLease(p.c, res.Returnval), nil
}
