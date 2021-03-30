package storage

import (
	"fmt"
)

type MetaStorage struct {
	storages []Storage
}

func NewMetaStorage(storages ...Storage) *MetaStorage {
	return &MetaStorage{
		storages: storages,
	}
}
func (m *MetaStorage) Write(record *Record) error {
	var err error
	for _, s := range  m.storages {
		curErr := s.Write(record)
		if curErr != nil {
			if err == nil {
				err = curErr
				continue
			}
			err = fmt.Errorf("%s; %w", err.Error(), curErr)
		}
	}
	return err
}
