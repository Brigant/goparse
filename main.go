package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/xuri/excelize/v2"
)

const lenMustHave = 2

type NotUsed struct {
	Name string
}

type Client struct { // Our example struct, you can use "-" to ignore a field
	Id            string `csv:"client_id"`
	Name          string `csv:"client_name"`
	Age           string `csv:"client_age"`
	NotUsedString string `csv:"-"`
	NotUsed       `csv:"-"`
}

type Message struct {
	Title   string `csv:"Title"`
	From    string `csv:"From"`
	To      string `csv:"To"`
	Date    string `csv:"Date"`
	NotUsed string `csv:"-"`
	Body    string `csv:"Body"`
}

type Contact struct {
	FirstName string `csv:"FirstName"`
	LastName  string `csv:"LastName"`
	Email     string `csv:"Email"`
	Chat      string `csv:"Chat"`
	Phone     string `csv:"Phone"`
	Company   string `csv:"Company"`
	Country   string `csv:"Country"`
	Comments  string `csv:"Comments"`
}

func main() {
	args, err := readArgs()
	if err != nil {
		log.Println(err.Error())

		return
	}

	listContacts, err := ReadCSVtoListContact(args["--in"])
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err := SaveToExelFile(args["--out"], listContacts); err != nil {
		log.Println(err)
	}
}

// Read the argument from commanline
// ex: --in=/path and --out=/path .
func readArgs() (map[string]string, error) {
	mArgs := map[string]string{"--in": "", "--out": ""}

	for i, val := range os.Args {
		if i > 0 {
			keyVal := strings.Split(val, "=")
			if len(keyVal) == 2 {
				mArgs[keyVal[0]] = keyVal[1]
			}
		}
	}

	if mArgs["--in"] == "" || mArgs["--out"] == "" {
		return nil, errors.New("you should specify two key-val, --in=/path and --out=/path")
	}

	return mArgs, nil
}

// Read the specified csv file and returns the contacts slice.
func ReadCSVtoListContact(fileName string) ([]Contact, error) {
	listContact := []Contact{}

	messageFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error while open file: %w", err)
	}
	defer messageFile.Close()

	messages := []Message{}

	if err := gocsv.UnmarshalFile(messageFile, &messages); err != nil { // Load clients from file
		return nil, fmt.Errorf("unmarshal file: %w", err)
	}

	if _, err := messageFile.Seek(0, 0); err != nil { // Go to the start of the file
		return nil, fmt.Errorf("seek: %w", err)
	}

	for _, message := range messages {
		listContact = append(listContact, parseBody(message.Body))
	}

	return listContact, nil
}

// Receives the contact slice and save it to the specified exel file.
func SaveToExelFile(fileName string, list []Contact) error {
	outputFile := excelize.NewFile()
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Create a new sheet.
	sheetName := "Contacts"
	index, err := outputFile.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("error while creatin new sheet: %w", err)
	}

	outputFile.SetCellValue(sheetName, "A1", "FirstName")
	outputFile.SetCellValue(sheetName, "B1", "LastName")
	outputFile.SetCellValue(sheetName, "C1", "Email")
	outputFile.SetCellValue(sheetName, "D1", "Chat")
	outputFile.SetCellValue(sheetName, "E1", "Phone")
	outputFile.SetCellValue(sheetName, "F1", "Company")
	outputFile.SetCellValue(sheetName, "G1", "Country")
	outputFile.SetCellValue(sheetName, "H1", "Comments")

	for i, contact := range list {
		row := strconv.Itoa(i + 2)
		outputFile.SetCellValue(sheetName, "A"+row, contact.FirstName)
		outputFile.SetCellValue(sheetName, "B"+row, contact.FirstName)
		outputFile.SetCellValue(sheetName, "C"+row, contact.Email)
		outputFile.SetCellValue(sheetName, "D"+row, contact.Chat)
		outputFile.SetCellValue(sheetName, "E"+row, contact.Phone)
		outputFile.SetCellValue(sheetName, "F"+row, contact.Company)
		outputFile.SetCellValue(sheetName, "G"+row, contact.Country)
		outputFile.SetCellValue(sheetName, "H"+row, contact.Comments)
	}

	outputFile.SetActiveSheet(index)

	// Save spreadsheet by the given path.
	if err := outputFile.SaveAs(fileName); err != nil {
		return fmt.Errorf("error while save xslx file: %w", err)
	}

	return nil
}

// The function works only with body column of the csv file.
// It checks and parse all key and values in the body.
func parseBody(body string) Contact {
	rows := strings.Split(body, "\n")
	contact := Contact{}

	if rows[10] == "Details:" {
		contact.FirstName = getValue(rows[11], "Name")
		contact.LastName = getValue(rows[12], "Last Name")
		contact.Email = getValue(rows[13], "Email")
		contact.Chat = getValue(rows[14], "Your Skype/WeChat/WhatsApp ID")
		contact.Phone = getValue(rows[15], "Phone number")
		contact.Company = getValue(rows[16], "Company Name")
		contact.Country = getValue(rows[17], "Country")
		contact.Comments = getValue(rows[18], "Comment")
	}
	return contact
}

func getValue(keyPairStr, keyName string) string {
	row := strings.Split(keyPairStr, ":")
	if len(row) != lenMustHave {
		return ""
	}

	key := strings.Trim(row[0], " ")
	if key != keyName {
		return ""
	}

	value := strings.Trim(row[1], " ")

	return value
}
