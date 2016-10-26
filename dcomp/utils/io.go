package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadYaml(fname string, config interface{}) error {
	yamlFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}
	return nil
}
