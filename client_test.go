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
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/luman75/govmomi/test"
	"github.com/luman75/govmomi/vim25/mo"
)

func TestLogin(t *testing.T) {
	u := test.URL()
	if u == nil {
		t.SkipNow()
	}

	c, err := NewClient(*u, true)
	if err != nil {
		t.Error(err)
	}

	f := func() error {
		var x mo.Folder
		err = mo.RetrieveProperties(c, c.ServiceContent.PropertyCollector, c.ServiceContent.RootFolder, &x)
		if err != nil {
			return err
		}
		if len(x.Name) == 0 {
			return errors.New("invalid response") // TODO: RetrieveProperties should propagate fault
		}
		return nil
	}

	// check cookie is valid with an sdk request
	if err := f(); err != nil {
		t.Error(err)
	}

	// check cookie is valid with a non-sdk request
	u.User = nil // turn off Basic auth
	u.Path = "/folder"
	r, err := c.Client.Get(u.String())
	if err != nil {
		t.Error(err)
	}
	if r.StatusCode != http.StatusOK {
		t.Error(r)
	}

	// sdk request should fail w/o a valid cookie
	c.Client.Jar = nil
	if err := f(); err == nil {
		t.Error("should fail")
	}

	// invalid login
	u.Path = "/sdk"
	u.User = url.UserPassword("ENOENT", "EINVAL")
	_, err = NewClient(*u, true)
	if err == nil {
		t.Error("should fail")
	}
}

func TestInvalidSdk(t *testing.T) {
	u := test.URL()
	if u == nil {
		t.SkipNow()
	}

	// a URL other than a valid /sdk should error, not panic
	u.Path = "/mob"
	_, err := NewClient(*u, true)
	if err == nil {
		t.Error("should fail")
	}
}
