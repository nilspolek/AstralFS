package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	functionservice "github.com/nilspolek/AstralFS/function-service"
	proxyservice "github.com/nilspolek/AstralFS/proxy-service"
	"github.com/nilspolek/AstralFS/repo"
	"github.com/nilspolek/goLog"
)

type REST struct {
	Router *mux.Router
	fns    functionservice.FunctionService
	ps     proxyservice.ProxyService
	Repo   *repo.Repo
}

type IDDTO struct {
	Id uuid.UUID `json:"id"`
}

type CreateResponseDTO struct {
	Status string    `json:"status"`
	Id     uuid.UUID `json:"id"`
	Path   string    `json:"path"`
}

func New(router *mux.Router, fns functionservice.FunctionService, ps proxyservice.ProxyService, repo *repo.Repo) REST {
	var rest REST

	rest.fns = fns
	rest.ps = ps

	router.HandleFunc("/function", rest.CreateFunction).Methods("POST")
	router.HandleFunc("/function", rest.DeleteFunction).Methods("DELETE")
	router.HandleFunc("/function", rest.GetFunctions).Methods("GET")
	router.HandleFunc("/functionAll", rest.DeleteAllFunctions).Methods("DELETE")

	rest.Router = router
	rest.Repo = repo
	return rest
}

func (rest *REST) GetFunctions(w http.ResponseWriter, r *http.Request) {
	goLog.Info("GET Request to /functions")
	fns, err := rest.fns.GetFunctions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(fns)
}

func (rest *REST) DeleteAllFunctions(w http.ResponseWriter, r *http.Request) {
	goLog.Info("DELETE Request to /functionsAll")
	fns, err := (*rest.Repo).GetFunctions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, fn := range fns {
		rest.ps.DeleteRoute(proxyservice.Route{
			Path: fn.Route,
		})
		rest.fns.DeleteFunction(fn.Id)
	}
	w.WriteHeader(http.StatusOK)
}

func (rest *REST) CreateFunction(w http.ResponseWriter, r *http.Request) {
	goLog.Info("POST Request to /functions")
	var (
		createFn functionservice.Function
		err      error
		port     int
	)
	createFn.Id = uuid.New()
	json.NewDecoder(r.Body).Decode(&createFn)

	port, err = rest.fns.CreateFunction(createFn)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err = rest.ps.AddRoute(proxyservice.Route{
		Path:   createFn.Route,
		Target: "http://localhost:" + strconv.Itoa(port),
	}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)

	resp := CreateResponseDTO{
		Status: "Ok",
		Path:   createFn.Route,
		Id:     createFn.Id,
	}
	json.NewEncoder(w).Encode(resp)
}

func (rest *REST) DeleteFunction(w http.ResponseWriter, r *http.Request) {
	goLog.Info("DELETE Request to /functions")
	var (
		id    IDDTO
		route string
	)
	json.NewDecoder(r.Body).Decode(&id)

	fns, err := rest.fns.GetFunctions()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, fn := range fns {
		if fn.Id == id.Id {
			route = fn.Route
			break
		}
	}

	if err = rest.fns.DeleteFunction(id.Id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err = rest.ps.DeleteRoute(proxyservice.Route{
		Path: route,
	}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
