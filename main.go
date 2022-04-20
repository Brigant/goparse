package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func findValue(str string, like string) string {
	result := ""
	if strings.Contains(str, like) {
		result = strings.Trim(str, " ")
		result = strings.Replace(result, like, "", -1)
		return result
	} else {
		return "NotFound"
	}
}

func main() {
	csvFile, _ := os.Open("input.csv")
	reader := csv.NewReader(csvFile)
	//	var ListContacts [][]string
	outputFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvwriter := csv.NewWriter(outputFile)
	columnTitle := []string{"EMAIL", "PHONE", "NAME", "LASTNAME", "CHAT", "COMPANY", "COUNTRY", "COMMNENT"}
	_ = csvwriter.Write(columnTitle)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		array := strings.Split(line[0], "\n")
		row := []string{"", "", "", "", "", "", "", ""}
		for _, item := range array {
			name := findValue(item, "* Name:")
			email := findValue(item, "* Email:")
			phone := findValue(item, "* Phone number:")
			sname := findValue(item, "* Last Name:")
			chat := findValue(item, "* Your Skype/WeChat/WhatsApp ID:")
			company := findValue(item, "* Company Name:")
			country := findValue(item, "* Country:")
			comment := findValue(item, "* Comment:")
			if name != "NotFound" {
				row[2] = name
			}
			if email != "NotFound" {
				row[0] = email
			}
			if phone != "NotFound" {
				row[1] = phone
			}
			if sname != "NotFound" {
				row[3] = sname
			}
			if chat != "NotFound" {
				row[4] = chat
			}
			if company != "NotFound" {
				row[5] = company
			}
			if country != "NotFound" {
				row[6] = country
			}
			if comment != "NotFound" {
				row[7] = comment
			}
		}

		if row[0] != "" {
			_ = csvwriter.Write(row)
		}
	}
	csvwriter.Flush()
	csvFile.Close()
	outputFile.Close()

}
