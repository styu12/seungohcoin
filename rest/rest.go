package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

type balanceResponse struct {
	Address string	`json:"address"`
	Balance int	`json:"balance"`
}

type urlDescription struct {
	URL url	`json:"url"`
	Method string	`json:"method"`	
	Description string	`json:"description"`
	Payload string `json:"payload,omitempty"`
}

type addTxPayload struct {
	To string
	Amount int
}

type errorResponse struct {
	ErrorMessage string `json"errorMessage"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: url("/"),
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: url("/status"),
			Method: "GET",
			Description: "See the status of the Blockchain",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "See All Blocks",
		},
		{
			URL: url("/blocks/{hash}"),
			Method: "GET",
			Description: "Find A Block",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "Add New Block",
			Payload: "data:string",
		},
		{
			URL: url("/balance/{address}"),
			Method: "GET",
			Description: "See the balance of an address",
		},
		{
			URL: url("/balance/{address}?total=true"),
			Method: "GET",
			Description: "See the total amount of that balance",
		},
		{
			URL: url("/mempool"),
			Method: "GET",
			Description: "Unconfirmed Transactions in Mempool",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func allBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}	else {
		encoder.Encode(block)
	}
	
}

func status(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	encoder.Encode(blockchain.Blockchain())
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		utils.HandleError(json.NewEncoder(rw).Encode(balanceResponse{address, blockchain.Blockchain().BalanceByAddress(address)}))
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blockchain().UTxOutsByAddress(address)))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"Not Enough Funds."})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", allBlocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening REST API on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}