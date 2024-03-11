// main.go
package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

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
func handler(w http.ResponseWriter, r *http.Request) {
	// Obtener los valores de los parámetros de la URL
	operando1 := r.URL.Query().Get("operando1")
	operando2 := r.URL.Query().Get("operando2")
	operador := r.URL.Query().Get("operador")
	fmt.Println(operador)
	// Si hay valores, convertirlos a números y calcular el resultado
	if operando1 != "" && operando2 != "" && operador != "" {
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
		http.ServeFile(w, r, "index.html")
	}
}

func main() {
	// Crear un servidor web en el puerto 8080
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}
