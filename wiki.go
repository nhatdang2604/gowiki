package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
)

const (
	TEXT_FILE_EXTENSION = ".txt"
)

type page struct {
	Title string
	Body []byte
}



//Save a page body to text file
func (p *page) save() error {
	filename := p.Title + TEXT_FILE_EXTENSION
	return os.WriteFile(filename, p.Body, 0600)
}

//Parse a page after reading it
func loadPage(title string) (*page, error) { 
	filename := title + TEXT_FILE_EXTENSION
	body, err := os.ReadFile(filename)

	//return the error if encoutering it
	if nil != err {
		return nil, err
	}

	//return the page with the nil error
	return &page{Title: title, Body: body}, err

}


//Handle the request to URL which has prefix '/view/'
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	const PREFIX = "/view/"
	title := request.URL.Path[len(PREFIX):]
	p, _ := loadPage(title)
	fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	p1 := page{Title: "TestPage", Body: []byte("Test Body")}
	p1.save()
	p2, _ := loadPage(p1.Title)
	fmt.Println(string(p2.Body))
}
