package imbot

type Messager interface {
	BuildURL(...string) string
	Reset() Messager
	SendText(url, text string) error
	SendMarkdown(url, title, markdown string) error
}
