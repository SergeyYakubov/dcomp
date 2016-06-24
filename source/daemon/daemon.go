package daemon

import (
	"log"
	"net/http"
)

func StartDaemon(args []string) {

	mux := NewRouter()
	//	mux.Schemes("https")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
