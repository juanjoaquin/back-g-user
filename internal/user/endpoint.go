// Aqui vamos a generar nuestros endpoints
// El Endpoint seria el equivalente al Controller

// 1. Vamos a crear una funcion llamada "MakeEndpoints". Esta se encargara de crear nuestros endpoints
// 2. Vamos a crear una struct, que va a tener todos los endpoints que nosotros vayamos a utilizar
package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/juanjoaquin/back-g-meta/pkg/meta"
	"github.com/juanjoaquin/back-g-response/response"
)

type (
	// Usamos el Endpoint de GoKit recomendable
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		// Aqui definimos los endpoints:
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	/* 	4. Vamos a definir nuestro request para arrancar.
	   	Vamos a crear una Struct donde vamos a tener los campos que vamos a recibir.
		Esto se lo debemos pasar al controlador que tenemos en el Create, de abajo.

	*/
	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	/*
		5. Vamos a generar un struct para los errores de las response:
		DEPRECADO
	*/
	/* 	ErrorRes struct {
		Error string `json:"error"`
	} */

	/* Nueva Struct  para el UPDATE */

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	/* Struct de Response general */

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"` // Esto es una interface porque le podemos mandar cualquier cosa relacionada a nuestro servicio (usuario, curso, etc).
		Err    string      `json:"err,omitempty"`  // Usamos el omitempty, esto que si viene vacio lo omite. No lo recibe.
		Meta   *meta.Meta  `json:"meta,omitempty"` // Devolvemos en la response el Meta (Total pages)
	}

	Config struct {
		LimPageDef string
	}
)

// 3. Esta es la función de MakeEndpoints, que va a devolver una estructura de Edpoints. Estos son los que vamos a poder utilizar en nuestro dominio.

// Ahora le pasaremos el Service. Este lo tendra como prop. También lo recibira todas las funciones que encapsula.
func MakeEndpoints(s Service, config Config) Endpoints {
	// Returnamos los endpoints
	return Endpoints{
		// Debemos indicar que cada endpoint representa cada funcion
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

// Estas seran una funcion privada, ya que empiezan con minuscula, porque el que vamos a usar es el de arriba
func makeDeleteEndpoint(s Service) Controller {
	// Definimos la funcion del Controller, que seria la que esta arriba de todo del Controller
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Aqui ira nuestra logica del endpoint

		// Es parecido al de Get By Id.

		// Usamos el path de Gorilla Mux
		path := mux.Vars(r)
		// Le pasamos el id
		id := path["id"]

		// Nos traemos el service.Delete y handleamos el error (CON LA NUEVA STRUCT)
		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Error al borrar el usuario"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})
	}
}

// Create Endpoint
// Aqui tambien le pasaremos ese servicio
func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		// Asignamos el nuevo valor con el Go Kit del Request del Context
		req := request.(CreateReq)

		/* 		var req CreateReq // Este es el viejo req para el CreateReq*/

		// Esto no lo usamos mas porque con el Middleware, es donde aplicamos esto. Se modifica

		/* 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		   			w.WriteHeader(400)
		   			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
		   			return
		   		}
		*/
		// La validacion tambien se modifica. No usamos mas esta

		/* 		if req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "First name is required"})
			return
		}
		*/

		//Esta es el nuevo tipo de validacion con nuestro Package. Especificamos cual es el tipo de error
		if req.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.LastName == "" {
			return nil, response.BadRequest("Last name is required")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Phone, req.Email) // Le pasamos el Context (ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		// Y aqui no lo usamos mas.
		/* 		json.NewEncoder(w).Encode(&Response{Status: 201, Data: user}) */

		// Debemos hacer un return de nuestro Package de Response
		return response.Created("success", user, nil), nil
	}
}

// Get All Endpoint
func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		//Para obtener el Query Params
		v := re.URL.Query()
		// Nos traemos el SearchParams, la Struct del Service.
		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
		}

		// Nos traemos el Limit y el Page desde las ENV.
		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		// Aqui aplicamos el Counter que hicimos despues de todo esto
		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}
		// Nos traemos el Package de Meta de la función New del propio package
		meta, err := meta.New(page, limit, count, config.LimPageDef) // Le debemos pasar tanto Page & Limit
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// Debemos hacer referencia al GetAll del Service
		users, err := s.GetAll(filters, meta.Offset(), meta.Limit()) // Pasamos el filtro al GetAll del Service. Y tambien el Meta de Offset y Limit

		// Si el error es != nill, manejamos con el w.WirteHeader la Bad Request
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		// Lo devolvemos con la nueva struct de Response & Devolvemos el package de Meta (previamente traido arriba)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: users, Meta: meta})
	}
}

// Get by id endpoint
func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// Se debe crear una variable y guardar el ID como parametro

		//Gorilla Max con Vars le pasamos nuestra request, y esta nos devuelve un path con los parametros
		path := mux.Vars(r)    // Aqui llamamos a mux (Gorilla Mux),Vars(r) / La r es el http.Request como parametro que tenemos
		id := path["id"]       // Especificamos que queremos el ID
		user, err := s.Get(id) // Declaramos al user, y llamamos al service ( s.Get() )

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Usuario no encontrado"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})

	}
}

// Update endpoint
func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		// Llamamos a la struct que creamos previamente
		var req UpdateReq

		// Decodificamos el body y lo validamos
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request ofrmat"})
			return
		}

		// Validamos los campos y si viene vacio
		if req.FirstName != nil && *req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "First name is required"})
			return
		}

		if req.LastName != nil && *req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Last name is required"})
			return
		}

		if req.Email != nil && *req.Email == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Email is required"})
			return
		}

		if req.Phone != nil && *req.Phone == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Phone is required"})
			return
		}

		// Y aca debemos hacer lo del Gorilla Mux
		path := mux.Vars(r)
		id := path["id"]

		// Vamos a returnar la capa de Servicio que tenemos. Pasandole el ID. En este caso sería: s.Update() con el Body que le habiamos pasado.
		if err := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User dosent exists"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "User updated succesfully"})
	}
}
