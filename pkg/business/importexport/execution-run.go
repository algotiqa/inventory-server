//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package importexport

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/core/messaging"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================
//===
//=== ExecutionCache
//===
//=============================================================================

type ExecutionCache struct {
	dataProducts    map[uint]*DataProduct
	brokerProducts  map[uint]*BrokerProduct
	tradingSessions map[uint]*TradingSession
	agentProfiles   map[uint]*AgentProfile
	tradingSystems  map[uint]*TradingSystem
}

//=============================================================================

func NewExecutionCache(ed *ExportedData) *ExecutionCache {
	ec := &ExecutionCache{
		dataProducts   : make(map[uint]*DataProduct),
		brokerProducts : make(map[uint]*BrokerProduct),
		tradingSessions: make(map[uint]*TradingSession),
		agentProfiles  : make(map[uint]*AgentProfile),
		tradingSystems : make(map[uint]*TradingSystem),
	}

	for _, dp := range ed.DataProducts {
		ec.dataProducts[dp.Id] = dp
	}

	for _, bp := range ed.BrokerProducts {
		ec.brokerProducts[bp.Id] = bp
	}

	for _, ts := range ed.TradingSessions {
		ec.tradingSessions[ts.Id] = ts
	}

	for _, ts := range ed.TradingSystems {
		ec.tradingSystems[ts.Id] = ts
	}

	for _, ap := range ed.AgentProfiles {
		ec.agentProfiles[ap.Id] = ap
	}

	return ec
}

//=============================================================================
//===
//=== Plan execution
//===
//=============================================================================

