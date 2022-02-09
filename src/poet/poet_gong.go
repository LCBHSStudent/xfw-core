package poet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"
)

// prefix[0]; suffix[1]
var gongFormat = [8][2]int{
	{1, 2}, {2, 1}, {2, 2}, {2, 3}, {3, 2}, {1, 4}, {4, 1}, {3, 4},
}

// count of sentence [0]; probility [1]
var gongConfig = [4][2]int{
	{2, 50}, {4, 42}, {6, 5}, {8, 3},
}

func getPinyinResult(pinyin *string) (res string, err error) {
	res = ""
	url := fmt.Sprintf("https://inputtools.google.com/request?text=%s&itc=zh-t-i0-pinyin&num=%d&cp=0&cs=1&ie=utf-8&oe=utf-8&app=demopage", *pinyin, 10)

	defer func(res *string, e *error) {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}(&res, &err)

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New("Failed to get result of PinYin: " + *pinyin)
		return
	}
	jsonBytes := []byte{}
	jsonBytes = append(jsonBytes, []byte("{\"value\":")...)
	jsonBytes = append(jsonBytes, body...)
	jsonBytes = append(jsonBytes, []byte("}")...)

	jsonValue := make(map[string]interface{})
	err = json.Unmarshal(jsonBytes, &jsonValue)
	if err != nil {
		return
	}
	candidatesItf := jsonValue["value"].([]interface{})[1].([]interface{})[0].([]interface{})[1]
	candidates := make([]string, 0, len(candidatesItf.([]interface{})))
	for _, candidate := range candidatesItf.([]interface{}) {
		candidates = append(candidates, candidate.(string))
	}
	res = candidates[rand.Intn(len(candidates))%6]

	return
}

func GenerateGongPoem() (gongPoem string) {
	rand.Seed(time.Now().UnixNano())
	var hzRegexp = regexp.MustCompile(`^[\p{Han}]*$`)

	gongPoem = "龚诗来力~\n\n"
	sentencesCount := 0
	base := 0
	prob := rand.Intn(100)
	for _, content := range gongConfig {
		base += content[1]
		if prob < base {
			sentencesCount = content[0]
			break
		}
	}

	allowDifferentFormat := rand.Intn(5) >= 2
	format := gongFormat[rand.Intn(len(gongFormat))]
	for i := 0; i < sentencesCount; i++ {
		if allowDifferentFormat {
			format = gongFormat[rand.Intn(len(gongFormat))]
		}
		for {
			prefix := ""
			suffix := ""
			for index, charCnt := range format {
				for j := 0; j < charCnt; j++ {
					if index == 0 {
						prefix += string(rune('a' + rand.Intn(26)))
					} else {
						suffix += string(rune('a' + rand.Intn(26)))
					}
				}
			}
			prefixRes, err := getPinyinResult(&prefix)
			if err != nil {
				log.Println(err)
				return
			}
			suffixRes, err := getPinyinResult(&suffix)
			if err != nil {
				log.Println(err)
				return
			}
			if utf8.RuneCountInString(prefixRes) == format[0] &&
				utf8.RuneCountInString(suffixRes) == format[1] &&
				hzRegexp.MatchString(prefixRes) &&
				hzRegexp.MatchString(suffixRes) {
				gongPoem += prefixRes + suffixRes + "\n"
				break
			}
		}
	}

	return
}
