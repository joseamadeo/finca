package pantallas

import (
	"finca/db"
	"fmt"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
)

var dbCooperativa, _ = db.NewDBCooperativo()

// globalCuentas y crearFabricaTextual se mantienen para la tabla de
// Cooperativa, que todavía no está conectada a datos propios (el llenado está
// pendiente, ver comentario de llenarTablaCooperativa más abajo).
var globalCuentas []*db.DBCuentas

func crearFabricaTextual(campo string) *gtk.ListItemFactory {
	factory := gtk.NewSignalListItemFactory()

	factory.ConnectSetup(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		label := gtk.NewLabel("")
		label.SetXAlign(0) // Alineación a la izquierda
		cell.SetChild(label)
	})

	factory.ConnectBind(func(obj *glib.Object) {
		cell := obj.Cast().(*gtk.ColumnViewCell)
		item := cell.Item()
		if item == nil {
			return
		}
		strObj := item.Cast().(*gtk.StringObject)
		var idx int
		fmt.Sscanf(strObj.String(), "%d", &idx)

		if idx < len(globalCuentas) {
			cuenta := globalCuentas[idx]
			if child := cell.Child(); child != nil {
				if label, ok := child.(*gtk.Label); ok {
					switch campo {
					case "Numero":
						label.SetText(cuenta.Numero)
					case "Nombre":
						label.SetText(cuenta.Nombre)
					case "Tipo":
						if cuenta.IngresoEgreso {
							label.SetText("Ingreso")
						} else {
							label.SetText("Egreso")
						}
					case "Imputable":
						if cuenta.Imputable {
							label.SetText("Sí")
						} else {
							label.SetText("No")
						}
					}
				}
			}
		}
	})

	return &factory.ListItemFactory
}

func NewCooperativa(contenedor *gtk.Box, win *gtk.Window, usuarioUM int) {
	store := gio.NewListStore(glib.TypeObject)
	sel := gtk.NewSingleSelection(store)
	sel.SetAutoselect(false)
	sel.SetCanUnselect(true)
	sel.SetSelected(gtk.InvalidListPosition)

	tabla := gtk.NewColumnView(sel)
	tabla.SetShowColumnSeparators(true)
	tabla.SetShowRowSeparators(true)

	// 1. Definición corregida del struct
	type ColumnaConfig struct {
		Titulo string
		Campo  string
		Expand bool
	}
	// 2. Inicialización consistente
	columnas := []ColumnaConfig{
		{"Cuenta", "Cuenta", false},
		{"NCuenta", "Nombre Cuenta", true},
		{"Número", "Numero", false},
		{"Período", "Periodo", false},
		{"Alta", "Alta", false},
		{"Baja", "Baja", false},
		{"Importe", "Importe", false},
		{"F.Vencimiento", "Vencimiento", false},
		{"F.Pago", "Pago", false},
	}

	for _, c := range columnas {
		col := gtk.NewColumnViewColumn(c.Titulo, crearFabricaTextual(c.Campo))
		col.SetExpand(c.Expand)
		tabla.AppendColumn(col)
	}
	//llenarTablaCooperativa(win, store, tabla)

	// 3. Recomendación: Usar un ScrolledWindow para la tabla
	scroll := gtk.NewScrolledWindow()
	scroll.SetChild(tabla)
	scroll.SetVExpand(true)

	notebook :=gtk.NewNotebook()

	gridFactura := gtk.NewGrid()
	gridFactura.SetMarginTop(10)
	gridFactura.SetMarginStart(20)
	gridFactura.SetRowSpacing(10)
	gridFactura.SetColumnSpacing(10)

	gridFactura.Attach(gtk.NewLabel("Nombre Cuenta:"), 0, 0, 1, 1)
	
	datosCuentas, _ := dbCooperativa.ObtenerCuentaContable()
	var opciones []string
	for _, c := range datosCuentas {
		opciones = append(opciones, fmt.Sprintf("%s - %s", c.Numero, c.Nombre))
	}
	
	dropdown := gtk.NewDropDownFromStrings(opciones)
	gridFactura.Attach(dropdown, 1, 0, 1, 1)

	txtNumeroFactura := entradas.NewNumerico("#.###")
	gridFactura.Attach(gtk.NewLabel("Número Factura:"), 0, 1, 1, 1)
	gridFactura.Attach(txtNumeroFactura.Object(), 1, 1, 1, 1)

	txtPeriodo := gtk.NewEntry()
	gridFactura.Attach(gtk.NewLabel("Período:"), 2, 1, 1, 1)
	gridFactura.Attach(txtPeriodo, 3, 1, 1, 1)

	txtBaja := entradas.NewNumerico("#.###")
	gridFactura.Attach(gtk.NewLabel("Baja Kw:"), 0, 3, 1, 1)
	gridFactura.Attach(txtBaja.Object(), 1, 3, 1, 1)

	txtAlta := entradas.NewNumerico("#.###")
	gridFactura.Attach(gtk.NewLabel("Alta Kw:"), 2, 3, 1, 1)
	gridFactura.Attach(txtAlta.Object(), 3, 3, 1, 1)

	txtImporte := entradas.NewNumerico("#.###,00")
	gridFactura.Attach(gtk.NewLabel("Importe:"), 0, 4, 1, 1)
	gridFactura.Attach(txtImporte.Object(), 1, 4, 1, 1)

	fchVenc := entradas.NewFecha()
	gridFactura.Attach(gtk.NewLabel("F. Vencimiento:"), 2, 4, 1, 1)
	gridFactura.Attach(fchVenc, 3, 4, 1, 1)

	cmdAgregar := gtk.NewButtonWithLabel("Agregar")
	cmdAgregar.SetHAlign(gtk.AlignCenter)
	cmdAgregar.SetHExpand(false)
	gridFactura.Attach(cmdAgregar, 0, 5, 4, 1)

	notebook.AppendPage(gridFactura, gtk.NewLabel("Factura"))

	// --- PESTAÑA 2: Espacio para otra ventana/vista ---
	label2 := gtk.NewLabel("Contenido de la segunda pestaña")
	notebook.AppendPage(label2, gtk.NewLabel("Pagos"))


	// Añadir el formulario al contenedor principal
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)
	vbox.Append(scroll)
	vbox.Append(notebook)
	contenedor.Append(vbox)

}
