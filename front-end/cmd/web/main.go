package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const webPort = "8081"

//go:embed templates
var templateFS embed.FS

func main() {
	// Handler for HomePage
	http.HandleFunc("/", http.HandlerFunc(HomePageHandler))

	fmt.Println("Starting Front End Service on Port", webPort)
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
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	// 'templateName' should be the first template to be passed to 'ParseFiles()'
	templates := []string{fmt.Sprintf("templates/%s", templateName)}
	templates = append(templates, partialTemplates...)

	tmpl, err := template.ParseFS(templateFS, templates...)
	if err != nil {
		// Write error to response
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// define and create data that we will be passing to go template
	var data struct {
		BrokerURL string
	}
	data.BrokerURL = os.Getenv("BROKER_URL")

	// fmt.Println(tmpl.Name())
	if err = tmpl.Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
