package ttlmap

type pqitem struct {
	key   string
	item  *Item
	index int
}

type pqueue []*pqitem

func (pq pqueue) Len() int {
	return len(pq)
}

func (pq pqueue) Less(i, j int) bool {
	pqi := pq[i].item
	pqj := pq[j].item
	if pqi.expires {
		if pqj.expires {
			return pqi.expiration.Before(pqj.expiration)
		}
		return true
	}
	return false
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *pqueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqitem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq pqueue) peek() *pqitem {
	if len(pq) == 0 {
		return nil
	}
	return pq[0]
}
