package humanize

import "sync"

type Humanizers map[string]*Humanizer

func (h Humanizers) Get(language string) *Humanizer {
	lock.Lock()
	defer lock.Unlock()
	humanizer, ok := h[language]
	if ok {
		return humanizer
	}
	return nil
}

func (h Humanizers) Set(language string, humanizer *Humanizer) {
	lock.Lock()
	h[language] = humanizer
	lock.Unlock()
}

func (h Humanizers) Delete(languages ...string) {
	lock.Lock()
	for _, language := range languages {
		_, ok := h[language]
		if ok {
			delete(h, language)
		}
	}
	lock.Unlock()
}

var (
	humanizers = Humanizers{}
	lock       = sync.RWMutex{}
)

func Cached() Humanizers {
	return humanizers
}

func Clear() {
	lock.Lock()
	humanizers = Humanizers{}
	lock.Unlock()
}
