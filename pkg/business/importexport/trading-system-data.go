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
	"github.com/algotiqa/inventory-server/pkg/db"
)

//=============================================================================
//===
//=== Data collection structures
//===
//=============================================================================

type TradingSystemsData struct {
	Ids             []uint
	TradingSystems  []db.TradingSystem
	DataProducts    map[uint]*db.DataProductFull
	BrokerProducts  map[uint]*db.BrokerProductFull
	TradingSessions map[uint]*db.TradingSession
	AgentProfiles   map[uint]*db.AgentProfile
}

//=============================================================================

func NewTradingSystemsData(tss []db.TradingSystem, dps []db.DataProductFull, bps []db.BrokerProductFull,
	sess []db.TradingSession, aps []db.AgentProfile) *TradingSystemsData {
	tsd := &TradingSystemsData{
		TradingSystems : tss,
		DataProducts   : make(map[uint]*db.DataProductFull),
		BrokerProducts : make(map[uint]*db.BrokerProductFull),
		TradingSessions: make(map[uint]*db.TradingSession),
		AgentProfiles  : make(map[uint]*db.AgentProfile),
	}

	products := buildDataProductMap(dps)
	brokers  := buildDataBrokerMap(bps)
	sessions := buildSessionMap(sess)
	profiles := buildProfileMap(aps)

	for _, ts := range tss {
		p,okp := products[ts.DataProductId   ]
		b,okb := brokers [ts.BrokerProductId ]
		s,oks := sessions[ts.TradingSessionId]

		if !okp || !okb || !oks {
			return nil
		}

		tsd.DataProducts   [p.Id] = p
		tsd.BrokerProducts [b.Id] = b
		tsd.TradingSessions[s.Id] = s

		if ts.AgentProfileId != nil {
			a,oka := profiles[*ts.AgentProfileId]
			if !oka {
				return nil
			}

			tsd.AgentProfiles[a.Id] = a
		}

		ss,ok := sessions[p.SessionId]
		if !ok {
			return nil
		}
		tsd.TradingSessions[ss.Id] = ss
	}

	return tsd
}

//=============================================================================

func buildDataProductMap(list []db.DataProductFull) map[uint]*db.DataProductFull {
	res := map[uint]*db.DataProductFull{}
	for _, x := range list {
		res[x.Id] = &x
	}

	return res
}

//=============================================================================

func buildDataBrokerMap(list []db.BrokerProductFull) map[uint]*db.BrokerProductFull {
	res := map[uint]*db.BrokerProductFull{}

	for _, x := range list {
		res[x.Id] = &x
	}

	return res
}

//=============================================================================

func buildSessionMap(list []db.TradingSession) map[uint]*db.TradingSession {
	res := map[uint]*db.TradingSession{}

	for _, x := range list {
		res[x.Id] = &x
	}

	return res
}

//=============================================================================

func buildProfileMap(list []db.AgentProfile) map[uint]*db.AgentProfile {
	res := map[uint]*db.AgentProfile{}

	for _, x := range list {
		res[x.Id] = &x
	}

	return res
}

//=============================================================================
