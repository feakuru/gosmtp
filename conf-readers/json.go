package conf-readers

import (
	"fmt"
	"io/ioutil"
    "strings"
	"log"
)

func ReadJSONConfig(filename string) map[string]string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
    content_array := strings.Split(string(content), ",")
    content_map := make(map[string]string)
    for i := 0; i < len(content_array); i++ {
        row := strings.Split(content_array[i], ": ")
        if len(row) == 2 {
			key := strings.Trim(row[0], "\" \n\r,{}")
			val := strings.Trim(row[1], "\" \n\r,{}")
            content_map[key] = val
        }
    }
	return content_map
}

func main() {
    fmt.Println(readConfig("sample.json"))
}
