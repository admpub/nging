package imbot

type Messager interface {
	SendText(url, text string) error
	SendMarkdown(url, title, markdown string) error
}
