package main

import (
	"fmt"
	"os"
	"html/template"
	"log"
	"net/http"
)

const (
	PORT = "8088"

	STATIC_FOLDER_PATH = "/static"

	TEXT_FILE_EXTENSION = ".txt"
	EDIT_TEMPLATE_FILENAME = "edit.html"


	VIEW_PREFIX = "/view/"
	EDIT_PREFIX = "/edit/"
	SAVE_PREFIX = "/save/"
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

//Handle the qreust to URL which has prefix '/edit/'
func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len(EDIT_PREFIX):]
	p, err := loadPage(title)
	if nil != err {
		p = &page{Title: title}
	}

	templ, _ := template.ParseFiles(STATIC_FILE_PATH + "/" + EDIT_TEMPLATE_FILENAME)
	templ.Execute(w, p)
}

func main() {
	http.HandleFunc(VIEW_PREFIX, viewHandler)
	http.HandleFunc(EDIT_PREFIX, editHandler)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}
