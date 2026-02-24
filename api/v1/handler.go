package v1

import (
	"elasticsearch/config"
	"elasticsearch/storage"
)

type handlerV1 struct {
	cfg  *config.Config
	strg storage.StorageI
}

type HandleV1 struct {
	Cfg  *config.Config
	Strg storage.StorageI
}

func New(h *HandleV1) *handlerV1 {
	return &handlerV1{
		cfg:  h.Cfg,
		strg: h.Strg,
	}
}
