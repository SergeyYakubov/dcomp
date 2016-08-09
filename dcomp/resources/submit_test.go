package resources

import (
	"testing"

	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type request struct {
	job        structs.JobInfo
	path       string
	cmd        string
	answercode int
	answer     string
	message    string
}

var submitTests = []request{
	{structs.JobInfo{JobDescription: structs.JobDescription{"image", "script", 20},
		Id: "578359205e935a20adb39a18"}, "POST", "jobs", http.StatusCreated, "12345", "submit job"},
}

func TestSubmitJob(t *testing.T) {
	/*p := resources.NewPlugin(new(resources.MockResource), new(database.Mockdatabase))
	mux := utils.NewRouter(p.listRoutes)
	for _, test := range submitTests {

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.job.ImageName == "nil" {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8002/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answercode, w.Code, test.message)
		if w.Code == http.StatusOK {
			assert.Contains(t, w.Body, test.answer, test.message)
		}
	}*/

}
