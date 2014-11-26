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
	"github.com/luman75/govmomi/vim25/mo"
	"github.com/luman75/govmomi/vim25/types"
)

type Folder struct {
	types.ManagedObjectReference

	c *Client
}

func NewFolder(c *Client, ref types.ManagedObjectReference) *Folder {
	return &Folder{
		ManagedObjectReference: ref,
		c: c,
	}
}

func (f Folder) Reference() types.ManagedObjectReference {
	return f.ManagedObjectReference
}

func (f Folder) Children() ([]Reference, error) {
	var mf mo.Folder

	err := f.c.Properties(f.Reference(), []string{"childEntity"}, &mf)
	if err != nil {
		return nil, err
	}

	var rs []Reference

	for _, e := range mf.ChildEntity {
		if r := NewReference(f.c, e); r != nil {
			rs = append(rs, r)
		}
	}

	return rs, nil
}

func (f Folder) CreateVM(config types.VirtualMachineConfigSpec, pool *ResourcePool, host *HostSystem) (*Task, error) {
	req := types.CreateVM_Task{
		This:   f.Reference(),
		Config: config,
		Pool:   pool.Reference(),
	}

	if host != nil {
		ref := host.Reference()
		req.Host = &ref
	}

	res, err := methods.CreateVM_Task(f.c, &req)
	if err != nil {
		return nil, err
	}

	return NewTask(f.c, res.Returnval), nil
}
