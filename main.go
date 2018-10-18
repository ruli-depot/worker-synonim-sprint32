package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

//struct json
type dict struct {
	CatID    int64    `json:"cat_id"`
	DataList []string `json:"data_list"`
}

var (
	dictionaryMap = map[string]bool{}
	dictionary    = map[string][]dict{}
	file          *os.File
)

func main() {

	// start
	start := time.Now()
	log.Println("initiate on :", start.Format("2006-01-02 15:04:05"))

	//make sure filepath is satisfied by user input
	if len(os.Args) < 3 {
		panic("Missing Paramater ! Command : go run main.go directory_to_csv_file directory_to_json_file ")
	}

	//get the json path
	jsonPath := os.Args[2]
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		panic(err)
	}

	//get the filepath
	filePath := os.Args[1]
	file, err = os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	log.Println("filePath :", filePath)

	// check file type
	if !strings.HasSuffix(filePath, ".csv") {
		panic("invalid csv file")
	}
	if !strings.HasSuffix(jsonPath, ".json") {
		panic("invalid json file")
	}

	//decode
	err = json.NewDecoder(jsonFile).Decode(&dictionary)
	if err != nil {
		panic(err)
	}

	// assign json to hash
	generateDictionaryMap()

	// loop csv to hash map
	generateResult()

	log.Println("done in", time.Since(start))
}

func generateDictionaryMap() {
	for _, val := range dictionary["data"] {
		for _, data := range val.DataList {
			dictionaryMap[data] = true
		}
	}

}

func generateResult() {

	path := "result/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	output, err := os.Create(path + "synonim_" + time.Now().Format("2006_02_01") + ".csv")
	if err != nil {
		panic(err)
	}

	writerCSV := csv.NewWriter(output)
	err = writerCSV.Write([]string{"KEYWORD", "0=SYNONIM NOT FOUND | 1= FOUND SYNONIM"})
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		fullWord := strings.Join(row, " ")
		isSynonimAvailable := "0"
		for _, word := range row {
			if dictionaryMap[word] {
				isSynonimAvailable = "1"
				break
			}
		}

		// log.Println(fullWord)
		err = writerCSV.Write([]string{fullWord, isSynonimAvailable})
		if err != nil {
			panic(err)
		}

	}
}
