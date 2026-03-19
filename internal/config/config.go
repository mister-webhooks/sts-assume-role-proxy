package config

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

type ServerConfiguration struct {
	SocketPath  string      `yaml:"socket_path"`
	AccessTable AccessTable `yaml:"rules"`
}

type AccessTable struct {
	table map[string]map[uint]string
}

func (at *AccessTable) UnmarshalYAML(node *yaml.Node) error {
	*at = *NewAccessTable()

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected root node of access table to be mapping")
	}

	if len(node.Content)%2 != 0 {
		return fmt.Errorf("mapping node contents must be key/value pairs")
	}

	for i := 0; i < len(node.Content); i = i + 2 {
		namespaceKeyNode := node.Content[i]
		namespaceValues := node.Content[i+1]

		if namespaceKeyNode.Tag != "!!str" {
			return fmt.Errorf("expected namespace to be a string, got %s", namespaceKeyNode.Tag)
		}

		if len(namespaceValues.Content)%2 != 0 {
			return fmt.Errorf("inner mapping node contents must be key/value pairs")
		}

		for j := 0; j < len(namespaceValues.Content); j = j + 2 {
			useridNode := namespaceValues.Content[j]
			arnNode := namespaceValues.Content[j+1]

			if useridNode.Tag != "!!int" {
				return fmt.Errorf("expected userid to be an int, got %s", useridNode.Tag)
			}

			userid, err := strconv.ParseUint(useridNode.Value, 10, 32)

			if err != nil {
				return err
			}

			if arnNode.Tag != "!!str" {
				return fmt.Errorf("expected value to be a string, got %s", arnNode.Tag)
			}

			at.Insert(namespaceKeyNode.Value, uint(userid), arnNode.Value)
		}
	}

	return nil
}

func NewAccessTable() *AccessTable {
	return &AccessTable{
		table: make(map[string]map[uint]string),
	}
}

func NewConfigurationFromYAML(data []byte) (*ServerConfiguration, error) {
	scfg := ServerConfiguration{
		AccessTable: *NewAccessTable(),
	}

	err := yaml.Load(data, &scfg)

	if err != nil {
		return nil, err
	}

	return &scfg, nil
}

func (at *AccessTable) Insert(namespace string, userid uint, arn string) {
	if _, ok := at.table[namespace]; !ok {
		at.table[namespace] = make(map[uint]string)
	}

	at.table[namespace][userid] = arn
}

func (at *AccessTable) Lookup(namespace string, userid uint) (string, bool) {
	nstable, present := at.table[namespace]

	if !present {
		return "", false
	}

	arn, present := nstable[userid]
	return arn, present
}
