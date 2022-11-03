package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type OriginalJsonType struct {
	SeriesNumber string `json:"Series Number"`
	FileName     string `json:"Filename"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	Gender       string `json:"Gender"`
	Attributes   string `json:"Attributes"`
	UUID         string `json:"UUID"`
}

type NewJsonType struct {
	SeriesNumber string `json:"Series Number"`
	FileName     string `json:"Filename"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	Gender       string `json:"Gender"`
	Attributes   string `json:"Attributes"`
	UUID         string `json:"UUID"`
	Sha256       string `json:"sha256"`
}

func makeMapFromString(updatedRow string) map[string]string {
	newMapfromstring := map[string]string{}
	if err := json.Unmarshal([]byte(updatedRow), &newMapfromstring); err != nil {
		panic(err)
	}
	return newMapfromstring
}

func computeArrayOfRowData(rowDataMap map[string]string) []string {
	return []string{
		rowDataMap["Series Number"],
		rowDataMap["Filename"],
		rowDataMap["Name"],
		rowDataMap["Description"],
		rowDataMap["Gender"],
		rowDataMap["Attributes"],
		rowDataMap["UUID"],
		rowDataMap["sha256"],
	}
}

func doesRowStartWithInteger(row []string) bool {
	var result bool = false
	if _, err := strconv.Atoi(row[0]); err == nil {
		result = true

	}

	return result
}

func generateJsonFromCSVRow(rowData []string) []byte {

	jsonObject := constructJsonStructFromRowData(rowData)

	jsonByte, err := json.Marshal(jsonObject)
	if err != nil {
		panic(err)
	}

	return jsonByte
}

func constructJsonStructFromRowData(rowData []string) OriginalJsonType {
	jsonObject := OriginalJsonType{
		SeriesNumber: rowData[0],
		FileName:     rowData[1],
		Name:         rowData[2],
		Description:  rowData[3],
		Gender:       rowData[4],
		Attributes:   rowData[5],
		UUID:         rowData[6],
	}

	return jsonObject

}

func makeNewRowWithJsonHash(rowData []string, hash string) []byte {
	jsonObject := constructJsonStructFromRowData(rowData)
	newJsonObject := NewJsonType{
		FileName:     jsonObject.FileName,
		SeriesNumber: jsonObject.SeriesNumber,
		Name:         jsonObject.Name,
		Description:  jsonObject.Description,
		Gender:       jsonObject.Gender,
		Attributes:   jsonObject.Attributes,
		UUID:         jsonObject.UUID,
		Sha256:       hash,
	}
	newJson, err := json.Marshal(newJsonObject)

	if err != nil {
		panic(err)
	}

	return newJson
}

func generateHashFromJson(jsonData []byte) string {
	hash := sha256.New()
	hash.Write(jsonData)
	return string(hash.Sum(nil))
}

func main() {
	file, err := os.Open("./temp-hng-csv-file-1.csv")

	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(file)
	headers, _ := reader.Read()
	records, _ := reader.ReadAll()

	writeFile, err := os.Create("./output.csv")
	if err != nil {
		fmt.Println(err)
	}

	writer := csv.NewWriter(writeFile)
	err = writer.Write(append(headers, "Sha256"))
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range records {
		var err error
		if doesRowStartWithInteger(row) == true {
			jsonFromCSV := generateJsonFromCSVRow(row)
			hash := generateHashFromJson(jsonFromCSV)
			updatedRow := string(makeNewRowWithJsonHash(row, hash))
			newRowData := computeArrayOfRowData(makeMapFromString(updatedRow))
			err = writer.Write(newRowData)

		} else {
			err = writer.Write(row)
		}

		if err != nil {
			fmt.Println(err)
		}

	}

}
