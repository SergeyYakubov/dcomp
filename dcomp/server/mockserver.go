// +build !release

package server

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		MockFuncGetAll,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}/",
		MockFuncGet,
	},
	utils.Route{
		"PatchJob",
		"PATCH",
		"/jobs/{jobID}/",
		MockFuncGet,
	},

	utils.Route{
		"DeleteJob",
		"DELETE",
		"/jobs/{jobID}/",
		MockFuncDelete,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		MockFuncSubmit,
	},
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		MockFuncEstimate,
	},
	utils.Route{
		"Authorize",
		"POST",
		"/authorize/",
		MockFuncAuthorize,
	},
	utils.Route{
		"ReceiveFile",
		"POST",
		"/jobfile/{jobID}/",
		MockFuncReceiveJobFile,
	},
}

func MockFuncReceiveJobFile(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func MockFuncSubmit(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	if strings.Contains(string(b), "/etc") {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{"ImageName":"submittedimage","Script":"aaa","NCPUs":1,"Id":"578359205e935a20adb39a18","Status":1}`)
	}

}

func MockFuncEstimate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"local":100,"slurm":0,"cloud":10,"batch":0}`)
}

func MockFuncGetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	showFinished := r.URL.Query().Get("finished")
	if showFinished == "true" {
		fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a19"}]`)
		return
	}
	fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a18"}]`)
}

func MockFuncAuthorize(w http.ResponseWriter, r *http.Request) {
	var t AuthorizationRequest

	d := json.NewDecoder(r.Body)
	d.Decode(&t)

	authType, user, err := SplitAuthToken(t.Token)

	if err != nil {
		http.Error(w, "cannot split token", http.StatusUnauthorized)
		return
	}

	if authType != "Basic" {
		http.Error(w, "wrong auth type", http.StatusUnauthorized)
		return
	}

	if user == "wronguser" {
		http.Error(w, "user not allowed", http.StatusUnauthorized)
		return
	}

	var tt AuthorizationResponce
	tt.UserName = user

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(tt)

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())

}

func MockFuncGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	jobID := vars["jobID"]
	jobLog := r.URL.Query().Get("log")

	logCompress := r.URL.Query().Get("compress")

	if jobLog == "true" {
		w.WriteHeader(http.StatusOK)
		if logCompress == "true" {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			fmt.Fprintf(gzr, `hello compressed`)
		} else {
			fmt.Fprintf(w, `hello`)
		}
		return
	}

	if jobID == "678359205e935a20adb39a18" {
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `{"Status":310}`)
		return
	} else if jobID == "878359205e935a20adb39a18" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `[{"Status":310}]`)
		return
	} else if jobID == "778359205e935a20adb39a18" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `[{"Status":103}]`)
		return
	} else if jobID != "578359205e935a20adb39a18" {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a18"}]`)
}

func MockFuncDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	if jobID != "578359205e935a20adb39a18" {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ProcessMockBasicAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		at, au, err := ExtractAuthInfo(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if at != "Basic" {
			http.Error(w, "wrong auth type", http.StatusUnauthorized)
			return
		}

		if au == "wronguser" {
			http.Error(w, "user not allowed", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func CreateMockServer(srv *Server) *httptest.Server {
	var ts *httptest.Server
	mux := utils.NewRouter(listRoutes)
	auth := srv.GetAuth()
	var newsrv func(http.Handler) *httptest.Server
	if srv.Tls {
		newsrv = httptest.NewTLSServer
	} else {
		newsrv = httptest.NewServer
	}
	switch auth := auth.(type) {
	case nil:
		ts = newsrv(http.HandlerFunc(mux.ServeHTTP))
	case *HMACAuth:
		ts = newsrv(ProcessHMACAuth(http.HandlerFunc(mux.ServeHTTP), auth.Key))
	case *BasicAuth:
		ts = newsrv(ProcessMockBasicAuth(http.HandlerFunc(mux.ServeHTTP)))
	case *JWTAuth:
		ts = newsrv(ProcessJWTAuth(http.HandlerFunc(mux.ServeHTTP), auth.Key))

	}
	srv.parseUrl(ts.URL)
	return ts
}
