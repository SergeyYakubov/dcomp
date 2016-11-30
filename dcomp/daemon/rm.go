package daemon

import (
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"net/http"
)

func deleteJobInResourceIfNeeded(job structs.JobInfo) error {
	if job.Status == structs.StatusWaitData {
		return nil
	}

	res := resources[job.Resource]
	_, err := res.Server.CommandDelete("jobs" + "/" + job.Id)

	return err

}

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := deleteJobInResourceIfNeeded(job); err != nil {
		http.Error(w, "cannot delete job in resource: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := db.DeleteRecordByID(job.Id); err != nil {
		http.Error(w, "cannot delete job: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}
