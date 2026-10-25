package pantallas

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Definimos una estructura para poder actualizar las etiquetas desde fuera si es necesario
type BarraEstado struct {
	Contenedor *gtk.Box
	LblAyuda   *gtk.Label
	LblReloj   *gtk.Label
}

func CrearBarraEstado(id int, nombre string) *BarraEstado {
	barra := gtk.NewBox(gtk.OrientationHorizontal, 10)
	barra.AddCSSClass("barraEstado")

	lblAyuda := gtk.NewLabel("Listo")
	lblAyuda.SetHExpand(true)
	lblAyuda.SetXAlign(0)

	lblUsuario := gtk.NewLabel(fmt.Sprintf("👤 %s (%d)", nombre, id))
	lblReloj := gtk.NewLabel("🕒 --:--")

	barra.Append(lblAyuda)
	barra.Append(lblUsuario)
	barra.Append(lblReloj)

	return &BarraEstado{
		Contenedor: barra,
		LblAyuda:   lblAyuda,
		LblReloj:   lblReloj,
	}
}