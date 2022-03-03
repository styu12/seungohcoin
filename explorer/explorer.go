package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/styu12/seungohcoin/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks []*blockchain.Block 
}

func handleHome(rw http.ResponseWriter, r *http.Request) {
	return
	// data := homeData{"Home", blockchain.Blockchain().AllBlocks()}
	// templates.ExecuteTemplate(rw, "home", data)
}

func handleAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}


func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", handleHome)
	handler.HandleFunc("/add", handleAdd)
	fmt.Printf("Listening on HTML Explorer http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port), handler))
}