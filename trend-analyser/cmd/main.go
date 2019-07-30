package main

import (
	"github.com/domahidizoltan/playground-workflow-engine/trend-analyser/pkg/zeebe"
	"github.com/domahidizoltan/playground-workflow-engine/trend-analyser/internal/config"
)

func main() {
	context := config.Bootstrap(config.AppConfig)
	zeebe.InitAndStart(context)
}
