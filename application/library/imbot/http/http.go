package http

import (
	"encoding/json"
	"github.com/admpub/marmot/miner"
)

func Send(url string, m interface{}) ([]byte,error) {
	worker := miner.NewAPI()
	worker.SetURL(url).SetMaxRetries(3)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	worker.SetBinary(body)
	return worker.PostJSON()
}
