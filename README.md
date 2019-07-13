# playground-workflow-engine
WIP: A playground project to learn Zeebe and Go

- stockprice-mockserver:  
`http://localhost:9000/api/company/real-time-price/{SYMBOL}[?trend={number}]`

- stockprice-service:  
`http://localhost:9100/api/stockdatas/{SYMBOL}[?offset={number}&limit={number}]`



**Zeebe order-process tutorial**
https://docs.zeebe.io/getting-started/README.html  

`docker-compose up`

- enter container and deploy workers

`docker exec -it zeebe sh`

```
zbctl create worker payment-service --handler cat &
zbctl create worker inventory-service --handler cat &
zbctl create worker shipment-service --handler cat &
```

- Zeebe simple monitor: [http://localhost:8000]()

- Zeebe monitor db: [http://localhost:81]()

- Deploy `order-process.bpmn` or create instance by using:  
`zbctl create instance order-process --variables '{"orderId": 12345}'`


**First draft (stockprice-service)**
- exec into zeebe container  
```zbctl create worker analyze-trend --handler cat```

- start stockprice-service  
```go run cmd/main.go```

- deploy temp.bpmn

- create new instance with payload `{"symbol": "TEST"}`