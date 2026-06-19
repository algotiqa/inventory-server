//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"slices"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetAgentProfiles(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.AgentProfile, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	list, err := db.GetAgentProfiles(tx, filter, offset, limit)

	if err != nil {
		return nil, err
	}

	return list, nil
}

//=============================================================================

func GetExternalRefs(tx *gorm.DB, c *auth.Context, id uint) ([]string, error) {
	c.Log.Info("GetExternalRefs: Getting external refs from profile", "id", id)

	ap, err := getAgentProfile(tx, c, id)
	if err != nil {
		return nil, err
	}

	list, err := callAgentToGetExternalRefs(c, ap)
	if err != nil {
		return nil, err
	}

	dbRefs,err := getExternalRefsForProfile(tx, c, id)
	if err != nil {
		return nil, err
	}

	res := []string{}

	for _, xref := range list {
		if _,ok := dbRefs[xref]; !ok {
			res = append(res, xref)
		}
	}

	slices.Sort(res)

	c.Log.Info("GetExternalRefs: Got new list of external refs", "id", id, "size", len(res))
	return res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func callAgentToGetExternalRefs(c *auth.Context, ap *db.AgentProfile) ([]string, error) {
	client := agentscanner.CreateClient(ap.SslCertRef, ap.SslKeyRef, "ca.crt")
	if client == nil {
		return nil, req.NewServerError("Cannot create client for agent: %v", ap.Id)
	}

	var list []string
	err := req.DoGet(client, ap.RemoteUrl + agentscanner.UrlTradingSystems, &list, "")
	if err != nil {
		c.Log.Error("callAgentToGetExternalRefs: Agent raised an error", "id", ap.Id, "error", err.Error())
		return nil, req.NewServiceUnavailableError("Agent raised an error : " + err.Error())
	}

	return list, nil
}

//=============================================================================

func getExternalRefsForProfile(tx *gorm.DB, c *auth.Context, id uint) (map[string]bool, error) {
	filter := map[string]any{}
	filter["username"]         = c.Session.Username
	filter["agent_profile_id"] = id
	list, err := db.GetTradingSystems(tx, filter, 0, 5000)

	if err != nil {
		return nil, err
	}

	result := map[string]bool{}
	for _,ts := range *list {
		result[ts.ExternalRef] = true
	}

	return result, nil
}

//=============================================================================
