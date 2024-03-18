package pdf

import (
	"bytes"
	"fmt"
	"calc/modelo"
	"github.com/jung-kurt/gofpdf"
)

func CrearReportePDF(historial []modelo.Historico) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16) // Cambia la fuente para el título
	pdf.SetY(10)                  // Posiciona el título a 10 mm del borde superior
	pdf.CellFormat(0, 10, "Registro de operaciones", "", 1, "C", false, 0, "")
	pdf.Ln(10) // Añade un poco de espacio antes de empezar con la tabla

	pdf.SetFont("Arial", "B", 12)

	pdf.SetFont("Arial", "B", 12)
	// Definir los anchos de las columnas
	anchos := []float64{42.5, 80.0, 40.0}
	// Aquí es donde debes calcular el margen izquierdo
	// Obtener el ancho y la altura de la página
	anchoPagina, _ := pdf.GetPageSize() // Asigna correctamente el ancho de la página

	// Calcular el ancho de la página restando los márgenes
	left, _, right, _ := pdf.GetMargins()    // Asigna correctamente los márgenes
	anchoPagina = anchoPagina - left - right // Realiza la resta en una línea separada

	// Continuar con el cálculo del margen izquierdo
	anchoTotalColumnas := 0.0
	for _, ancho := range anchos {
		anchoTotalColumnas += ancho
		fmt.Println(anchoTotalColumnas)
	}
	margenIzquierdo := (anchoPagina - anchoTotalColumnas) * 2
	fmt.Println(40)

	pdf.SetLeftMargin(margenIzquierdo - 10)

	// Definir los encabezados de la tabla
	encabezados := []string{"Fecha", "Operacion"}

	// Añadir los encabezados
	for i, encabezado := range encabezados {
		pdf.CellFormat(anchos[i], 7, encabezado, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Añadir los datos de la tabla
	var rows = [][]string{}
	for _, v := range historial {
		row := []string{v.Fecha, v.Operacion}
		rows = append(rows, row)
	}

	for _, linea := range rows {
		for i, dato := range linea {
			pdf.CellFormat(anchos[i], 15, dato, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	var Reporte bytes.Buffer

	err := pdf.Output(&Reporte)
	if err != nil {
		return nil, err
	}

	return &Reporte, nil
}
