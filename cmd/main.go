package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	bootsrap "github.com/juanjoaquin/back-g-user/internal/pkg"
	"github.com/juanjoaquin/back-g-user/internal/user"
)

func main() {
	router := mux.NewRouter()
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

	userRepository := user.NewRepo(l, db) // Importamos el Logger (l)

	// Al haber hecho lo de la capa de servicio. Va a necesitar recibir un servicio, nosotros debemos especificarlo
	userService := user.NewService(l, userRepository) // Este userService se lo debemos pasar al endpoint. En este caso, le pasamos el repository // Importamos el Logger (l)
	userEndpoint := user.MakeEndpoints(userService, user.Config{LimPageDef: pagLimDef})

	router.HandleFunc("/users", userEndpoint.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoint.Get).Methods("GET") // La rutas dinamicas se usan con /{"Nombre de lo que deseamos dinamico"}

	router.HandleFunc("/users", userEndpoint.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEndpoint.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoint.Delete).Methods("DELETE")

	// Obtenemos el puerto a traves de la ENV, y no hardcodeado
	port := os.Getenv("PORT")
	// Generamos una address
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
