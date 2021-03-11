package imbot

type Messager interface {
	BuildURL(...string) string
	Reset() Messager
	SendText(url, text string, atMobiles ...string) error
	SendMarkdown(url, title, markdown string, atMobiles ...string) error
}

type Message struct {
	Name     string
	Label    string
	Messager Messager
}

var messagers = map[string]*Message{}

func Register(name string, label string, m Messager) {
	messagers[name] = &Message{Name: name, Label: label, Messager: m}
}

func Messagers() map[string]*Message {
	return messagers
}

func Open(name string) *Message {
	m, _ := messagers[name]
	return m
}

func Unregister(names ...string) {
	for _, name := range names {
		if _, ok := messagers[name]; ok {
			delete(messagers, name)
		}
	}
}
