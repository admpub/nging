package hook

import "sort"

type (
	Hook  func() error
	Hooks map[string][]Hook
)

func New() Hooks {
	return Hooks{}
}

func (h Hooks) On(ev string, fn Hook) {
	if _, ok := h[ev]; !ok {
		h[ev] = []Hook{}
	}
	h[ev] = append(h[ev], fn)
}

func (h Hooks) Size(ev string) int {
	if _, ok := h[ev]; !ok {
		return 0
	}
	return len(h[ev])
}

func (h Hooks) Names() []string {
	names := make([]string, len(h))
	var i int
	for name := range h {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

func (h Hooks) Off(ev string) {
	if _, ok := h[ev]; !ok {
		return
	}
	delete(h, ev)
}

func (h Hooks) Fire(ev string) error {
	if _, ok := h[ev]; !ok {
		return nil
	}
	var err error
	for _, hook := range h[ev] {
		if err = hook(); err != nil {
			return err
		}
	}
	return err
}
