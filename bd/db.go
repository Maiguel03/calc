package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	/*dbHost := "dpg-cnn7jl6n7f5s73da6f70-a/calc_t4wj"
	dbPort := "5432"
	dbUser := "dek"
	dbPassword := "rubqwjCMeRfKkwvGsGeYudqOpmKtbo34"
	dbName := "calc_t4wj"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s sslmode=require",
		dbHost, dbPort, dbUser, dbName, dbPassword)*/

	db, err := sql.Open("postgres", "postgres://dek:rubqwjCMeRfKkwvGsGeYudqOpmKtbo34@dpg-cnn7jl6n7f5s73da6f70-a/calc_t4wj")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Loguear(usuario string, contraseña string) bool {


	db, err := ConnectDB()
	defer db.Close()


	if err != nil {
		log.Fatal("Error conectando a la base de datos", err)
	}
	log.Println("Conexion exitosa")

	rows, err := db.Query("SELECT * FROM username where username = $1  and pass = $2;", usuario, contraseña)

	if err != nil {
		log.Fatal("Ocurrio un error al hacer la consulta", err)
	}

	for rows.Next() {
		var user string
		var pass string

		// Escanea los valores de las columnas en las variables
		err := rows.Scan(&user, &pass)

		if err != nil {
			log.Fatal("Error al escanear fila:", err)
		}

		if user == usuario && pass == contraseña {
			return true
		}
	}

	return false
}

