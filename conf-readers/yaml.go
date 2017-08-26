package conf-readers

import (
	"fmt"
	"io/ioutil"
    "strings"
	"log"
)

func ReadYAMLConfig(filename string) map[string]string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
    content_array := strings.Split(string(content), "\n")
    content_map := make(map[string]string)
    for i := 0; i < len(content_array); i++ {
        row := strings.Split(content_array[i], ": ")
        if len(row) == 2 {
            content_map[row[0]] = row[1]
        }
    }
	return content_map
}

func main() {
    fmt.Println(readConfig("sample.yaml"))
}
