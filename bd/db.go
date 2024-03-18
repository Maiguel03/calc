package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"calc/modelo"
	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	dbHost := "dpg-cnn7jl6n7f5s73da6f70-a.oregon-postgres.render.com"
	dbPort := "5432"
	dbUser := "dek"
	dbPassword := "rubqwjCMeRfKkwvGsGeYudqOpmKtbo34"
	dbName := "calc_t4wj"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s sslmode=require",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Loguear(usuario string, contraseña string) bool {

	db, err := connectDB()
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

// Funcion para recoger historial de operaciones
func RecogerHistorial() []modelo.Historico {
	db, err := connectDB()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM historico")

	defer rows.Close()

	var Resultados []modelo.Historico
	var r modelo.Historico

	for rows.Next() {
		err := rows.Scan(&r.Operacion, &r.Fecha)
		if err != nil {
			panic(err)
		}

		fecha := r.Fecha

		f, err := time.Parse(time.RFC3339, fecha)
		if err != nil {
			panic(err)
		}

		//Obtenemos la zona horaria local
		zonaHoraria, err := time.LoadLocation("America/New_York")
		if err != nil {
			panic(err)
		}
		fechaCorregida := f.In(zonaHoraria)
		formatoPersonalizado := "2006-01-02 15:04:05"
		fc := fechaCorregida.Format(formatoPersonalizado)
		r.Fecha = fc
		Resultados = append(Resultados, r)
		fmt.Println(Resultados)

	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return Resultados
}
//Funcion para enviar un registro al historial de operaciones
func LlenarHistorial(resultado string, fecha string) error {
	fmt.Println("LLenar Historial: ", resultado, fecha)
	db, err := connectDB()

	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO historico (operacion, fecha) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Ejecuta la sentencia con los valores deseados
	_, err = stmt.Exec(resultado, fecha)
	if err != nil {
		return err
	}

	fmt.Println("Registro insertado correctamente")
	return nil
}
