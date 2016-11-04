package utils

import (
	"bytes"
	"compress/gzip"
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

func CompressString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}
