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
	"strconv"
	"strings"

	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

const (
	MatchSep = " vs "
	MatchLF  = " | "
)

//=============================================================================

func CreateImportOverview(tx *gorm.DB, user string, pack *InMemoryPackage) (*ImportOverviewResponse, error) {
	plan := &ImportOverviewResponse{
		TradingSystems: createTradingSystemItems(pack.Data.TradingSystems),
	}

	//--- Add plan for data products

	dpRefs,err := matchDataProducts(tx, user, pack.Data.DataProducts)
	if err != nil {
		return nil, err
	}

	plan.ReferencedItems = append(plan.ReferencedItems, dpRefs...)

	//--- Add plan for broker products

	bpRefs,err := matchBrokerProducts(tx, user, pack.Data.BrokerProducts)
	if err != nil {
		return nil, err
	}

	plan.ReferencedItems = append(plan.ReferencedItems, bpRefs...)

	//--- Add plan for agent profiles

	apRefs,err := matchAgentProfiles(tx, user, pack.Data.AgentProfiles)
	if err != nil {
		return nil, err
	}

	plan.ReferencedItems = append(plan.ReferencedItems, apRefs...)

	//--- Add plan for trading sessions

	tsRefs,err := matchTradingSessions(tx, user, pack.Data.TradingSessions)
	if err != nil {
		return nil, err
	}

	plan.ReferencedItems = append(plan.ReferencedItems, tsRefs...)

	return plan, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func createTradingSystemItems(list []*TradingSystem) []*TradingSystemItem {
	var res []*TradingSystemItem

	for _, ts := range list {
		tsi := &TradingSystemItem{
			Id       : ts.Id,
			Name     : ts.Name,
			Timeframe: ts.Timeframe,
		}
		res = append(res, tsi)
	}

	return res
}

//=============================================================================
//=== Data products
//=============================================================================

func matchDataProducts(tx *gorm.DB, user string, list []*DataProduct) ([]*ReferencedItem, error) {
	filter := make(map[string]any)
	filter["username"] = user

	dps,err := db.GetDataProductsFull(tx, filter, 0, 5000)
	if err != nil {
		return nil, err
	}

	var res []*ReferencedItem

	for _, dp := range list {
		planItem := &ReferencedItem{
			Id          : dp.Id,
			Symbol      : dp.Symbol,
			Name        : dp.Name,
			SystemCode  : dp.SystemCode,
			ExchangeCode: dp.ExchangeCode,
			ItemType    : ReferencedItemTypeData,
		}

		options := matchDataProduct(dp, dps)

		if options != nil {
			planItem.Status  = RIStatusExisting
			planItem.Options = options
			planItem.MappedTo= options[0].Id
		} else {
			planItem.Status = RIStatusNoConnection
		}

		res = append(res, planItem)
	}

	return res, nil
}

//=============================================================================

func matchDataProduct(dp *DataProduct, list *[]db.DataProductFull) []*ReferencedOption {
	var res []*ReferencedOption

	for _, dbDp := range *list {
		if matchesDataProduct(dp, &dbDp) {
			pio := &ReferencedOption{
				Id         : dbDp.Id,
				Name       : dbDp.ConnectionCode,
				MatchNotes : createDataProductMatchNotes(dp, &dbDp),
				dataProduct: &dbDp.DataProduct,
			}
			res = append(res, pio)
		}
	}

	return res
}

//=============================================================================

func matchesDataProduct(dp *DataProduct, dbDp *db.DataProductFull) bool {
	return (dbDp.Symbol == dp.Symbol) && (dbDp.ExchangeId == dp.ExchangeId)
}

//=============================================================================

func createDataProductMatchNotes(dp *DataProduct, dbDp *db.DataProductFull) string {
	var sb strings.Builder

	if dp.SystemCode != dbDp.SystemCode {
		addMatchDiff(&sb, "systemCode", dp.SystemCode, dbDp.SystemCode)
	}
	if dp.MarketType != dbDp.MarketType {
		addMatchDiff(&sb, "marketType", dp.MarketType, dbDp.MarketType)
	}
	if dp.ProductType != dbDp.ProductType {
		addMatchDiff(&sb, "productType", dp.ProductType, dbDp.ProductType)
	}
	if dp.Months != dbDp.Months {
		addMatchDiff(&sb, "months", dp.Months, dbDp.Months)
	}
	if dp.RolloverTrigger != dbDp.RolloverTrigger {
		addMatchDiff(&sb, "rolloverTrigger", string(dp.RolloverTrigger), string(dbDp.RolloverTrigger))
	}

	return sb.String()
}

//=============================================================================
//=== Broker products
//=============================================================================

func matchBrokerProducts(tx *gorm.DB, user string, list []*BrokerProduct) ([]*ReferencedItem, error) {
	filter := make(map[string]any)
	filter["username"] = user

	bps,err := db.GetBrokerProductsFull(tx, filter, 0, 5000)
	if err != nil {
		return nil, err
	}

	var res []*ReferencedItem

	for _, bp := range list {
		planItem := &ReferencedItem{
			Id          : bp.Id,
			Symbol      : bp.Symbol,
			Name        : bp.Name,
			SystemCode  : bp.SystemCode,
			ExchangeCode: bp.ExchangeCode,
			ItemType    : ReferencedItemTypeBroker,
		}

		options := matchBrokerProduct(bp, bps)

		if options != nil {
			planItem.Status  = RIStatusExisting
			planItem.Options = options
			planItem.MappedTo= options[0].Id
		} else {
			planItem.Status = RIStatusNoConnection
		}

		res = append(res, planItem)
	}

	return res, nil
}

//=============================================================================

func matchBrokerProduct(bp *BrokerProduct, list *[]db.BrokerProductFull) []*ReferencedOption {
	var res []*ReferencedOption

	for _, dbBp := range *list {
		if matchesBrokerProduct(bp, &dbBp) {
			pio := &ReferencedOption{
				Id           : dbBp.Id,
				Name         : dbBp.ConnectionCode,
				MatchNotes   : createBrokerProductMatchNotes(bp, &dbBp),
				brokerProduct: &dbBp.BrokerProduct,
			}
			res = append(res, pio)
		}
	}

	return res
}

//=============================================================================

func matchesBrokerProduct(bp *BrokerProduct, dbBp *db.BrokerProductFull) bool {
	return (dbBp.Symbol == bp.Symbol) && (dbBp.ExchangeId == bp.ExchangeId)
}

//=============================================================================

func createBrokerProductMatchNotes(bp *BrokerProduct, dbBp *db.BrokerProductFull) string {
	var sb strings.Builder

	if bp.SystemCode != dbBp.SystemCode {
		addMatchDiff(&sb, "systemCode", bp.SystemCode, dbBp.SystemCode)
	}
	if bp.MarketType != dbBp.MarketType {
		addMatchDiff(&sb, "marketType", bp.MarketType, dbBp.MarketType)
	}
	if bp.ProductType != dbBp.ProductType {
		addMatchDiff(&sb, "productType", bp.ProductType, dbBp.ProductType)
	}
	if bp.PointValue != dbBp.PointValue {
		addMatchDiff(&sb, "pointValue", s64(bp.PointValue), s64(dbBp.PointValue))
	}
	if bp.CostPerOperation != dbBp.CostPerOperation {
		addMatchDiff(&sb, "costPerOperation", s64(bp.CostPerOperation), s64(dbBp.CostPerOperation))
	}
	if bp.MarginValue != dbBp.MarginValue {
		addMatchDiff(&sb, "marginValue", s64(bp.MarginValue), s64(dbBp.MarginValue))
	}
	if bp.Increment != dbBp.Increment {
		addMatchDiff(&sb, "increment", s64(bp.Increment), s64(dbBp.Increment))
	}

	return sb.String()
}

//=============================================================================
//=== Agent profiles
//=============================================================================

func matchAgentProfiles(tx *gorm.DB, user string, list []*AgentProfile) ([]*ReferencedItem, error) {
	filter := make(map[string]any)
	filter["username"] = user

	aps,err := db.GetAgentProfiles(tx, filter, 0, 5000)
	if err != nil {
		return nil, err
	}

	var res []*ReferencedItem

	for _, ap := range list {
		planItem := &ReferencedItem{
			Id          : ap.Id,
			Symbol      : "",
			Name        : ap.Name,
			SystemCode  : "",
			ExchangeCode: "",
			ItemType    : ReferencedItemTypeProfile,
		}

		options := matchAgentProfile(ap, aps)

		if options != nil {
			planItem.Status  = RIStatusExisting
			planItem.Options = options
			planItem.MappedTo= options[0].Id

		} else {
			planItem.Status = RIStatusNew
		}

		res = append(res, planItem)
	}

	return res, nil
}

//=============================================================================

func matchAgentProfile(ap *AgentProfile, list *[]db.AgentProfile) []*ReferencedOption {
	var res []*ReferencedOption

	for _, dbAp := range *list {
		if matchesAgentProfile(ap, &dbAp) {
			pio := &ReferencedOption{
				Id          : dbAp.Id,
				Name        : dbAp.Name,
				agentProfile: &dbAp,
			}
			res = append(res, pio)
		}
	}

	return res
}

//=============================================================================

func matchesAgentProfile(ap *AgentProfile, dbAp *db.AgentProfile) bool {
	return dbAp.Name == ap.Name
}

//=============================================================================
//=== Trading sessions
//=============================================================================

func matchTradingSessions(tx *gorm.DB, user string, list []*TradingSession) ([]*ReferencedItem, error) {
	tss,err := db.GetTradingSessions(tx, user)
	if err != nil {
		return nil, err
	}

	var res []*ReferencedItem

	for _, s := range list {
		planItem := &ReferencedItem{
			Id          : s.Id,
			Symbol      : "",
			Name        : s.Name,
			SystemCode  : "",
			ExchangeCode: "",
			ItemType    : ReferencedItemTypeSession,
		}

		options := matchTradingSession(s, tss)

		if options != nil {
			planItem.Status  = RIStatusExisting
			planItem.Options = options
			planItem.MappedTo= options[0].Id

		} else {
			planItem.Status        = RIStatusNew
			planItem.sessionConfig = s.Session
		}

		res = append(res, planItem)
	}

	return res, nil
}

//=============================================================================

func matchTradingSession(s *TradingSession, list *[]db.TradingSession) []*ReferencedOption {
	var res []*ReferencedOption

	for _, dbS := range *list {
		if matchesTradingSession(s, &dbS) {
			pio := &ReferencedOption{
				Id            : dbS.Id,
				Name          : dbS.Name,
				tradingSession: &dbS,
			}
			res = append(res, pio)
		}
	}

	return res
}

//=============================================================================

func matchesTradingSession(s *TradingSession, dbS *db.TradingSession) bool {
	return dbS.Session == s.Session
}

//=============================================================================
//=== Other functions
//=============================================================================

func addMatchDiff(sb *strings.Builder, name string, value1, value2 string) {
	sb.WriteString(name +": \""+ value1 +"\""+MatchSep+"\""+ value2 +"\""+MatchLF)
}

//=============================================================================

func s32(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

//=============================================================================

func s64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

//=============================================================================
