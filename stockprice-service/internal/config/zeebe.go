package config

type Zeebe struct {
	Server
}

const zeebePath = "zeebe"

func GetZeebe() Zeebe {
	return Zeebe {
		Server: GetServer(zeebePath),
	}
}