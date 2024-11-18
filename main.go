package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nilspolek/AstralFS/api/rest"
	functionservice "github.com/nilspolek/AstralFS/function-service"
	proxyservice "github.com/nilspolek/AstralFS/proxy-service"
	"github.com/nilspolek/AstralFS/repo"
	sqliterepo "github.com/nilspolek/AstralFS/repo/sqliteRepo"
	"github.com/nilspolek/goLog"
)

func main() {
	router := mux.NewRouter()

	fns, err := functionservice.New()
	if err != nil {
		goLog.Error(err.Error())
	}
	rep, err := sqliterepo.New("./database.db")
	if err != nil {
		goLog.Error(err.Error())
	}

	fns, err = repo.New(&rep, fns)

	ps := proxyservice.New(router)
	r := rest.New(router, fns, ps, &rep)

	http.ListenAndServe(":8080", r.Router)
}
