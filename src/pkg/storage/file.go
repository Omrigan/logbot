package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
)

type file struct {
	f    *os.File
	csvW *csv.Writer
}

func (m *file) Flush() error {
	m.csvW.Flush()
	return m.f.Sync()
}
func (m *file) Close() error {
	m.csvW.Flush()
	return m.f.Close()
}

type FileConfig struct {
	Dir string `yaml:"dir"`
}

type FileStorage struct {
	cfg *FileConfig
}

func NewFileStorage(cfg *FileConfig) *FileStorage {
	return &FileStorage{
		cfg: cfg,
	}
}

func (m *FileStorage) Write(record *Record) error {
	err := os.Mkdir(m.cfg.Dir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	f := &file{}
	fileName := fmt.Sprintf("%s.csv", record.TS.Format("2006-01-02"))
	name := path.Join(m.cfg.Dir, fileName)
	f.f, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.csvW = csv.NewWriter(f.f)

	err = f.csvW.Write(record.Values())
	if err != nil {
		return err
	}

	return f.Close()
}
