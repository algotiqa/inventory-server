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

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.TradingSystemFull, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	if details {
		return db.GetTradingSystemsFull(tx, filter, offset, limit)
	}

	return db.GetTradingSystems(tx, filter, offset, limit)
}

//=============================================================================

func AddTradingSystem(tx *gorm.DB, c *auth.Context, tss *TradingSystemSpec) (*db.TradingSystem, error) {
	c.Log.Info("AddTradingSystem: Adding a new trading system", "name", tss.Name)

	//TODO: validate type

	var ts db.TradingSystem
	ts.Username = c.Session.Username
	ts.DataProductId = tss.DataProductId
	ts.BrokerProductId = tss.BrokerProductId
	ts.TradingSessionId = tss.TradingSessionId
	ts.AgentProfileId = tss.AgentProfileId
	ts.Name = tss.Name
	ts.Timeframe = tss.Timeframe
	ts.StrategyType = tss.StrategyType
	ts.Overnight = tss.Overnight
	ts.Tags = tss.Tags
	ts.ExternalRef = tss.ExternalRef

	//--- External systems are already finalized

	if ts.AgentProfileId != nil {
		ts.Finalized = true
	}

	err := db.AddTradingSystem(tx, &ts)
	if err != nil {
		c.Log.Error("AddTradingSystem: Could not add a new trading system", "error", err.Error())
		return nil, err
	}

	err = sendChangeMessage(tx, c, &ts, msg.TypeCreate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("AddTradingSystem: Trading system added", "id", ts.Id)
	return &ts, err
}

//=============================================================================

func UpdateTradingSystem(tx *gorm.DB, c *auth.Context, id uint, tss *TradingSystemSpec) (*db.TradingSystem, error) {
	c.Log.Info("UpdateTradingSystem: Updating a trading system", "id", id, "name", tss.Name)

	ts, err := getTradingSystem(tx, c, id, "UpdateTradingSystem")
	if err != nil {
		return nil, err
	}

	//TODO: validate type

	ts.DataProductId = tss.DataProductId
	ts.BrokerProductId = tss.BrokerProductId
	ts.TradingSessionId = tss.TradingSessionId
	ts.AgentProfileId = tss.AgentProfileId
	ts.Name = tss.Name
	ts.Timeframe = tss.Timeframe
	ts.StrategyType = tss.StrategyType
	ts.Overnight = tss.Overnight
	ts.Tags = tss.Tags
	ts.ExternalRef = tss.ExternalRef

	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		c.Log.Error("UpdateTradingSystem: Could not update a trading system", "error", err.Error(), "id", ts.Id)
		return nil, err
	}

	err = sendChangeMessage(tx, c, ts, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdateTradingSystem: Trading system updated", "id", ts.Id, "name", ts.Name)
	return ts, nil
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, c *auth.Context, id uint) (*db.TradingSystem, error) {
	c.Log.Info("DeleteTradingSystem: Deleting trading system", "id", id)

	ts, err := getTradingSystem(tx, c, id, "DeleteTradingSystem")
	if err != nil {
		return nil, err
	}

	err = db.DeleteTradingSystem(tx, id)
	if err != nil {
		c.Log.Error("DeleteTradingSystem: Cannot delete trading system", "id", id, "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	tsm := TradingSystemMessage{}
	tsm.TradingSystem = ts
	err = msg.SendMessage(msg.ExInventory, msg.SourceTradingSystem, msg.TypeDelete, &tsm)

	if err != nil {
		c.Log.Error("DeleteTradingSystem: Could not publish the delete message", "id", id, "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	c.Log.Info("DeleteTradingSystem: Trading system deleted", "id", id, "name", ts.Name)
	return ts, nil
}

//=============================================================================

const (
	ResponseStatusOk      = "ok"
	ResponseStatusSkipped = "skipped"
)

//-----------------------------------------------------------------------------

type FinalizationResponse struct {
	Status string `json:"status"`
}

//-----------------------------------------------------------------------------

func FinalizeTradingSystem(tx *gorm.DB, c *auth.Context, id uint) (*FinalizationResponse, error) {
	c.Log.Info("FinalizeTradingSystem: Finalizing trading system", "id", id)

	ts, err := getTradingSystem(tx, c, id, "FinalizeTradingSystem")
	if err != nil {
		return nil, err
	}

	if ts.Finalized {
		return &FinalizationResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	ts.Finalized = true
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		c.Log.Error("FinalizeTradingSystem: Cannot finalize trading system", "id", id, "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	err = sendChangeMessage(tx, c, ts, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("FinalizeTradingSystem: Trading system finalized", "id", ts.Id, "name", ts.Name)
	return &FinalizationResponse{
		Status: ResponseStatusOk,
	}, nil
}

//=============================================================================

type TradingSystemReloadRequest struct {
	Name string `json:"name"`
}

//-----------------------------------------------------------------------------

func ReloadTradesFromAgent(tx *gorm.DB, c *auth.Context, id uint) (*TradingSystemReloadResponse, error) {
	c.Log.Info("ReloadTradesFromAgent: Reloading trades for trading system", "id", id)

	ts, err := getTradingSystem(tx, c, id, "ReloadTradingSystem")
	if err != nil {
		return nil, err
	}

	if ts.AgentProfileId == nil {
		return nil, req.NewBadRequestError("Trading system has no agent associated")
	}

	var ap *db.AgentProfile
	ap, err = getAgentProfile(tx, c, *ts.AgentProfileId)
	if err != nil {
		return nil, err
	}

	ats, err := callAgentToReloadSystem(c, ap, ts)
	if err != nil {
		return nil, err
	}

	location, err := agentscanner.GetLocation(tx, ts)
	if err != nil {
		c.Log.Error("ReloadTradesFromAgent: Cannot retrieve timezone for trading system.", "id", ts.Id, "error", err)
		return nil, req.NewServerErrorByError(err)
	}

	err = agentscanner.SendTrades(ts, ats, location, true)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	res := NewTradingSystemReloadResponse()
	for _, tl := range ats.TradeLists {
		res.TradeCount[tl.FileName] = len(tl.Trades)
	}

	c.Log.Info("ReloadTradesFromAgent: Trading system has been reloaded", "id", id, "fileCount", len(ats.TradeLists))
	return res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getTradingSystem(tx *gorm.DB, c *auth.Context, id uint, funcName string) (*db.TradingSystem, error) {
	ts, err := db.GetTradingSystemById(tx, id)
	if err != nil {
		c.Log.Error(funcName+": Could not retrieve trading system", "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	if ts == nil {
		c.Log.Error(funcName+": Trading system was not found", "id", id)
		return nil, req.NewNotFoundError("Trading system was not found: %v", id)
	}

	if ts.Username != c.Session.Username {
		c.Log.Error(funcName+": Trading system not owned by user", "id", id)
		return nil, req.NewForbiddenError("Trading system is not owned by user: %v", id)
	}

	return ts, nil
}

//=============================================================================

func sendChangeMessage(tx *gorm.DB, c *auth.Context, ts *db.TradingSystem, msgType int) error {
	dp, err := db.GetDataProductById(tx, ts.DataProductId)
	if err != nil {
		c.Log.Error("sendChangeMessage: Could not retrieve data product of TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	bp, err := db.GetBrokerProductById(tx, ts.BrokerProductId)
	if err != nil {
		c.Log.Error("sendChangeMessage: Could not retrieve broker product of TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	ex, err := db.GetExchangeById(tx, bp.ExchangeId)
	if err != nil {
		c.Log.Error("sendChangeMessage: Could not retrieve exchange of TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	cu, err := db.GetCurrencyById(tx, ex.CurrencyId)
	if err != nil {
		c.Log.Error("sendChangeMessage: Could not retrieve currency of TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	se, err := db.GetTradingSessionById(tx, ts.TradingSessionId)
	if err != nil {
		c.Log.Error("sendChangeMessage: Could not retrieve trading session of TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	var ap *db.AgentProfile

	if ts.AgentProfileId != nil {
		ap, err = db.GetAgentProfileById(tx, *ts.AgentProfileId)
		if err != nil {
			c.Log.Error("sendChangeMessage: Could not retrieve agent profile of TS", "error", err.Error(), "id", ts.Id)
			return err
		}
	}

	tsm := TradingSystemMessage{ts, dp, bp, cu, se, ap, ex}
	err = msg.SendMessage(msg.ExInventory, msg.SourceTradingSystem, msgType, &tsm)

	if err != nil {
		c.Log.Error("sendChangeMessage: Could not publish the update message for TS", "error", err.Error(), "id", ts.Id)
		return err
	}

	return nil
}

//=============================================================================

func getAgentProfile(tx *gorm.DB, c *auth.Context, id uint) (*db.AgentProfile, error) {
	ap, err := db.GetAgentProfileById(tx, id)
	if err != nil {
		c.Log.Error("getAgentProfile: Could not retrieve agent profile", "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	if ap == nil {
		c.Log.Error("getAgentProfile: Agent profile was not found", "id", id)
		return nil, req.NewNotFoundError("Agent profile was not found: %v", id)
	}

	if ap.Username != c.Session.Username {
		c.Log.Error("getAgentProfile: Agent profile not owned by user", "id", id)
		return nil, req.NewForbiddenError("Agent profile is not owned by user: %v", id)
	}

	return ap, nil
}

//=============================================================================

func callAgentToReloadSystem(c *auth.Context, ap *db.AgentProfile, ts *db.TradingSystem) (*agentscanner.TradingSystem, error) {
	client := agentscanner.CreateClient(ap.SslCertRef, ap.SslKeyRef, "ca.crt")
	if client == nil {
		return nil, req.NewServerError("Cannot create client for agent: %v", ap.Id)
	}

	var ats agentscanner.TradingSystem
	params := &TradingSystemReloadRequest{
		Name: ts.ExternalRef,
	}
	err := req.DoPost(client, ap.RemoteUrl+agentscanner.UrlReloadTradingSystem, params, &ats, "")
	if err != nil {
		c.Log.Error("ReloadTradingSystem: Agent raised an error", "id", ap.Id, "error", err.Error())
		return nil, req.NewServiceUnavailableError("Agent raised an error : " + err.Error())
	}

	return &ats, nil
}

//=============================================================================
