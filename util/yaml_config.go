package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const configFilePath = "../share/config.yaml"

var config map[interface{}] interface{}

func init() {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config = make(map[interface{}] interface{})

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetObjectByKey(key string) interface{} {
	return config[key]
}
