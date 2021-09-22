package poet

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	util "github.com/LCBHSStudent/xfw-core/util"
)

var poemList []string
var poemCount int64

func init() {
	poemPath := util.GetObjectByKey("poem-file-path").(string)

	err := filepath.Walk(poemPath,
    func(path string, info os.FileInfo, err error) error {
    	if err != nil {
        	return err
    	}
    	if strings.HasSuffix(path, ".pt") {
			poemList = append(poemList, path)
		}
    	return nil
	})

	if err != nil {
		log.Fatal(err)
	} else {
		poemCount = int64(len(poemList))
		log.Printf("Amount of loaded poem: %v", len(poemList))
	}
}

// const [13]int
var lineCountProbability = []int {32, 15, 10, 10, 8, 6, 4, 3, 3, 2, 2, 2, 2}
var filter = regexp.MustCompile("(title.*)|(date.*)|([\\w\n].*)")
	
func GetPoetry() string {
	
	bProb, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}

	prob := int(bProb.Uint64())
	var probSum = 0
	var lineCount int

	for i, v := range lineCountProbability {
		probSum += v
		if prob <= probSum {
			lineCount = i + 2
			break
		}
	}

	var ret string

	for i := 0; i < lineCount; i++ {
		bFile, err := rand.Int(rand.Reader, big.NewInt(poemCount))
		if err != nil {
			log.Fatal(err)
		}

		fileIdx := int(bFile.Uint64())

		lines := util.ReadLine(poemList[fileIdx], filter)
		
		bFile, err = rand.Int(rand.Reader, big.NewInt(int64(len(lines))))
		if err != nil {
			log.Fatal(err)
		} else {
			ret += lines[bFile.Uint64()]
			if i != lineCount - 1 {
				ret += "\n"
			}
		}
	}


	return ret
}
