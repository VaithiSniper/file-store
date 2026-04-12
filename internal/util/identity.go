package util

type ServerIdentity struct {
	Id   string
	Name string
}

var GlobalServerIdentity ServerIdentity

func SetIdentity(serverIdentity ServerIdentity) {
	GlobalServerIdentity = serverIdentity
}

func GetIdentity() ServerIdentity {
	return GlobalServerIdentity
}
