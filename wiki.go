package main

import (
	"fmt"
	"os"
)

type page struct {
	Title string
	Body []byte
}

//Save a page body to text file
func (p *page) save() error {
	const filename = p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}


