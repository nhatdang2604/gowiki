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

	STATIC_FOLDER_PATH = "static/"

	TEXT_FILE_EXTENSION = ".txt"
	EDIT_TEMPLATE_FILENAME = "edit.html"
	VIEW_TEMPLATE_FILENAME = "view.html"

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
	p, err := loadPage(title)
	
	if nil != err {
		url := EDIT_PREFIX + title
		http.Redirect(writer, request, url, http.StatusFound)
	}

	filePath := STATIC_FOLDER_PATH + VIEW_TEMPLATE_FILENAME
	renderTemplate(writer, filePath, p)
}

//Handle the request to URL which has prefix '/edit/'
func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len(EDIT_PREFIX):]
	p, err := loadPage(title)
	if nil != err {
		p = &page{Title: title}
	}
	
	fmt.Println(err)
	filePath := STATIC_FOLDER_PATH + EDIT_TEMPLATE_FILENAME
	renderTemplate(writer, filePath, p)
}

//Handle the execution after saving an editted page
func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len(SAVE_PREFIX):]
	
	//Name of the form's param in the /static/edit.html
	const bodyParam = "body"
	body := request.FormValue(bodyParam)
	
	//Get and save the current page
	p:= &page{Title: title, Body: []byte(body)}
	err := p.save()

	//Handler error after saving the editted page
	if nil != err {
		throwInternalError(writer, err)
		return
	}

	//Redirect to the view page
	url := VIEW_PREFIX + title
	http.Redirect(writer, request, url, http.StatusFound)
}

//Render a page to the html template into the browser, with a given
// response, a html template file path, and the page to render
func renderTemplate(writer http.ResponseWriter, filePath string, p *page) {
	templ, err := template.ParseFiles(filePath)
	
	//Handle error come frome template.ParseFiles()
	if nil != err {
		throwInternalError(writer, err)
		return
	}

	err = templ.Execute(writer, p)

	//Handle error come from templ.Execute()
	if nil != err {
		throwInternalError(writer, err)		
	}
}

func throwInternalError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}

func main() {
	http.HandleFunc(VIEW_PREFIX, viewHandler)
	http.HandleFunc(EDIT_PREFIX, editHandler)
	http.HandleFunc(SAVE_PREFIX, saveHandler)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}
