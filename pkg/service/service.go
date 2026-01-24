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

package service

import (
	"log/slog"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/auth/roles"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/gin-gonic/gin"
)

//=============================================================================

func Init(router *gin.Engine, cfg *app.Config, logger *slog.Logger) {

	ctrl := auth.NewOidcController(cfg.Authentication.Authority, req.GetClient("bf"), logger, cfg)

	//--- Inventory

	router.GET("/api/inventory/v1/currencies", ctrl.Secure(getCurrencies, roles.Admin_User_Service))
	router.GET("/api/inventory/v1/exchanges", ctrl.Secure(getExchanges, roles.Admin_User_Service))

	router.GET("/api/inventory/v1/data-products", ctrl.Secure(getDataProducts, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/data-products", ctrl.Secure(addDataProduct, roles.Admin_User_Service))
	router.GET("/api/inventory/v1/data-products/:id", ctrl.Secure(getDataProductById, roles.Admin_User_Service))
	router.PUT("/api/inventory/v1/data-products/:id", ctrl.Secure(updateDataProduct, roles.Admin_User_Service))

	router.GET("/api/inventory/v1/broker-products", ctrl.Secure(getBrokerProducts, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/broker-products", ctrl.Secure(addBrokerProduct, roles.Admin_User_Service))
	router.GET("/api/inventory/v1/broker-products/:id", ctrl.Secure(getBrokerProductById, roles.Admin_User_Service))
	router.PUT("/api/inventory/v1/broker-products/:id", ctrl.Secure(updateBrokerProduct, roles.Admin_User_Service))

	router.GET("/api/inventory/v1/trading-systems", ctrl.Secure(getTradingSystems, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/trading-systems", ctrl.Secure(addTradingSystem, roles.Admin_User_Service))
	router.PUT("/api/inventory/v1/trading-systems/:id", ctrl.Secure(updateTradingSystem, roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/trading-systems/:id", ctrl.Secure(deleteTradingSystem, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/trading-systems/:id/finalize", ctrl.Secure(finalizeTradingSystem, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/trading-systems/:id/reload-trades", ctrl.Secure(reloadTradesFromAgent, roles.Admin_User_Service))

	router.GET("/api/inventory/v1/trading-sessions", ctrl.Secure(getTradingSessions, roles.Admin_User_Service))
	router.GET("/api/inventory/v1/agent-profiles", ctrl.Secure(getAgentProfiles, roles.Admin_User_Service))

	//--- Administration

	router.GET("/api/inventory/v1/connections", ctrl.Secure(getConnections, roles.Admin_User_Service))
	router.GET("/api/inventory/v1/connections/:id", ctrl.Secure(getConnectionById, roles.Admin_User_Service))
	router.POST("/api/inventory/v1/connections", ctrl.Secure(addConnection, roles.Admin_User_Service))
	router.PUT("/api/inventory/v1/connections/:id", ctrl.Secure(updateConnection, roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/connections/:id", ctrl.Secure(deleteConnection, roles.Admin_User_Service))
}

//=============================================================================
