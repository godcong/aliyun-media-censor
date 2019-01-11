package util

import (
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

/*RandomKind RandomKind */
type RandomKind int

/*random kinds */
const (
	RandomNum      RandomKind = iota // 纯数字
	RandomLower                      // 小写字母
	RandomUpper                      // 大写字母
	RandomLowerNum                   // 数字、小写字母
	RandomUpperNum                   // 数字、大写字母
	RandomAll                        // 数字、大小写字母
)

/*RandomString defines */
var (
	RandomString = map[RandomKind]string{
		RandomNum:      "0123456789",
		RandomLower:    "abcdefghijklmnopqrstuvwxyz",
		RandomUpper:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomLowerNum: "0123456789abcdefghijklmnopqrstuvwxyz",
		RandomUpperNum: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomAll:      "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

//GenerateRandomString2 随机字符串
func GenerateRandomString2(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//GenerateRandomString 随机字符串
func GenerateRandomString(size int, kind ...RandomKind) string {
	bytes := RandomString[RandomAll]
	if kind != nil {
		if k, b := RandomString[kind[0]]; b == true {
			bytes = k
		}
	}
	var result []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func UnmarshalJSON(reader io.Reader, v interface{}) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	log.Println(string(bytes))
	err = jsoniter.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

func MarshalJSON(v interface{}) ([]byte, error) {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func FileList(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func Files(path string, pattern ...string) []string {
	if pattern == nil {
		path += "/*.ts"
	} else {
		path += "/" + pattern[0]
	}
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil
	}
	return matches
}

func IsDir(name string) bool {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return false
	}
	return fi.IsDir()
}
