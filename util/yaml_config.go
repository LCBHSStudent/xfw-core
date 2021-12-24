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

	createIntMapObjectFromList("group-enable-send-randomgck")
	createIntMapObjectFromList("user-black-list")
}

func GetObjectByKey(key string) interface{} {
	return config[key]
}

func createIntMapObjectFromList(key string) {
	if _, ok := config[key]; ok {
		mapObject := make(map[int64]bool)
		originSlice := config[key].([]interface{})

		for index := range originSlice {
			mapObject[int64(originSlice[index].(int))] = false
		}
		config[key] = mapObject
	}
}
