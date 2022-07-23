package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
)

const (
	PORT = "8088"

	TEXT_FILE_EXTENSION = ".txt"
	
	VIEW_PREFIX = "/view/"
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
	title := request.URL.Path[len(VIEW_PREFIX):]
	p, _ := loadPage(title)
	fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc(VIEW_PREFIX, viewHandler)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}
