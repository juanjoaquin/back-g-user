package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	bootsrap "github.com/juanjoaquin/back-g-user/internal/pkg"
	"github.com/juanjoaquin/back-g-user/internal/pkg/handler"
	"github.com/juanjoaquin/back-g-user/internal/user"
)

func main() {
	// router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootsrap.InitLogger()

	db, err := bootsrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}
	// Debemos definir el Context para pasarselo al handler
	ctx := context.Background()

	userRepository := user.NewRepo(l, db) // Importamos el Logger (l)

	// Al haber hecho lo de la capa de servicio. Va a necesitar recibir un servicio, nosotros debemos especificarlo
	userService := user.NewService(l, userRepository) // Este userService se lo debemos pasar al endpoint. En este caso, le pasamos el repository // Importamos el Logger (l)
	/* 	userEndpoint := user.MakeEndpoints(userService, user.Config{LimPageDef: pagLimDef})
	 */

	// Generamos el handler. Que sera la funcion de NewUserHTTPServer
	handler := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userService, user.Config{LimPageDef: pagLimDef}))

	/* 	router.HandleFunc("/users", userEndpoint.GetAll).Methods("GET")
	   	router.HandleFunc("/users/{id}", userEndpoint.Get).Methods("GET") // La rutas dinamicas se usan con /{"Nombre de lo que deseamos dinamico"}

	   	// router.HandleFunc("/users", userEndpoint.Create).Methods("POST")
	   	router.HandleFunc("/users/{id}", userEndpoint.Update).Methods("PATCH")
	   	router.HandleFunc("/users/{id}", userEndpoint.Delete).Methods("DELETE") */

	// Obtenemos el puerto a traves de la ENV, y no hardcodeado
	port := os.Getenv("PORT")
	// Generamos una address
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(handler), // Aqui pasamos la funcion de abajo accessControl que ya devuelve el Handler
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Definimos un canal donde ya generamos el server. Vamos a ejecutar una Go Routine.
	errCh := make(chan error)
	go func() {
		l.Println("listn in:", address)
		errCh <- srv.ListenAndServe() // Si hay un error nos marca el error en el canal
	}()
	// Especificamos el error que se ha generado:
	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.ListenAndServe())

}

/*  Definimos una funcion por tema de CORS. Para habilitar todas las opciones como permitidas */
func accessControl(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)

			return
		}
		handler.ServeHTTP(w, r)
	})

}
