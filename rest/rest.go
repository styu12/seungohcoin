package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/styu12/seungohcoin/blockchain"
	"github.com/styu12/seungohcoin/utils"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) { 
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL url	`json:"url"`
	Method string	`json:"method"`	
	Description string	`json:"description"`
	Payload string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: url("/"),
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "See All Blocks",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "Add New Block",
			Payload: "data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func allBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utils.HandleError(err)
	block := blockchain.GetBlockchain().GetBlock(height)
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(block)
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", allBlocks).Methods("GET", "POST")
	router.HandleFunc("/block/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}