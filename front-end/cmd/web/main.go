package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	// Handler for HomePage
	http.HandleFunc("/", http.HandlerFunc(HomePageHandler))

	fmt.Println("Starting Front End Service on Port 80")
	addr := fmt.Sprintf(":%s", webPort)
	// Starting the server on given port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic(err)
	}
}

// Homepage handler function
func HomePageHandler(res http.ResponseWriter, req *http.Request) {
	render(res, "test.page.gohtml")
}

func render(res http.ResponseWriter, templateName string) {
	partialTemplates := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	// 'templateName' should be the first template to be passed to 'ParseFiles()'
	templates := []string{fmt.Sprintf("./cmd/web/templates/%s", templateName)}
	templates = append(templates, partialTemplates...)

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		// Write error to response
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(tmpl.Name())
	if err = tmpl.Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
