package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/marc-campbell/nicedishy/pkg/handlers"
)

func Start() {
	log.Printf("nicedishy api version %s\n", os.Getenv("VERSION"))

	r := mux.NewRouter()

	r.Methods("OPTIONS").HandlerFunc(handlers.CORS)

	/**********************************************************************
	* Unauthenticated routes
	**********************************************************************/

	r.HandleFunc("/healthz", handlers.Healthz)
	r.HandleFunc("/api/v1/login", handlers.Login)
	r.Path("/api/v1/login/callback").Methods("POST").HandlerFunc(handlers.LoginCallback)

	/**********************************************************************
	* Static routes
	**********************************************************************/

	/**********************************************************************
	* Authenticated routes
	**********************************************************************/
	sessionAuthQuietRouter := r.PathPrefix("").Subrouter()
	sessionAuthQuietRouter.Use(handlers.RequireValidSessionQuietMiddleware)
	sessionAuthQuietRouter.Path("/api/v1/dishy").Methods("POST").HandlerFunc(handlers.CreateDishy)
	sessionAuthQuietRouter.Path("/api/v1/dishies").Methods("GET").HandlerFunc(handlers.ListDishies)
	sessionAuthQuietRouter.Path("/api/v1/dishy/{id}/token").Methods("GET").HandlerFunc(handlers.GetDishyToken)
	sessionAuthQuietRouter.Path("/api/v1/dishy/{id}").Methods("DELETE").HandlerFunc(handlers.DeleteDishy)

	tokenAuthQuietRouter := r.PathPrefix("").Subrouter()
	tokenAuthQuietRouter.Use(handlers.RequireValidTokenQuietMiddleware)
	tokenAuthQuietRouter.Path("/api/v1/stats").Methods("POST").HandlerFunc(handlers.StoreData)

	srv := &http.Server{
		Handler: r,
		Addr:    ":3000",
	}

	fmt.Printf("Starting nicedishy-api on port %d...\n", 3000)

	log.Fatal(srv.ListenAndServe())
}
