package notice

import "sync"

var (
	poolMessage = sync.Pool{
		New: func() interface{} {
			return &Message{}
		},
	}
	poolProgressInfo = sync.Pool{
		New: func() interface{} {
			return &ProgressInfo{}
		},
	}
)

func acquireMessage() *Message {
	return poolMessage.Get().(*Message)
}

func releaseMessage(m *Message) {
	poolMessage.Put(m)
}

func acquireProgressInfo() *ProgressInfo {
	return poolProgressInfo.Get().(*ProgressInfo)
}

func releaseProgressInfo(m *ProgressInfo) {
	m.reset()
	poolProgressInfo.Put(m)
}
