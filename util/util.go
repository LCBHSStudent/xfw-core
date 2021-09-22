package util

import (
	"bufio"
	"log"
	"os"
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
