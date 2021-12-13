package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const configFilePath = "../share/config.yaml"

var config map[interface{}]interface{}

func init() {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config = make(map[interface{}]interface{})

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := config["group-enable-send-randomgck"]; ok {
		mapObject := make(map[int64]bool)
		originSlice := config["group-enable-send-randomgck"].([]interface{})

		for index := range originSlice {
			mapObject[int64(originSlice[index].(int))] = false
		}
		config["group-enable-send-randomgck"] = mapObject
	}
}

func GetObjectByKey(key string) interface{} {
	return config[key]
}
