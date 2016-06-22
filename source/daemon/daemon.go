package daemon

import (
	"fmt"
	"log"
	"net/http"
)

func StartDaemon(args []string) {

	mux := NewRouter()
	//	mux.Schemes("https")
	fmt.Println("hello")

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
