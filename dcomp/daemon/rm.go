package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]

	var jobs []structs.JobInfo
	if err := db.GetRecordsByID(jobID, &jobs); err != nil || len(jobs) != 1 {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusNotFound)
		return
	}

	job := jobs[0]
	res := resources[job.Resource]
	_, err := res.Server.CommandDelete("jobs" + "/" + job.Id)

	if err != nil {
		http.Error(w, "cannot delete job in resource: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := db.DeleteRecordByID(jobID); err != nil {
		http.Error(w, "cannot delete job: "+err.Error(), http.StatusNotFound)
	}
}
