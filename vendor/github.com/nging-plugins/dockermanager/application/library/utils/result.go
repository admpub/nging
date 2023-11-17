package utils

type Result struct {
	completed bool
	Status    string `json:"status"`
	ID        string `json:"id"`
}

func (r *Result) SetCompleted(completed bool) {
	r.completed = completed
}

func (r *Result) Completed() bool {
	return r.completed
}
