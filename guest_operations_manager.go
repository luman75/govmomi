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

import "github.com/luman75/govmomi/vim25/mo"

type GuestOperationsManager struct {
	c *Client
}

func (m GuestOperationsManager) AuthManager() (*GuestAuthManager, error) {
	var g mo.GuestOperationsManager

	err := m.c.Properties(*m.c.ServiceContent.GuestOperationsManager, []string{"authManager"}, &g)
	if err != nil {
		return nil, err
	}

	return &GuestAuthManager{*g.AuthManager, m.c}, nil
}

func (m GuestOperationsManager) FileManager() (*GuestFileManager, error) {
	var g mo.GuestOperationsManager

	err := m.c.Properties(*m.c.ServiceContent.GuestOperationsManager, []string{"fileManager"}, &g)
	if err != nil {
		return nil, err
	}

	return &GuestFileManager{*g.FileManager, m.c}, nil
}

func (m GuestOperationsManager) ProcessManager() (*GuestProcessManager, error) {
	var g mo.GuestOperationsManager

	err := m.c.Properties(*m.c.ServiceContent.GuestOperationsManager, []string{"processManager"}, &g)
	if err != nil {
		return nil, err
	}

	return &GuestProcessManager{*g.ProcessManager, m.c}, nil
}
