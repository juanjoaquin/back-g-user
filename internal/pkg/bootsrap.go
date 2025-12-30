package bootsrap

// Basicamente Bootsrap en go es un modulo de arranque, para centralizarlo todo en la app.
import (
	"fmt"
	"log"
	"os"

	"github.com/juanjoaquin/back-g-domain/domain" // Hay que hacer un go get con el link del repo
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Esta funcion será la conexión de la DB. Que lo traemos del package de GORM.
func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		// Con el elemento de: os. Es donde nos emparejamos a las ENV
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	/* Para la conexion a la DB, debemos usar el gorm package
	Con la funcion Open, y el package mysql
	*/
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	/* IMPORTANTE: ACTIVAR ESTAS VARIABLES DE ENTORNO EN LA .ENV */
	// Este es el DEBUG de la DB en caso de que venga en true
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	// Debemos hacer el Auto Migrate a traves de las variables de entorno
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			return nil, err
		}
	}

	return db, nil

}

// Esta funcion es el Logger de la DB
func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
