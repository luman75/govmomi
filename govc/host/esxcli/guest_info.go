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

package esxcli

import (
	"strings"

	"github.com/luman75/govmomi"
	"github.com/luman75/govmomi/vim25/mo"
	"github.com/luman75/govmomi/vim25/types"
)

type hostInfo struct {
	*Executor
	wids map[string]string
}

type GuestInfo struct {
	c     *govmomi.Client
	hosts map[string]*hostInfo
}

func NewGuestInfo(c *govmomi.Client) *GuestInfo {
	return &GuestInfo{
		c:     c,
		hosts: make(map[string]*hostInfo),
	}
}

func (g *GuestInfo) hostInfo(ref *types.ManagedObjectReference) (*hostInfo, error) {
	// cache exectuor and uuid -> worldid map
	if h, ok := g.hosts[ref.Value]; ok {
		return h, nil
	}

	host := govmomi.NewHostSystem(g.c, *ref)

	e, err := NewExecutor(g.c, host)
	if err != nil {
		return nil, err
	}

	res, err := e.Run([]string{"vm", "process", "list"})
	if err != nil {
		return nil, err
	}

	ids := make(map[string]string, len(res.Values))

	for _, process := range res.Values {
		// Normalize uuid, esxcli and mo.VirtualMachine have different formats
		uuid := strings.Replace(process["UUID"][0], " ", "", -1)
		uuid = strings.Replace(uuid, "-", "", -1)

		ids[uuid] = process["WorldID"][0]
	}

	h := &hostInfo{e, ids}
	g.hosts[ref.Value] = h

	return h, nil
}

// IpAddress attempts to find the guest IP address using esxcli.
// ESX hosts must be configured with the /Net/GuestIPHack enabled.
// For example:
// $ govc host.esxcli -- system settings advanced set -o /Net/GuestIPHack -i 1
func (g *GuestInfo) IpAddress(vm *govmomi.VirtualMachine) (string, error) {
	var mvm mo.VirtualMachine
	err := g.c.Properties(vm.ManagedObjectReference, []string{"runtime.host", "config.uuid"}, &mvm)
	if err != nil {
		return "", err
	}

	h, err := g.hostInfo(mvm.Runtime.Host)
	if err != nil {
		return "", err
	}

	// Normalize uuid, esxcli and mo.VirtualMachine have different formats
	uuid := strings.Replace(mvm.Config.Uuid, "-", "", -1)

	if wid, ok := h.wids[uuid]; ok {
		res, err := h.Run([]string{"network", "vm", "port", "list", "--world-id", wid})
		if err != nil {
			return "", err
		}

		if len(res.Values) == 1 {
			if ip, ok := res.Values[0]["IPAddress"]; ok {
				return ip[0], nil
			}
		}
	}

	return "", nil
}