func ExecutePlan(tx *gorm.DB, c *auth.Context, pack *InMemoryPackage, plan *ImportPlan, res *ImportOverviewResponse) error {
	cache         := NewExecutionCache(pack.Data)
	username      := c.Session.Username
	addedProfiles := make(map[uint]*db.AgentProfile)
	addedSessions := make(map[uint]*db.TradingSession)

	exchanges,errE := db.GetExchangesAsMap(tx)
	if errE != nil {
		return req.NewServerErrorByError(errE)
	}

	currencies,errC := db.GetCurrenciesAsMap(tx)
	if errC != nil {
		return req.NewServerErrorByError(errC)
	}

	for _, sts := range plan.TradingSystems {
		ts,ok := cache.tradingSystems[sts.Id]
		if !ok {
			return req.NewBadRequestError("trading system not found in package: %v", sts.Id)
		}

		dp,errD := findDataProduct(ts.DataProductId, plan.ReferencedItems, res.ReferencedItems)
		if errD != nil {
			return errD
		}

		bp,errB := findBrokerProduct(ts.BrokerProductId, plan.ReferencedItems, res.ReferencedItems)
		if errB != nil {
			return errB
		}

		ap,errP := findAgentProfile(tx, username, ts.AgentProfileId, plan.ReferencedItems, res.ReferencedItems, addedProfiles)
		if errP != nil {
			return errP
		}

		se,errS := findTradingSession(tx, username, ts.TradingSessionId, plan.ReferencedItems, res.ReferencedItems, addedSessions)
		if errS != nil {
			return errS
		}

		dbTs := convertTradingSystem(ts)
		//--- Update name with the new one provided by the user
		dbTs.Name             = sts.Name
		dbTs.Username         = username
		dbTs.DataProductId    = dp.Id
		dbTs.BrokerProductId  = bp.Id
		dbTs.TradingSessionId = se.Id

		if ap != nil {
			dbTs.AgentProfileId = &ap.Id
		}

		err := db.AddTradingSystem(tx, dbTs)
		if err != nil {
			c.Log.Error("ExecutePlan: Could not add a new trading system", "error", err.Error())
			return err
		}

		ex,okE := exchanges[bp.ExchangeId]
		if !okE {
			return req.NewServerError("exchange not found in database: %v", bp.ExchangeId)
		}

		cu,okC := currencies[ex.CurrencyId]
		if !okC {
			return req.NewServerError("currency not found in database: %v", ex.CurrencyId)
		}

		tsm := messaging.TradingSystemMessage{
			TradingSystem : dbTs,
			DataProduct   : dp,
			BrokerProduct : bp,
			Currency      : cu,
			TradingSession: se,
			AgentProfile  : ap,
			Exchange      : ex,
			PortfolioPack : ts.portfolioData,
			StoragePack   : ts.storageData,
		}

		err = msg.SendMessage(msg.ExInventory, msg.SourceTradingSystem, msg.TypeCreate, &tsm, tx)
		if err != nil {
			c.Log.Error("ExecutePlan: Could not publish the update message for TS", "error", err.Error(), "id", ts.Id)
			return err
		}
	}

	return nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func findDataProduct(dpIdInPackage uint, planRefs []*SelectedReference, ovRefs []*ReferencedItem) (*db.DataProduct, error) {
	for _, ref := range planRefs {
		if dpIdInPackage == ref.Id && ref.ItemType == ReferencedItemTypeData {
			for _, ovRef := range ovRefs {
				if dpIdInPackage == ovRef.Id && ovRef.ItemType == ReferencedItemTypeData {
					for _, opt := range ovRef.Options {
						if ref.MappedTo == opt.Id {
							return opt.dataProduct, nil
						}
					}
				}
			}
		}
	}

	return nil, req.NewBadRequestError("Can't map data product in package: %v", dpIdInPackage)
}

//=============================================================================

func findBrokerProduct(bpIdInPackage uint, planRefs []*SelectedReference, ovRefs []*ReferencedItem) (*db.BrokerProduct, error) {
	for _, ref := range planRefs {
		if bpIdInPackage == ref.Id && ref.ItemType == ReferencedItemTypeBroker {
			for _, ovRef := range ovRefs {
				if bpIdInPackage == ovRef.Id && ovRef.ItemType == ReferencedItemTypeBroker {
					for _, opt := range ovRef.Options {
						if ref.MappedTo == opt.Id {
							return opt.brokerProduct, nil
						}
					}
				}
			}
		}
	}

	return nil, req.NewBadRequestError("Can't map broker product in package: %v", bpIdInPackage)
}

//=============================================================================

func findAgentProfile(tx *gorm.DB, username string, apIdInPackage *uint, planRefs []*SelectedReference,
					  ovRefs []*ReferencedItem, added map[uint]*db.AgentProfile) (*db.AgentProfile, error) {
	if apIdInPackage == nil {
		return nil, nil
	}

	for _, ovRef := range ovRefs {
		if *apIdInPackage == ovRef.Id && ovRef.ItemType == ReferencedItemTypeProfile {
			if ovRef.Status == RIStatusNew {
				//--- Profile must be added

				ap,ok := added[*apIdInPackage]
				if !ok {
					ap = &db.AgentProfile{
						Username    : username,
						Name        : ovRef.Name,
						RemoteUrl   : "unknown",
						ScanInterval: 0,
					}

					err := db.AddAgentProfile(tx,ap)
					if err != nil {
						return nil,req.NewServerErrorByError(err)
					}
					added[*apIdInPackage] = ap
				}
				return ap, nil
			}

			//--- There must be an existing profile

			for _, ref := range planRefs {
				if *apIdInPackage == ref.Id && ref.ItemType == ReferencedItemTypeProfile {
					for _, opt := range ovRef.Options {
						if ref.MappedTo == opt.Id {
							return opt.agentProfile, nil
						}
					}
				}
			}
		}
	}

	return nil, req.NewBadRequestError("Can't map agent profile in package: %v", *apIdInPackage)
}

//=============================================================================

func findTradingSession(tx *gorm.DB, username string, seIdInPackage uint, planRefs []*SelectedReference,
						ovRefs []*ReferencedItem, added map[uint]*db.TradingSession) (*db.TradingSession, error) {
	for _, ovRef := range ovRefs {
		if seIdInPackage == ovRef.Id && ovRef.ItemType == ReferencedItemTypeSession {
			if ovRef.Status == RIStatusNew {
				//--- Session must be added

				se,ok := added[seIdInPackage]
				if !ok {
					se = &db.TradingSession{
						Username: username,
						Name    : ovRef.Name,
						Session : ovRef.sessionConfig,
					}

					err := db.AddTradingSession(tx,se)
					if err != nil {
						return nil,req.NewServerErrorByError(err)
					}
					added[seIdInPackage] = se
				}
				return se, nil
			}

			//--- There must be an existing profile

			for _, ref := range planRefs {
				if seIdInPackage == ref.Id && ref.ItemType == ReferencedItemTypeSession {
					for _, opt := range ovRef.Options {
						if ref.MappedTo == opt.Id {
							return opt.tradingSession, nil
						}
					}
				}
			}
		}
	}

	return nil, req.NewBadRequestError("Can't map trading session in package: %v", seIdInPackage)
}

//=============================================================================

func convertTradingSystem(ts *TradingSystem) *db.TradingSystem {
	var dbTs db.TradingSystem
	dbTs.Timeframe        = ts.Timeframe
	dbTs.StrategyType     = ts.StrategyType
	dbTs.Overnight        = ts.Overnight
	dbTs.Tags             = ts.Tags
	dbTs.ExternalRef      = ts.ExternalRef
	dbTs.InSampleFrom     = ts.InSampleFrom
	dbTs.InSampleTo       = ts.InSampleTo
	dbTs.EngineCode       = ts.EngineCode
	dbTs.Finalized = ts.Finalized

	return &dbTs
}

//=============================================================================
