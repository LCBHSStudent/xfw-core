package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
)

func ReadLine(filepath string, re *regexp.Regexp) []string {
	lines := make([]string, 0)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := readLine(reader)
		if err != nil {
			break
		} else {
			if len(line) == 0 {
				continue
			}
			if re != nil && re.MatchString(line) {
				continue
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func readLine(reader *bufio.Reader) (string, error) {
	line, isprefix, err := reader.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = reader.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

func GetStructureDataInfo(itf interface{}, structName string) string {

	structType := reflect.TypeOf(itf)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		panic("can not use not-structure arguments" +
			" in function \"GetStructureDataInfo\"")
	}

	structValue := reflect.ValueOf(itf)

	info := fmt.Sprintf("\t%s {\n" /*Green(*/, structName /*)*/)

	for i := 0; i < structType.NumField(); i++ {
		// 若想层层展开结构体，针对reflect.Struct进行递归调用
		typeInfo := structType.Field(i).Name
		valueInfo := structValue.Field(i)

		info += fmt.Sprintf("\t\t%20s:\t %v\n", typeInfo, valueInfo)
	}
	info += "\t}"
	return info
}
