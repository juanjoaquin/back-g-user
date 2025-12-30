package user

import (
	"log"

	"github.com/juanjoaquin/back-g-domain/domain"
)

// Nuestro servicio lo vamos a menajar con interfaces. Esto nos facilitara para mockearlo, o utilizarlo de forma mas generica
type Service interface {
	/* 	1. Vamos a definirle los metodos de los Endpoints que fuimos utilizando.
	   	Le pasaremos tambien los elementos del body del Create por ejemplo */
	Create(firstName, lastName, email, phone string) (*domain.User, error)
	GetAll(filters Filters, offset, limit int) /* Pasamos el Filtrado de params */ ([]domain.User, error) // Get All
	Get(id string) (*domain.User, error)                                                                  // Get by User ID
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
}

// Struct de Filter params:
type Filters struct {
	FirstName string
	LastName  string
}

/* 2. Vamos a definir una struct, está sera en privado */
type service struct {
	log *log.Logger
	// Ahora debemos pasar el Repository
	repo Repository
}

/*
 3. Haremos una funcion llamada: NewService
    Esta lo que hara sera crear un nuevo servicio, que esta ser la interface.
*/
func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

/* 4. Vamos a generar un metodo, que esto se lo deberemos pasar a la funcion de NewService. */
func (s service) Create(firstName, lastName, email, phone string) (*domain.User, error) {

	s.log.Println("Create user service")

	/* Ahora para crear el endpoint, pasamos los valores que tenemos del User */
	user := domain.User{
		// Una vez terminado esto, se lo debemos pasar al Repositorio
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	/* Una vez pasado el repo, dentro de nuestro create. Debemos pasarle el repository. Debemos ejecutar el metodo Create del propio Repo */
	/* Una vez creado el User en el Repositorio, debemos hacer una validacion de que si el Repo da error, este service lo handlea */
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

/* Get All de los Users */
func (s service) GetAll(filters Filters, offset, limit int) /* Pasamos el Search Params */ ([]domain.User, error) {

	/* Traemos a los Users y usamos nos traemos el .GetAll() de la Interface del Service (s.repo), que previamente declaramos en nuestro Repository (GetAll) */
	users, err := s.repo.GetAll(filters, offset, limit) // Tambien le pasamos el Search Params

	// Handleo error
	if err != nil {
		return nil, err
	}

	// Returno la entidad completa con el nil
	return users, nil

}

func (s service) Get(id string) (*domain.User, error) {
	user, err := s.repo.Get(id)

	// Handleo error
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	return s.repo.Update(id, firstName, lastName, email, phone)
}

// Pasamos el Count en el Service
func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
