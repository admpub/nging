package server

type Auth interface {
	CheckPasswd(string, string) (bool, error)
}
