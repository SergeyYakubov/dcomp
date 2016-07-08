package daemon

import (
	"log"
	"net/http"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

func Start(args []string) {
	mux := utils.NewRouter(ListRoutes)
	SetServerConfiguration()
	log.Fatal(http.ListenAndServe(":8000", mux))
}
