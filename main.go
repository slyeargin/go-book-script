package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

const TIMELAYOUT = "2006-01-02"

type Book struct {
	Title        string   `json:"title"`
	Authors      []string `json:"authors"`
	DateFinished string   `json:"date_finished"`
	Isbn         string   `json:"isbn,omitempty"`
}

func main() {
	logfile, _ := os.Create("app.log")

	// open file
	f, err := os.Open("imports/goodreads-export.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	importedList := createBookList(data)

	// fill in missing ISBNs
	isbnList := getMissingISBNs(importedList)

	// convert to json
	jsonData, err := json.MarshalIndent(isbnList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// save json file
	err = os.WriteFile("goodreads.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()
	log.SetOutput(logfile)
}

type Config struct {
	GoogleBooksAPIKey string
}

func createBookList(data [][]string) []Book {
	// convert csv lines to array of structs
	var bookList []Book
	for i, line := range data {
		if i > 0 { // omit header line
			var book Book
			if line[14] == "" {
				// only save read books
				continue
			} else {
				book.DateFinished = strings.ReplaceAll(line[14], "/", "-")
			}
			book.Title = line[1]
			var authors []string
			authors = append(authors, line[2])
			if line[4] != "" {
				authors = append(authors, line[4])
			}
			book.Authors = authors
			book.Isbn = cleanIsbn(line[6])
			bookList = append(bookList, book)
		}
	}
	return bookList
}

func cleanIsbn(isbn string) string {
	return strings.TrimFunc(isbn, func(r rune) bool {
		return !unicode.IsNumber(r)
	})
}

type GoogleBooksResponse struct {
	TotalItems int          `json:"totalItems"`
	Items      []GoogleBook `json:"items"`
}

type GoogleBook struct {
	VolumeInfo GoogleBookVolumeInfo `json:"volumeInfo"`
}

type GoogleBookVolumeInfo struct {
	IndustryIdentifiers []IndustryIdentifier `json:"industryIdentifiers"`
	ImageLinks          map[string]string
}

type IndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func getMissingISBNs(bookList []Book) []Book {
	for i, book := range bookList {
		if book.Isbn != "" {
			continue
		}

		requestUrl := strings.ReplaceAll("https://www.googleapis.com/books/v1/volumes?q=intitle:"+book.Title+"+inauthor:"+book.Authors[0], " ", "%20")

		resp, err := http.Get(requestUrl)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		googleBooksResponse := GoogleBooksResponse{}
		err = json.Unmarshal(body, &googleBooksResponse)

		if err != nil {
			continue
		}

		if googleBooksResponse.TotalItems == 0 {
			continue
		}

		bestMatch := googleBooksResponse.Items[0].VolumeInfo
		for _, identifier := range bestMatch.IndustryIdentifiers {
			if identifier.Type == "ISBN_13" {
				bookList[i].Isbn = identifier.Identifier
			}
		}
	}

	return bookList
}
