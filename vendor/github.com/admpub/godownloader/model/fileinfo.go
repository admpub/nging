package model

type FileInfo struct {
	Size     int64    `json:"Size"`
	FileName string   `json:"FileName"`
	Url      string   `json:"Url"`
	Pipes    []string `json:"Pipes"`
}
