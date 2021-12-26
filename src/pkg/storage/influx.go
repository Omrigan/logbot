package storage

import (
	"context"
	"strconv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxConfig struct {
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
	Org      string `yaml:"org"`
	Bucket   string `yaml:"bucket"`
}

type InfluxStorage struct {
	cfg      *InfluxConfig
	writeAPI api.WriteAPIBlocking
}

func NewInfluxStorage(cfg *InfluxConfig) *InfluxStorage {
	storage := &InfluxStorage{
		cfg: cfg,
	}
	client := influxdb2.NewClient(cfg.Endpoint, cfg.Token)
	storage.writeAPI = client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	return storage
}

func (m *InfluxStorage) Write(record *Record) error {
	fields := map[string]interface{}{"count": 1, "comment": record.Comment}
	if record.Param != "" {
		f, err := strconv.ParseFloat(record.Param, 64)
		if err == nil {
			fields["param_float"] = f
		} else {
			fields["param"] = record.Param
		}
	}

	p := influxdb2.NewPoint(record.Item,
		map[string]string{},
		fields,
		record.TS)
	// write point immediately
	return m.writeAPI.WritePoint(context.Background(), p)
}
