package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	functionservice "github.com/nilspolek/AstralFS/function-service"
	proxyservice "github.com/nilspolek/AstralFS/proxy-service"
)

type REST struct {
	Router *mux.Router
	fns    functionservice.FunctionService
	ps     proxyservice.ProxyService
}

type IDDTO struct {
	Id uuid.UUID `json:"id"`
}

type CreateResponseDTO struct {
	Status string    `json:"status"`
	Id     uuid.UUID `json:"id"`
	Path   string    `json:"path"`
}

func New(router *mux.Router, fns functionservice.FunctionService, ps proxyservice.ProxyService) REST {
	var rest REST

	rest.fns = fns
	rest.ps = ps

	router.HandleFunc("/function", rest.CreateFunction).Methods("POST")
	router.HandleFunc("/function", rest.DeleteFunction).Methods("DELETE")
	router.HandleFunc("/function", rest.GetFunctions).Methods("GET")

	rest.Router = router
	return rest
}

func (rest *REST) GetFunctions(w http.ResponseWriter, r *http.Request) {
	fns, err := rest.fns.GetFunctions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(fns)

}

func (rest *REST) CreateFunction(w http.ResponseWriter, r *http.Request) {
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
