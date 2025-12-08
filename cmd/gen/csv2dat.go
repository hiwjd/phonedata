package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type CSVRecord struct {
	Prefix     string
	Phone      string
	Province   string
	City       string
	ISP        string
	TelCode    string
	PostalCode string
	AreaCode   string
}

const (
	VERSION            = "2511"
	INT_LEN            = 4
	CHAR_LEN           = 1
	HEAD_LENGTH        = 8
	PHONE_INDEX_LENGTH = 9
)

var ispMap = map[string]byte{
	"中国移动":      0x01,
	"中国联通":      0x02,
	"中国电信":      0x03,
	"中国广电":      0x07,
	"中国移动虚拟运营商": 0x06,
	"中国联通虚拟运营商": 0x05,
	"中国电信虚拟运营商": 0x04,
	"中国广电虚拟运营商": 0x08,
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: csv2dat <input_csv> <output_dat>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := readCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d records from CSV\n", len(records))

	err = writeDatFile(records, outputFile)
	if err != nil {
		fmt.Printf("Error writing DAT file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote %s\n", outputFile)
}

func readCSV(filename string) ([]CSVRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var csvRecords []CSVRecord
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 8 {
			continue
		}

		csvRecords = append(csvRecords, CSVRecord{
			Prefix:     record[0],
			Phone:      record[1],
			Province:   record[2],
			City:       record[3],
			ISP:        record[4],
			TelCode:    record[5],
			PostalCode: record[6],
			AreaCode:   record[7],
		})
	}

	return csvRecords, nil
}

func writeDatFile(records []CSVRecord, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Slice(records, func(i, j int) bool {
		phoneI, _ := strconv.Atoi(records[i].Phone)
		phoneJ, _ := strconv.Atoi(records[j].Phone)
		return phoneI < phoneJ
	})

	// First pass: build data section and collect unique data entries
	dataSection := make([]byte, 0)
	dataIndexMap := make(map[string]int32)

	// Also collect valid records for indexing
	var validRecords []CSVRecord

	for _, record := range records {
		// Check if phone is valid
		_, err := strconv.Atoi(record.Phone)
		if err != nil {
			continue
		}

		// Data format: province|city|zipcode|areaCode
		dataKey := fmt.Sprintf("%s|%s|%s|%s", record.Province, record.City, record.PostalCode, record.AreaCode)

		if _, exists := dataIndexMap[dataKey]; !exists {
			offset := int32(len(dataSection))
			dataIndexMap[dataKey] = offset

			dataStr := dataKey + "\000"
			dataSection = append(dataSection, []byte(dataStr)...)
		}

		validRecords = append(validRecords, record)
	}

	// Calculate offsets
	// Data section starts right after header
	dataSectionOffset := int32(HEAD_LENGTH)
	// Index section starts after data section
	indexOffset := dataSectionOffset + int32(len(dataSection))

	// Write index entries to buffer
	var indexData []byte
	for _, record := range validRecords {
		phoneNum, err := strconv.Atoi(record.Phone)
		if err != nil {
			continue
		}

		dataKey := fmt.Sprintf("%s|%s|%s|%s", record.Province, record.City, record.PostalCode, record.AreaCode)
		relativeOffset := dataIndexMap[dataKey]
		absoluteOffset := dataSectionOffset + relativeOffset

		ispType := ispMap[record.ISP]
		if ispType == 0 {
			ispType = 0x03
		}

		indexEntry := make([]byte, PHONE_INDEX_LENGTH)
		binary.LittleEndian.PutUint32(indexEntry[0:4], uint32(phoneNum))
		binary.LittleEndian.PutUint32(indexEntry[4:8], uint32(absoluteOffset))
		indexEntry[8] = ispType

		indexData = append(indexData, indexEntry...)
	}

	// Write header
	header := make([]byte, HEAD_LENGTH)
	copy(header[0:4], VERSION)
	// first record offset should point to where the first index entry is
	binary.LittleEndian.PutUint32(header[4:8], uint32(indexOffset))
	_, err = file.Write(header)
	if err != nil {
		return err
	}

	// Write data section first
	_, err = file.Write(dataSection)
	if err != nil {
		return err
	}

	// Write index section
	_, err = file.Write(indexData)
	if err != nil {
		return err
	}

	return nil
}
