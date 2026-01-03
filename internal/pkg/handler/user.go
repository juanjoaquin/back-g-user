package handler

// Basicamente es donde se handlea el funcionamiento del router. Es decir, de los Endpoints propios. El routeo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanjoaquin/back-g-response/response"
	"github.com/juanjoaquin/back-g-user/internal/user"
)

// Definimos la funcion. Recibira el Context y los Endpoints definidos.
func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	router := mux.NewRouter()

	// Manejo de Errores con Go Kit
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	//No usamos router.HandleFunc() como estabamos usando. Usaremos Handle de Gorilla Mux
	// Tampoco nos traeremos el userEndpoint. Usaremos el httptransport.NewServer() de Go Kit
	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create), // Debemos hacer una conversion. Llamamos al Endpoint de GO KIT, y lo encapsulamos dentro del endpoints.Create del Controller
		decodeCreateUser,
		encodeResponse,
		opts..., // Tambien le pasamos el OPTS del Middleware para descrifar los errores
	)).Methods("POST")
	/* EXPLICACION DE LOS PARAMETROS Y LA FUNCIONES:
	El handle primero va a enviar al Decode el POST para crear el usuario. Este Decode ejecuta la funcion, y hace la conversion.
	En caso de que no puede generara un error. Si esta OK, enviara la Request 200.

	Despues pasa el encodeResponse donde recibe la respuesta, y accede a la creacion y al status 200.
	*/

	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUsers,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

// Esta funcion se encarga de hacer un Decode dentro del request cuando nosotros hagamos el store de un User
func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Definimos la Request del CreateReq
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error())) // Le pasamos el package del Response
	}

	return req, nil
}

// Hacemos un Enconde del Response.
// Esto lo que va a devolver despues el Endpoint una vez que retorne
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	// Hacemos un reconverse de nuestro Package de Response
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())       // Esto tambien
	return json.NewEncoder(w).Encode(r) // Retornamos el response
}

// Aqui pasara por otra instancia donde decodifica el Error. En caso de haber un error por ejemplo un 400. Lo descifra.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response) // Hacemos una conversion del Error al Response
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp) // No debemos hacer un return. Solo mapearle al response que recibimos por parametro, lo que queremos retornar al cliente

}

func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format '%v'", err.Error()))
	}
	path := mux.Vars(r)
	req.ID = path["id"]
	return req, nil

}

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	req := user.DeleteReq{
		ID: path["id"],
	}

	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}
	return req, nil
}

func decodeGetAllUsers(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil

}
