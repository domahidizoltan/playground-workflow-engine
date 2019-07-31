package zeebe

import (
	"log"

	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/position"
	sClient "github.com/domahidizoltan/playground-workflow-engine/trading-service/pkg/client"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
)

const buyJobTypeName = "buy-position"
const sellJobTypeName = "sell-position"

type zeebeClient struct {
	client          zbc.ZBClient
	positionService position.PositionService
}

var zClient zeebeClient

func InitAndStart(context config.Context) {
	conf := config.AppConfig.Zeebe
	client, err := zbc.NewZBClient(conf.Host + ":" + conf.Port)
	if err != nil {
		panic(err)
	}

	zClient = zeebeClient{
		client:          client,
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
	data := toStockData(variables["stockdata"])
	if pos, e := zClient.positionService.Buy(username, data); e != nil {
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

	log.Printf("%v", variables)
	username := variables["username"].(string)
	data := toStockData(variables["stockdata"])
	if e := zClient.positionService.Sell(username, data); e != nil {
		log.Printf("Could not sell position of %s: %s", data.Symbol, e.Error())
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

func toStockData(stockdata interface{}) sClient.StockData {
	sd := stockdata.(map[string]interface{})
	symbol, _ := sd["symbol"].(string)
	price, _ := sd["price"].(float64)

	return sClient.StockData{
		Symbol: symbol,
		Price:  float32(price),
	}
}
