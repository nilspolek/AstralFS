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

	// fn := functionservice.Function{
	// 	Id:    uuid.New(),
	// 	Image: "nilspolek/echo-server",
	// 	Port:  8080,
	// }

	// fs, err := functionservice.New()

	// if err != nil {
	// 	goLog.Error(err.Error())
	// }

	// router := mux.NewRouter()
	// ps := proxyservice.New(router)
	// port, err := fs.CreateFunction(fn)

	// ps.AddRoute(proxyservice.Route{
	// 	Target: "http://localhost:" + strconv.Itoa(port),
	// 	Path:   "/test",
	// })

	// if err != nil {
	// 	goLog.Error(err.Error())
	// }

	// err = http.ListenAndServe(":8090", router)

	// if err != nil {
	// 	goLog.Error("%v", err)
	// }
}
