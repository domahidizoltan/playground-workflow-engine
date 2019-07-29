package main

import (
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/pkg/server"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/pkg/zeebe"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/config"
)

func main() {
	finish := make(chan struct{})

	context := config.Bootstrap(config.AppConfig)
	go zeebe.InitAndStart(context)
	go server.InitAndStart(context)

	<- finish
}
