//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================


package platform

import (
	"sync"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
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

	client := req.GetDefaultClient()
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
