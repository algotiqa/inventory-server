//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
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

	ctrl := auth.NewOidcController(cfg.Authentication.Authority, req.GetDefaultClient(), logger, cfg)

	//--- Inventory

	router.GET   ("/api/inventory/v1/currencies",          ctrl.Secure(getCurrencies, roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/exchanges",           ctrl.Secure(getExchanges,  roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/data-products",       ctrl.Secure(getDataProducts,    roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/data-products",       ctrl.Secure(addDataProduct,     roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/data-products/:id",   ctrl.Secure(getDataProductById, roles.Admin_User_Service))
	router.PUT   ("/api/inventory/v1/data-products/:id",   ctrl.Secure(updateDataProduct,  roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/data-products/:id",   ctrl.Secure(deleteDataProduct,  roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/broker-products",     ctrl.Secure(getBrokerProducts,    roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/broker-products",     ctrl.Secure(addBrokerProduct,     roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/broker-products/:id", ctrl.Secure(getBrokerProductById, roles.Admin_User_Service))
	router.PUT   ("/api/inventory/v1/broker-products/:id", ctrl.Secure(updateBrokerProduct,  roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/broker-products/:id", ctrl.Secure(deleteBrokerProduct,  roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/trading-systems",                   ctrl.Secure(getTradingSystems,     roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/trading-systems",                   ctrl.Secure(addTradingSystem,      roles.Admin_User_Service))
	router.PUT   ("/api/inventory/v1/trading-systems/:id",               ctrl.Secure(updateTradingSystem,   roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/trading-systems/:id",               ctrl.Secure(deleteTradingSystem,   roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/trading-systems/:id/finalize",      ctrl.Secure(finalizeTradingSystem, roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/trading-systems/:id/reload-trades", ctrl.Secure(reloadTradesFromAgent, roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/trading-systems/export",            ctrl.Secure(exportTradingSystems,  roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/trading-systems/import/overview",   ctrl.Secure(createImportOverview,  roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/trading-systems/import/execute",    ctrl.Secure(executeImportPlan,     roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/trading-sessions",                  ctrl.Secure(getTradingSessions, roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/agent-profiles",                    ctrl.Secure(getAgentProfiles,    roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/agent-profiles",                    ctrl.Secure(addAgentProfile,     roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/agent-profiles/:id",                ctrl.Secure(getAgentProfileById, roles.Admin_User_Service))
	router.PUT   ("/api/inventory/v1/agent-profiles/:id",                ctrl.Secure(updateAgentProfile,  roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/agent-profiles/:id",                ctrl.Secure(deleteAgentProfile,  roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/agent-profiles/:id/external-refs",  ctrl.Secure(getExternalRefs,     roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/agent-profiles/:id/package",        ctrl.Secure(getAgentPackage,     roles.Admin_User_Service))

	//--- Administration

	router.GET   ("/api/inventory/v1/connections",     ctrl.Secure(getConnections,    roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/connections/:id", ctrl.Secure(getConnectionById, roles.Admin_User_Service))
	router.POST  ("/api/inventory/v1/connections",     ctrl.Secure(addConnection,     roles.Admin_User_Service))
	router.PUT   ("/api/inventory/v1/connections/:id", ctrl.Secure(updateConnection,  roles.Admin_User_Service))
	router.DELETE("/api/inventory/v1/connections/:id", ctrl.Secure(deleteConnection,  roles.Admin_User_Service))
}

//=============================================================================
