FROM golang:alpine AS build

ENV GO111MODULE=on

RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser

RUN go get github.com/golang/protobuf/protoc-gen-go@v1.3.2

WORKDIR $GOPATH/src/github.com/domahidizoltan/playground-workflow-engine/stockprice-service
COPY . .

RUN mkdir /app
RUN cp ./config.yml /app/config.yml
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/service cmd/main.go

###
FROM scratch

WORKDIR /app
COPY --from=build /app .

COPY --from=build /etc/passwd /etc/passwd
USER appuser

EXPOSE 9100
ENTRYPOINT ["/app/service"]