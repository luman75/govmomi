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

package disk

import (
	"flag"

	"github.com/luman75/govmomi/govc/cli"
	"github.com/luman75/govmomi/govc/flags"
	"github.com/luman75/govmomi/vim25/types"
)

type attach struct {
	*flags.DatastoreFlag
	*flags.VirtualMachineFlag

	persist    bool
	link       bool
	disk       string
	controller string
}

func init() {
	cli.Register("vm.disk.attach", &attach{})
}

func (cmd *attach) Register(f *flag.FlagSet) {
	f.BoolVar(&cmd.persist, "persist", true, "Persist attached disk")
	f.BoolVar(&cmd.link, "link", true, "Link specified disk")
	f.StringVar(&cmd.controller, "controller", "", "Disk controller")
	f.StringVar(&cmd.disk, "disk", "", "Disk path name")
}

func (cmd *attach) Process() error { return nil }

func (cmd *attach) Run(f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	ds, err := cmd.Datastore()
	if err != nil {
		return err
	}

	devices, err := vm.Device()
	if err != nil {
		return err
	}

	controller, err := devices.FindDiskController(cmd.controller)
	if err != nil {
		return err
	}

	disk := devices.CreateDisk(controller, ds.Path(cmd.disk))
	backing := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)

	if cmd.link {
		if cmd.persist {
			backing.DiskMode = string(types.VirtualDiskModeIndependent_persistent)
		} else {
			backing.DiskMode = string(types.VirtualDiskModeIndependent_persistent)
		}

		disk = devices.ChildDisk(disk)
		return vm.AddDevice(disk)
	}

	if cmd.persist {
		backing.DiskMode = string(types.VirtualDiskModePersistent)
	} else {
		backing.DiskMode = string(types.VirtualDiskModeNonpersistent)
	}

	return vm.AddDevice(disk)
}
