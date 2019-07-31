package zeebe

import (
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/stockprice"
	spClient "github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/client"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
	"log"
)

const jobTypeName = "fetch-stock-price"

type zeebeClient struct {
	stockpriceClient spClient.StockPriceClient
	stockdataService stockprice.StockDataService
}

var zClient zeebeClient

func InitAndStart(context config.Context) {
	conf := config.AppConfig.Zeebe
	client, err := zbc.NewZBClient(conf.Host + ":" + conf.Port)
	if err != nil {
		panic(err)
	}

	zClient = zeebeClient {
		stockdataService: context.StockDataService,
		stockpriceClient: spClient.NewStockPriceClient(),
	}

	workerName := jobTypeName + ":" + config.AppConfig.Http.Server.Port
	jobWorker := client.NewJobWorker().
		JobType(jobTypeName).
		Handler(handleJob).
		Name(workerName).
		Open()
	defer jobWorker.Close()

	log.Println("Worker " + workerName + " started")
	jobWorker.AwaitClose()
}

func handleJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println("Failed to fetch variables")
		failJob(client, job)
		return
	}

	variables["stockdata"] = fetchAndGetLastValuesOfSymbol(variables["symbol"].(string))
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

func fetchAndGetLastValuesOfSymbol(symbol string) stockprice.StockData {
	log.Println("Processing symbol:", symbol)
	stockdata := zClient.stockpriceClient.FetchStockData(symbol)
	zClient.stockdataService.Save(stockdata)
	return stockdata
}
