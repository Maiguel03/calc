// main.go
package main

import (
	db "calc/bd"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	//"strconv"
	"time"
)

type Verficar struct {
	LoginFallido bool
}

type Historico struct {
	Fecha     string
	Operacion string
}

func RecogerHistorial() []Historico {
	db, err := db.ConnectDB()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM historico")

	defer rows.Close()

	var Resultados []Historico
	var r Historico

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

var templates = template.Must(template.New("T").ParseGlob("static/*.html"))

func renderTemplate(rw http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(rw, name, data)
	if err != nil {
		panic(err)
	}
}

// Funcion para recoger historial de operaciones

func llenarHistorial(resultado string, fecha string) error {
	fmt.Println("LLenar Historial: ", resultado, fecha)
	db, err := db.ConnectDB()

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

// Función para calcular el resultado de una operación
func calcular(operando1, operando2 float64, operador string) float64 {
	var resultado float64
	fmt.Println(operador)
	switch operador {
	case "+":
		resultado = operando1 + operando2
		resultado = math.Round(resultado*100000000) / 100000000
	case "-":
		resultado = operando1 - operando2
		resultado = math.Round(resultado*100000000) / 100000000
	case "*":
		resultado = operando1 * operando2
		resultado = math.Round(resultado*100000000) / 100000000
	case "/":
		resultado = operando1 / operando2
		resultado = math.Round(resultado*100000000) / 100000000

	case "sqrt":
		fmt.Println(operando1, operador)
		resultado = math.Sqrt(operando1)
		resultado = math.Round(resultado*100000000) / 100000000

	case "!":
		var factorial int64
		var numero int64
		numero = int64(operando1)
		factorial = 1
		fmt.Println(operador, operando1, numero)
		for i := int64(1); i <= numero; i++ {
			factorial *= i
			fmt.Println(factorial)
		}

		resultado = float64(factorial)
		fmt.Println(resultado)
		resultado = math.Round(resultado*100000000) / 100000000

	case "log":
		resultado = math.Log10(operando1)
		resultado = math.Round(resultado*100000000) / 100000000

	case "^":
		resultado = math.Pow(operando1, operando2)
		resultado = math.Round(resultado*100000000) / 100000000
	}

	return resultado
}

// Función para manejar las peticiones del cliente
func calculadora(w http.ResponseWriter, r *http.Request) {
	// Obtener los valores de los parámetros de la URL
	operando1 := r.URL.Query().Get("operando1")
	operando2 := r.URL.Query().Get("operando2")
	operador := r.URL.Query().Get("operador")
	timestamp := r.URL.Query().Get("timestamp")
	result := r.URL.Query().Get("resultado")
	fmt.Println("timestamp:", timestamp, "resultado:", result)
	//result := resultado
	fmt.Println(timestamp, result)
	// Si hay valores, convertirlos a números y calcular el resultado
	if result != "" && timestamp != "" {
		err := llenarHistorial(result, timestamp)
		if err != nil {
			panic(err)
		}
	}
	if operando1 != "" && operando2 != "" && operador != "" {

		fmt.Println(timestamp, result)

		num1, err1 := strconv.ParseFloat(operando1, 64)
		num2, err2 := strconv.ParseFloat(operando2, 64)
		if err1 == nil && err2 == nil {
			resultado := calcular(num1, num2, operador)
			// Enviar el resultado como respuesta
			w.Write([]byte(strconv.FormatFloat(resultado, 'f', -1, 64)))
		} else {
			// Enviar un mensaje de error como respuesta
			w.Write([]byte("Error: los operandos deben ser números"))
		}
	} else {
		// Si no hay valores, servir el archivo index.html
		Resultados := RecogerHistorial()
		renderTemplate(w, "index.html", Resultados)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var logueado bool
	if r.Method == "POST" {
		usuario := r.FormValue("usuario")
		contraseña := r.FormValue("contrasena")
		fmt.Println(usuario, contraseña)

		logueado = db.Loguear(usuario, contraseña)

		if logueado {
			http.Redirect(w, r, "/calculadora", http.StatusSeeOther)
		}

		Login := Verficar{LoginFallido: true}
		renderTemplate(w, "index.html", Login)
		return

	}
	renderTemplate(w, "index.html", nil)
}

func main() {
	// Crear un servidor web en el puerto 8080
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", Login)
	http.HandleFunc("/calculadora", calculadora)

	http.ListenAndServe(":8080", nil)
}
