package main

import (
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/server"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/zeebe"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
)

func main() {
	finish := make(chan struct{})

	context := config.Bootstrap(config.AppConfig)
	go zeebe.InitAndStart(context)
	go server.InitAndStart(context)

	<- finish
}
