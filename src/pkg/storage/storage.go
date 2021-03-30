package storage

import (
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type Record struct {
	Item    string    `yaml:"item"`
	Param   string    `yaml:"param"`
	Comment string    `yaml:"comment"`
	TS      time.Time `yaml:"ts"`
}

func (m *Record) Time() string {
	return m.TS.Format(time.RFC822)
}

func (m *Record) Values() []string {
	return []string{m.Item, m.Param, m.Comment, m.Time()}
}
func (m *Record) String() string {
	return strings.Join(m.Values(), ",")
}
func (m *Record) YAML() string {
	out, err := yaml.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

type Storage interface {
	Write(record *Record) error
}

