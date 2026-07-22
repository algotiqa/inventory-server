//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package main

import (
	"log/slog"

	"github.com/algotiqa/core/boot"
	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/algotiqa/inventory-server/pkg/core/messaging/system"
	"github.com/algotiqa/inventory-server/pkg/core/process"
	"github.com/algotiqa/inventory-server/pkg/service"
)

//=============================================================================

const component = "inventory-server"
var   version   = "dev"

//=============================================================================

func main() {
	cfg := &app.Config{}
	boot.ReadConfig(component, cfg)
	logger := boot.InitLogger(component, version, &cfg.Application)
	engine := boot.InitEngine(logger, &cfg.Application)
	initClients()
	dbms.InitDatabase(&cfg.Database)
	msg.InitMessaging(&cfg.Messaging)
	service.Init(engine, cfg, logger)
	process.Init(cfg)
	system.InitMessageListener()
	boot.RunHttpServer(engine, &cfg.Application)
}

//=============================================================================

func initClients() {
	slog.Info("Initializing clients...")
	req.AddDefaultClient("ca.crt", "server.crt", "server.key")
}

//=============================================================================
