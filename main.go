// main.go
package main

import (
	db "calc/bd"
	"calc/pdf"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
)

var templates = template.Must(template.New("T").ParseGlob("static/*.html"))

func renderTemplate(rw http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(rw, name, data)
	if err != nil {
		panic(err)
	}
}

// Función para calcular el resultado de una operación
func calcular(operando1, operando2 float64, operador string) float64 {
	var resultado float64
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
		resultado = math.Sqrt(operando1)
		resultado = math.Round(resultado*100000000) / 100000000

	case "!":
		var factorial int64
		var numero int64
		numero = int64(operando1)
		factorial = 1
		for i := int64(1); i <= numero; i++ {
			factorial *= i
		}
		resultado = float64(factorial)
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

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if r.Method == "POST" {
		ValorFecha := r.FormValue("ordenar-fecha")
		fmt.Println(ValorFecha)

		if ValorFecha != "todos"{
		Mapa := make(map[string]interface{})
		Mapa["Resultados"] = db.RecogerHistorialPorFecha(ValorFecha)
		Mapa["Fechas"] = db.RecogerFechas()
		renderTemplate(w, "index.html", Mapa)
		return
		}
	}
	
	operando1 := r.URL.Query().Get("operando1")
	operando2 := r.URL.Query().Get("operando2")
	operador := r.URL.Query().Get("operador")
	timestamp := r.URL.Query().Get("timestamp")
	result := r.URL.Query().Get("resultado")

	// Si hay valores, convertirlos a números y calcular el resultado
	if result != "" && timestamp != "" {
		err := db.LlenarHistorial(result, timestamp)
		if err != nil {
			panic(err)
		}
	}
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
		Fechas := db.RecogerFechas()
		Resultados := db.RecogerHistorialCompleto()
		Mapa := make(map[string]interface{})
		Mapa["Resultados"] = Resultados
		Mapa["Fechas"] = Fechas
		renderTemplate(w, "index.html", Mapa)
		return
	}
}

func Descargar(w http.ResponseWriter, r *http.Request) {

	//Recoger registros de la base de datos

	Registro := db.RecogerHistorialCompleto()
	// Generar el PDF y guardar el contenido en un buffer
	pdfContent, err := pdf.CrearReportePDF(Registro)
	if err != nil {
		http.Error(w, "Error al generar el PDF", http.StatusInternalServerError)
		return
	}

	loc, _ := time.LoadLocation("America/New_York")
	fechaNY := time.Now().In(loc)

	// Formatear la fecha y hora al formato deseado
	fechaFormato := fechaNY.Format("02-01-2006 15:04:05")

	// Crear el nombre del archivo con la fecha y hora actual
	nombreArchivo := fmt.Sprintf("Registro_Operaciones_%s.pdf", fechaFormato)

	// Configurar encabezados para la descarga
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nombreArchivo))
	w.Header().Set("Content-Type", "application/pdf")

	// Escribir el contenido del PDF en la respuesta
	_, err = w.Write(pdfContent.Bytes())
	if err != nil {
		http.Error(w, "Error al escribir el PDF en la respuesta", http.StatusInternalServerError)
		return
	}

}

func main() {
	// Crear un servidor web en el puerto 8080
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", calculadora)
	http.HandleFunc("/descargar", Descargar)

	http.ListenAndServe(":8080", nil)
}
