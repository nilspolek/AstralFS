package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nilspolek/AstralFS/api/rest"
	functionservice "github.com/nilspolek/AstralFS/function-service"
	proxyservice "github.com/nilspolek/AstralFS/proxy-service"
	"github.com/nilspolek/goLog"
)

func main() {
	router := mux.NewRouter()

	fns, err := functionservice.New()
	if err != nil {
		goLog.Error(err.Error())
	}

	ps := proxyservice.New(router)
	r := rest.New(router, fns, ps)

	http.ListenAndServe(":8080", r.Router)

}
