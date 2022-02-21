package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/styu12/seungohcoin/explorer"
	"github.com/styu12/seungohcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 승오코인! \n\n")
	fmt.Printf("Please use the following flags: \n\n")
	fmt.Printf("-port:   Set port of the server.\n")
	fmt.Printf("-mode:   Choose between 'rest' and 'html'\n")
	os.Exit(0)
}


func Start() {
	if len(os.Args) < 2 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server.")
	mode := flag.String("mode", "rest", "Choose between 'rest' and 'html', or just 'all'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	case "all":
		go explorer.Start(*port)
		rest.Start(*port + 1000)
	default:
		usage()
	}
}