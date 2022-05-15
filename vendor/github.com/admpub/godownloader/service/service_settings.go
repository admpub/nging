package service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/admpub/godownloader/model"
)

type DownloadSettings struct {
	FI model.FileInfo           `json:"FileInfo"`
	Dp []model.DownloadProgress `json:"DownloadProgress"`
}

type ServiceSettings struct {
	Ds []DownloadSettings
}

func LoadFromFile(s string) (*ServiceSettings, error) {
	sb, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}

	var ss ServiceSettings
	err = json.Unmarshal(sb, &ss)
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

func (s *ServiceSettings) SaveToFile(fp string) error {
	dat, err := json.Marshal(s)
	if err != nil {
		return err
	}
	log.Println("info: try save settings")
	log.Println(string(dat))
	err = ioutil.WriteFile(fp, dat, 0664)
	if err != nil {
		return err
	}
	return nil
}
