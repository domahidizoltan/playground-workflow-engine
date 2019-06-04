FROM golang:alpine AS build

ENV GO111MODULE=on

RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/domahidizoltan/playground-workflow-engine
COPY . .

RUN mkdir -p /app/configs
RUN cp config.yml /app/config.yml
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/mockserver cmd/main.go

###
FROM scratch

WORKDIR /app
COPY --from=build /app .

COPY --from=build /etc/passwd /etc/passwd
USER appuser

EXPOSE 9000
ENTRYPOINT ["/app/mockserver"]