package daemon

import "net/http"

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := resources[job.Resource]
	_, err = res.Server.CommandDelete("jobs" + "/" + job.Id)

	if err != nil {
		http.Error(w, "cannot delete job in resource: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := db.DeleteRecordByID(job.Id); err != nil {
		http.Error(w, "cannot delete job: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}
