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

package serial

import (
	"flag"

	"github.com/luman75/govmomi/govc/cli"
	"github.com/luman75/govmomi/govc/flags"
)

type disconnect struct {
	*flags.VirtualMachineFlag

	device string
}

func init() {
	cli.Register("device.serial.disconnect", &disconnect{})
}

func (cmd *disconnect) Register(f *flag.FlagSet) {
	f.StringVar(&cmd.device, "device", "", "serial port device name")
}

func (cmd *disconnect) Process() error { return nil }

func (cmd *disconnect) Run(f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	devices, err := vm.Device()
	if err != nil {
		return err
	}

	d, err := devices.FindSerialPort(cmd.device)
	if err != nil {
		return err
	}

	return vm.EditDevice(devices.DisconnectSerialPort(d))
}
