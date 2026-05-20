//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

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
