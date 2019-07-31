package zeebe

import (
	"log"
	"github.com/domahidizoltan/playground-workflow-engine/trend-analyser/internal/config"
	ta "github.com/domahidizoltan/playground-workflow-engine/trend-analyser/internal/trendanalyser"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
)

const jobTypeName = "trend-analyser"

type zeebeClient struct {
	trendAnalyser ta.TrendAnalyser
}

var zClient zeebeClient

func InitAndStart(context config.Context) {
	conf := config.AppConfig.Zeebe
	client, err := zbc.NewZBClient(conf.Host + ":" + conf.Port)
	if err != nil {
		panic(err)
	}

	zClient = zeebeClient{
		trendAnalyser: context.TrendAnalyser,
	}

	jobWorker := client.NewJobWorker().
		JobType(jobTypeName).
		Handler(handleJob).
		Name(jobTypeName).
		Open()
	defer jobWorker.Close()

	log.Println("Worker " + jobTypeName + " started")
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

	stockdata := variables["stockdata"]
	data := toStockData(stockdata)
	log.Printf("Received %v", data)
	zClient.trendAnalyser.SetLastPrice(data.Symbol, data.Price)

	variables["trend"] = nil
	if zClient.trendAnalyser.IsReadyToAnalyse(data.Symbol) {
		variables["trend"] = zClient.trendAnalyser.Analyse(data.Symbol)
		log.Printf("Sending trend %v for %s", variables["trend"], data.Symbol)
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)

	if err != nil {
		log.Println("Failed to set output fields")
		failJob(client, job)
		return
	}
	request.Send()

	log.Println("Complete job", jobKey, "of type", job.Type)
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}

func toStockData(stockdata interface{}) ta.StockData {
	sd := stockdata.(map[string]interface{})
	symbol, _ := sd["symbol"].(string)
	price, _ := sd["price"].(float64)

	return ta.StockData{
		Symbol: symbol,
		Price:  float32(price),
	}
}
