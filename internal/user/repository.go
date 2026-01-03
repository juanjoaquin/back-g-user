package user

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/juanjoaquin/back-g-domain/domain" // Hay que hacer un go get con el link del repo
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error                                                                 // Le pasamos como puntero al User
	GetAll(ctx context.Context, filters Filters, offset int, limit int) /* Pasamos el Filtrado */ ([]domain.User, error) // El Get all, nos devuelve un array de usuarios
	Get(ctx context.Context, id string) (*domain.User, error)                                                            // El Get by ID, nos devuelve un ID, y un puntero de User
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(ctx context.Context, filters Filters) (int, error) // Devuelve la cantidad de registros
}

// Esta struct va hacer referencia a la DB de GORM
type repo struct {
	log *log.Logger
	db  *gorm.DB
}

// Creamos una funcion que va a instanciar este repo.

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

// Creamos el Metodo Create
func (repo *repo) Create(ctx context.Context, user *domain.User) error {
	repo.log.Println(user) // Este loguer es el que nosotros le hemos pasado. Es lo que imprimira al pegarle al POST

	/* Esto no lo usamos mas. Le quitamos la responsabilidad al Repository y lo usamos con GORM-HOOKS para crear el ID.
	Ahora se encarga el Dominio */
	/* Nuestro ID es UUID, entonces debemos definirle ese UUID desde esta capa (tenemos que usar el package de google/uuid) */
	// user.ID = uuid.New().String()

	/* Tenemos que hacer del objeto  de "db" el metodo "Create", llamando a nuestra Struct (repo) que le debemos pasar la entidad del User */
	result := repo.db.WithContext(ctx).Create(user) // Aca le pasamos el Context

	// Tenemos 2 tipos de manejos de error. Este en el que le decimos, que si el resultado da error, y es distinto a null que lo tire:

	if result.Error != nil {
		repo.log.Println("[ERROR]-[REPOSITORY]-[CREATE]", result.Error)
		return result.Error
	}

	// O este donde seteamos con la funcion propia en la creacion del User, y no una vez previamente declara como en la primera opcion
	/* 	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	} */

	repo.log.Println("User creado exitosamente", user.ID)

	return nil
}

// Creamo el Metodo Get All
func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var u []domain.User // Declaramos la variable user. Que sera un vector de usuarios

	// Debemos traernos el Model del User
	tx := repo.db.WithContext(ctx).Model(u)
	// Nos traemos el filtrado, y se lo pasamos
	tx = applyFilters(tx, filters)
	// Con GORM especificamos tanto el limit & el offset
	tx = tx.Limit(limit).Offset(offset)

	/* Utilizamos la funcion de nuestro repo, para tener la DB, y ejecutar el metodo "Model"
	Con esto especificamos el Modelo que vamos a utilizar. En este caso el User, con su puntero */
	/* result := repo.db.Model(&u).Order("created_at desc").Find(&u) */ // Le aplicamos un orderBy, y un Find para encontrar el user

	//Ahora con el filtrado, le pasamos directamente el tx.Order, y no como esta arriba, es lo mismo, pero le aplicamos el filtrado
	result := tx.Order("created_at desc").Find(&u)

	// Hanldeamos el error
	if result.Error != nil {
		repo.log.Println("[ERROR]-[REPOSITORY]-[GET-ALL]", result.Error)

		return nil, result.Error
	}

	// Returnamos el user y el nil
	return u, nil
	//////////////////////////////////////
	// ABAJO DE TODO APLICAMOS UNA FUNCION PARA EL APLICADO DE FILTROS
	/////////////////////////////////////
}

// Creamo el Metodo Get By ID
func (repo *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	/* Primero debemos generar una estructura User para poder pasarle el ID a GORM */
	user := domain.User{ID: id}

	/* Para buscar la informacion, utilizamos el .First() con el puntero en el User.  */
	if err := repo.db.WithContext(ctx).First(&user).Error; err != nil {
		repo.log.Println(err)
		return nil, ErrUserNotFound{id}
	} // First es el primer elemento que encuentra

	// Devolvemos al puntero del User, tanto como el nil. No se devuelve el result
	return &user, nil

}

// Creamos el Metodo DELETE
func (repo *repo) Delete(ctx context.Context, id string) error {
	/* Primero debemos generar una estructura User para poder pasarle el ID a GORM */
	user := domain.User{ID: id}

	result := repo.db.WithContext(ctx).Delete(&user)

	// El metodo que se usa es el .DELETE

	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	// Esto se usa solo con RESULT. En caso de que venga con Rows = 0. Lanzamos el mensaje del error.
	if result.RowsAffected == 0 {
		repo.log.Printf("user %s doesnt exists", id)
		return ErrUserNotFound{id}
	}

	// Devolvemos nil. No se devuelve el result
	return nil

}

// Creamos el Metodo UPDATE

func (repo *repo) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = *firstName
	}

	if lastName != nil {
		values["last_name"] = *lastName
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}
	result := repo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	// Esto se usa solo con RESULT. En caso de que venga con Rows = 0. Lanzamos el mensaje del error.
	if result.RowsAffected == 0 {
		repo.log.Printf("user %s doesnt exists", id)
		return ErrUserNotFound{id}
	}

	return nil
}

// FUNCION PARA EL APLICADO DE FILTROS
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.FirstName != "" { // Basicamente que si viene vacio, no pasa nada, y que lo devuelva en lower o uppercase
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName) // Query de GORM para la consulta
	}

	if filters.LastName != "" { // Basicamente que si viene vacio, no pasa nada, y que lo devuelva en lower o uppercase
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName) // Query de GORM para la consulta
	}

	return tx
}

// FUNCION PARA EL CONTADOR DEL REGISTRO
func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err) // Imprimimos posiblemente los errores
		return 0, err
	}
	return int(count), nil
}
