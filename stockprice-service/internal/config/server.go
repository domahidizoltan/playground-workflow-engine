package config

type Http struct {
	Server
}

const httpPath = "http"

func GetHttp() Http {
	return Http {
		Server: GetServer(httpPath),
	}
}
