package daemon

import (
	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"net/http"
)

func deleteJobInResourceIfNeeded(job structs.JobInfo) error {
	if job.Status == structs.StatusWaitData {
		return nil
	}

	if job.Resource == "" {
		return nil
	}

	res := resources[job.Resource]

	b, status, err := res.Server.CommandDelete("jobs" + "/" + job.Id)

	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New(b.String())
	}

	return nil

}

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// do not check for errors, delete job in database anyway
	deleteJobInResourceIfNeeded(job)

	if err := db.DeleteRecordByID(job.Id); err != nil {
		http.Error(w, "cannot delete job: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
