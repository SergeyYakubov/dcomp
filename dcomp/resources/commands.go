package resources

import (
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func (p *plugin) SubmitJob(w http.ResponseWriter, r *http.Request) {

	var t structs.JobInfo

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	res, err := p.resource.SubmitJob(t.JobDescription)
	if err != nil {
		http.Error(w, "cannot submit job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = p.database.CreateRecord(t.Id, res)
	if err != nil {
		http.Error(w, "cannot create record job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//	b := new(bytes.Buffer)
	//	json.NewEncoder(b).Encode(res)
	//	w.Write(b.Bytes())
}
