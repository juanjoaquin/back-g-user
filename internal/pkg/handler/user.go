package handler

// Basicamente es donde se handlea el funcionamiento del router. Es decir, de los Endpoints propios. El routeo

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanjoaquin/back-g-response/response"
	"github.com/juanjoaquin/back-g-user/internal/user"
)

// Definimos la funcion. Recibira el Context y los Endpoints definidos.
func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	router := mux.NewRouter()

	//No usamos router.HandleFunc() como estabamos usando. Usaremos Handle de Gorilla Mux
	// Tampoco nos traeremos el userEndpoint. Usaremos el httptransport.NewServer() de Go Kit
	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create), // Debemos hacer una conversion. Llamamos al Endpoint de GO KIT, y lo encapsulamos dentro del endpoints.Create del Controller
		decodeCreateUser,
		encodeResponse,
	)).Methods("POST")
	/* EXPLICACION DE LOS PARAMETROS Y LA FUNCIONES:
	El handle primero va a enviar al Decode el POST para crear el usuario. Este Decode ejecuta la funcion, y hace la conversion.
	En caso de que no puede generara un error. Si esta OK, enviara la Request 200.

	Despues pasa el encodeResponse donde recibe la respuesta, y accede a la creacion y al status 200.
	*/

	return router
}

// Esta funcion se encarga de hacer un Decode dentro del request cuando nosotros hagamos el store de un User
func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Definimos la Request del CreateReq
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
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
