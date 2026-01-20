//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package platform

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
	"sync"
)

//=============================================================================
//===
//=== Properties
//===
//=============================================================================

var systems = struct {
	sync.RWMutex
	m map[string]*System
	l *[]System
}{}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func GetSystem(c *auth.Context, code string) (*System, error) {
	c.Log.Info("GetSystem: Getting system", "code", code)

	systems.Lock()
	defer systems.Unlock()

	if systems.l != nil {
		c.Log.Info("GetSystem: Returning cached data")
		return systems.m[code], nil
	}

	err := loadSystems(c)

	if err != nil {
		c.Log.Info("GetSystem: Could not retrieve system")
		return nil, err
	}

	c.Log.Info("GetSystem: Returning system")
	return systems.m[code], nil
}

//=============================================================================

func GetSystems(c *auth.Context) (*[]System, error) {
	c.Log.Info("GetSystems: Getting systems...")

	systems.Lock()
	defer systems.Unlock()

	if systems.l != nil {
		c.Log.Info("GetSystems: Returning cached data")
		return systems.l, nil
	}

	err := loadSystems(c)

	if err != nil {
		c.Log.Info("GetSystems: Could not retrieve systems")
		return nil, err
	}

	c.Log.Info("GetSystems: Returning systems", "systems", len(*systems.l))
	return systems.l, nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func loadSystems(c *auth.Context) error {
	c.Log.Info("loadSystems: Retrieving systems from system adapter...")

	var systemList SystemList

	client := req.GetClient("bf")
	url := c.Config.(*app.Config).Platform.System + "/v1/adapters"
	err := req.DoGet(client, url, &systemList, c.Token)

	if err != nil {
		c.Log.Error("loadSystems: Got an error from system adapter ", "error", err.Error())
		return req.NewServerError("Cannot communicate with system-adapter: %v", err.Error())
	}

	sysMap := map[string]*System{}

	for _, s := range systemList.Result {
		ss := s
		sysMap[s.Code] = &ss
	}

	c.Log.Info("loadSystems: Systems loaded", "systems", len(systemList.Result))
	systems.m = sysMap
	systems.l = &systemList.Result
	return nil
}

//=============================================================================
