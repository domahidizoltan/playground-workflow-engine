package zeebe

import (
	"log"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/position"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
)

const buyJobTypeName = "buy-stock"
const sellJobTypeName = "sell-stock"

type zeebeClient struct {
	client zbc.ZBClient
	positionService position.PositionService
}

var zClient zeebeClient

func InitAndStart(context config.Context) {
	conf := config.AppConfig.Zeebe
	client, err := zbc.NewZBClient(conf.Host + ":" + conf.Port)
	if err != nil {
		panic(err)
	}

	log.Println("Zeebe client is running on port " + conf.Port)

	zClient = zeebeClient {
		client: client,
		positionService: context.PositionService,
	}

	go startWorker(buyJobTypeName, handleBuyJob)
	go startWorker(sellJobTypeName, handleSellJob)
}

func startWorker(jobTypeName string, handleJob func(worker.JobClient, entities.Job)) {
	workerName := jobTypeName + ":" + config.AppConfig.Http.Server.Port
	jobWorker := zClient.client.NewJobWorker().
		JobType(jobTypeName).
		Handler(handleJob).
		Name(workerName).
		Open()
	defer jobWorker.Close()

	log.Println("Worker " + workerName + " started")
	jobWorker.AwaitClose()
}

func handleBuyJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println("Failed to fetch variables")
		failJob(client, job)
		return
	}

	username := variables["username"].(string)
	symbol := variables["symbol"].(string)
	if pos, e := zClient.positionService.Buy(username, symbol); e != nil {
		log.Printf("Could not buy position %v: %s", pos, e.Error())
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		log.Println("Failed to set output fields")
		failJob(client, job)
		return
	}
	
	log.Println("Complete job", jobKey, "of type", job.Type)
	request.Send()
}

func handleSellJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println("Failed to fetch variables")
		failJob(client, job)
		return
	}

	username := variables["username"].(string)
	symbol := variables["symbol"].(string)
	if e := zClient.positionService.Sell(username, symbol); e != nil {
		log.Printf("Could not sell position of %s: %s", symbol, e.Error())
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		log.Println("Failed to set output fields")
		failJob(client, job)
		return
	}
	
	log.Println("Complete job", jobKey, "of type", job.Type)
	request.Send()
}


func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
