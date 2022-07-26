package main

import (
	"regexp"
	"os"
	"html/template"
	"log"
	"net/http"
	"errors"
)

const (
	PORT = "8088"

	STATIC_FOLDER_PATH = "static/"
	DATA_FOLDER_PATH = "data/"

	TEXT_FILE_EXTENSION = ".txt"
	EDIT_TEMPLATE_FILENAME = "edit.html"
	VIEW_TEMPLATE_FILENAME = "view.html"

	VIEW_PREFIX = "/view/"
	EDIT_PREFIX = "/edit/"
	SAVE_PREFIX = "/save/"
)

var (
	templates = template.Must(template.ParseFiles(
		STATIC_FOLDER_PATH + EDIT_TEMPLATE_FILENAME, 
		STATIC_FOLDER_PATH + VIEW_TEMPLATE_FILENAME))

	validPath = regexp.MustCompile("^(" + 
					VIEW_PREFIX + "|" + 
					EDIT_PREFIX + "|" + 
					SAVE_PREFIX + 
					")([a-zA-Z0-9]+)$")
)

type page struct {
	Title string
	Body []byte
}



//Save a page body to text file
func (p *page) save() error {
	filename := p.Title + TEXT_FILE_EXTENSION
	filepath := DATA_FOLDER_PATH + filename
	return os.WriteFile(filepath, p.Body, 0600)
}

//Parse a page after reading it
func loadPage(title string) (*page, error) { 
	filename := title + TEXT_FILE_EXTENSION
	filepath := DATA_FOLDER_PATH + filename
	body, err := os.ReadFile(filepath)

	//return the error if encoutering it
	if nil != err {
		return nil, err
	}

	//return the page with the nil error
	return &page{Title: title, Body: body}, err

}


//Handle the request to URL which has prefix '/view/'
func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	p, err := loadPage(title)
	
	if nil != err {
		url := EDIT_PREFIX + title
		http.Redirect(writer, request, url, http.StatusFound)
	}

	fileName := VIEW_TEMPLATE_FILENAME
	renderTemplate(writer, fileName, p)
}

//Handle the request to URL which has prefix '/edit/'
func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	p, err := loadPage(title)
	if nil != err {
		p = &page{Title: title}
	}
	
	fileName := EDIT_TEMPLATE_FILENAME
	renderTemplate(writer, fileName, p)
}

//Handle the execution after saving an editted page
func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {	
	
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

	//Redirect to the view page`
	url := VIEW_PREFIX + title
	http.Redirect(writer, request, url, http.StatusFound)
}

//Render a page to the html template into the browser, with a given
// response, a html template file path, and the page to render
func renderTemplate(writer http.ResponseWriter, filePath string, p *page) {
	
	err := templates.ExecuteTemplate(writer, filePath, p)
	if nil != err {
		throwInternalError(writer, err)		
	}
}

//Throw the StatusInternalServerError
func throwInternalError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}

//Validate the title come from URL
func getTitle(writer http.ResponseWriter, request *http.Request) (string, error) {
	matchers := validPath.FindStringSubmatch(request.URL.Path)
	if nil == matchers {
		http.NotFound(writer, request)
		return "", errors.New("Invalid Page Title")
	}

	return matchers[2], nil
}

//Apply closure to wrapp the handlers function
func makeHandler(handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, err := getTitle(writer, request)
		
		//No need to handle error, because the getTitle() method already do that
		if nil != err {
			return
		}

		handler(writer, request, title)
	}
}

func main() {
	http.HandleFunc(VIEW_PREFIX, makeHandler(viewHandler))
	http.HandleFunc(EDIT_PREFIX, makeHandler(editHandler))
	http.HandleFunc(SAVE_PREFIX, makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}
