package pantallas

import (
	"finca/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
	"github.com/joseamadeo/appgtk/estructuras"
)

var dbCuentas, _ = db.NewDBCuentas()
var cuentaSeleccionada = &db.DBCuentas{}

func NewCuentas(contenedor *gtk.Box, win *gtk.Window, usuarioUM int) {
	lista, err := dbCuentas.ListaCuentas()
	if err != nil {
		mostrarErrorCuentas(win, "Error de Conexión",
			fmt.Sprintf("No se pudo obtener la lista de cuentas, el error es: %v", err))
		return
	}

	// Definición de columnas para la tabla genérica de appgtk
	columnas := []estructuras.ColumnaConfig[*db.DBCuentas]{
		{
			Titulo:  "Número",
			Extraer: func(c *db.DBCuentas) string { return c.Numero },
			Comparar: func(a, b *db.DBCuentas) int {
				return strings.Compare(a.Numero, b.Numero)
			},
		},
		{
			Titulo:  "Nombre",
			Extraer: func(c *db.DBCuentas) string { return c.Nombre },
			Comparar: func(a, b *db.DBCuentas) int {
				return strings.Compare(a.Nombre, b.Nombre)
			},
		},
		{
			Titulo: "I/E",
			Extraer: func(c *db.DBCuentas) string {
				if c.IngresoEgreso {
					return "Ingreso"
				}
				return "Egreso"
			},
		},
		{
			Titulo: "Imp.",
			Extraer: func(c *db.DBCuentas) string {
				if c.Imputable {
					return "Sí"
				}
				return "No"
			},
		},
	}

	tabla := entradas.NewTabla(lista, columnas)
	tabla.Widget.SetVExpand(true)
	tabla.Widget.SetMinContentHeight(300) // Altura mínima inicial

	// FORMULARIO
	formGrid := gtk.NewGrid()
	formGrid.SetColumnSpacing(15)
	formGrid.SetRowSpacing(10)
	formGrid.SetMarginTop(20)
	formGrid.SetMarginStart(10)

	// --- Definir los Widgets del Formulario ---
	txtNumero := entradas.NewCuentaContable()            // Para el Número
	txtNombre := entradas.NewNombrePropio("Nombre de la cuenta", false, 60)
	swIngEgr := gtk.NewSwitch() // Para Ingreso/Egreso
	swIngEgr.SetHAlign(gtk.AlignStart)
	swIngEgr.SetVAlign(gtk.AlignCenter)
	swImputable := gtk.NewSwitch() // Para Imputable
	swImputable.SetHAlign(gtk.AlignStart)
	swImputable.SetVAlign(gtk.AlignCenter)

	// --- Ubicar en la Grilla (Columna, Fila, Ancho, Alto) ---
	formGrid.Attach(gtk.NewLabel("Número:"), 0, 0, 1, 1)
	formGrid.Attach(txtNumero, 1, 0, 1, 1)

	formGrid.Attach(gtk.NewLabel("Nombre:"), 0, 1, 1, 1)
	formGrid.Attach(txtNombre, 1, 1, 1, 1)

	formGrid.Attach(gtk.NewLabel("Ingreso/Egreso:"), 0, 2, 1, 1)
	formGrid.Attach(swIngEgr, 1, 2, 1, 1)

	formGrid.Attach(gtk.NewLabel("Imputable:"), 0, 3, 1, 1)
	formGrid.Attach(swImputable, 1, 3, 1, 1)

	// Botones de acción
	boxAcciones := gtk.NewBox(gtk.OrientationHorizontal, 10)
	btnGuardar := gtk.NewButtonWithLabel("Guardar Nuevo")
	btnLimpiar := gtk.NewButtonWithLabel("Nuevo / Limpiar")

	boxAcciones.Append(btnGuardar)
	boxAcciones.Append(btnLimpiar)

	formGrid.Attach(boxAcciones, 1, 4, 1, 1)

	// al cambiar la selección
	tabla.AlSeleccionar(func(cuenta *db.DBCuentas) {
		cuentaSeleccionada = cuenta
		// llena los datos en el formulario
		txtNumero.SetCuenta(cuenta.Numero)
		txtNombre.SetTexto(cuenta.Nombre)
		swIngEgr.SetActive(cuenta.IngresoEgreso)
		swImputable.SetActive(cuenta.Imputable)
		btnGuardar.SetLabel("Guardar Cuenta Seleccionada")
	})

	// click de los botones
	btnLimpiar.ConnectClicked(func() {
		tabla.Seleccion.SetSelected(gtk.InvalidListPosition)
		cuentaSeleccionada = &db.DBCuentas{}
		txtNumero.SetCuenta("")
		txtNombre.SetTexto("")
		swIngEgr.SetActive(false)
		swImputable.SetActive(false)
		btnGuardar.SetLabel("Guardar Nuevo")
	})

	btnGuardar.ConnectClicked(func() {
		// Verifica si los datos están correctos
		dialog := gtk.NewMessageDialog(
			win,
			gtk.DialogModal,
			gtk.MessageInfo,
			gtk.ButtonsOKCancel,
		)
		msj := "Se encontrado los siguientes errores:"
		if !txtNumero.EsValido() {
			msj += "\n - El Cuenta Contable esta incorrecto"
		}
		if len(txtNombre.GetTexto()) < 4 {
			msj += "\n - Debe ingresar el nombre de la cuenta"
		}
		if msj != "Se encontrado los siguientes errores:" {
			dialog.SetTitle("Faltante de datos.")
			msg := gtk.NewLabel(msj)
			dialog.ContentArea().Append(msg)
			dialog.SetIconName("dialog-information")
			dialog.ConnectResponse(func(r int) {
				dialog.Destroy()
				return
			})
			dialog.Present()
			return
		}
		dbCuenta, err := db.NewDBCuentas()
		if err != nil {
			dialog.SetTitle("Error de conexión.")
			msg := gtk.NewLabel("No se puede tener acceso a la base de datos")
			dialog.ContentArea().Append(msg)
			dialog.SetIconName("dialog-information")
			dialog.ConnectResponse(func(r int) {
				dialog.Destroy()
				return
			})
			dialog.Present()
			return
		}
		dbCuenta.Numero = txtNumero.GetCuenta()
		dbCuenta.Nombre = txtNombre.GetTexto()
		dbCuenta.Imputable = swImputable.Active()
		dbCuenta.IngresoEgreso = swIngEgr.Active()
		dbCuenta.UsuarioUM = usuarioUM
		if btnGuardar.Label() == "Guardar Nuevo" {
			err = dbCuenta.AgregarCuenta()
			if err != nil {
				dialog.SetTitle("Error al guardar.")
				msg := gtk.NewLabel(fmt.Sprintf("No se pudo guardar la cuenta Nueva, el error es: %v", err))
				dialog.ContentArea().Append(msg)
				dialog.SetIconName("dialog-error")
				dialog.ConnectResponse(func(r int) {
					dialog.Destroy()
					return
				})
				dialog.Present()
				return
			}
		} else {
			err = dbCuenta.ModificarCuenta(cuentaSeleccionada.Numero)
			if err != nil {
				dialog.SetTitle("Error al guardar.")
				msg := gtk.NewLabel(fmt.Sprintf("No se pudo guardar la cuenta Actualizada, el error es: %v", err))
				dialog.ContentArea().Append(msg)
				dialog.SetIconName("dialog-error")
				dialog.ConnectResponse(func(r int) {
					dialog.Destroy()
					return
				})
				dialog.Present()
				return
			}
		}
		// Recargar la tabla y limpiar el formulario
		recargarCuentas(win, tabla)
		tabla.Seleccion.SetSelected(gtk.InvalidListPosition)
		cuentaSeleccionada = &db.DBCuentas{}

		txtNumero.SetCuenta("")
		txtNombre.SetTexto("")
		swIngEgr.SetActive(false)
		swImputable.SetActive(false)
		btnGuardar.SetLabel("Guardar Nuevo")
	})

	// Añadir el formulario al contenedor principal
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)
	vbox.Append(tabla.Widget)
	vbox.Append(formGrid)
	contenedor.Append(vbox)
}

// recargarCuentas vuelve a leer la lista de cuentas desde la base de datos y
// refresca el contenido de la tabla genérica manteniendo los índices del store
// sincronizados con el slice de datos.
func recargarCuentas(win *gtk.Window, tabla *entradas.MiTabla[*db.DBCuentas]) {
	lista, err := dbCuentas.ListaCuentas()
	if err != nil {
		mostrarErrorCuentas(win, "Error de Conexión",
			fmt.Sprintf("No se pudo obtener la lista de cuentas, el error es: %v", err))
		return
	}
	tabla.Datos = lista
	tabla.Store.Splice(0, tabla.Store.NItems(), nil)
	for i := range lista {
		tabla.Store.Append(strconv.Itoa(i))
	}
}

// mostrarErrorCuentas muestra un diálogo modal de error estándar
func mostrarErrorCuentas(win *gtk.Window, titulo, mensaje string) {
	dialog := gtk.NewMessageDialog(
		win,
		gtk.DialogModal,
		gtk.MessageError,
		gtk.ButtonsOK,
	)
	dialog.SetTitle(titulo)
	dialog.ContentArea().Append(gtk.NewLabel(mensaje))
	dialog.SetIconName("dialog-error")
	dialog.ConnectResponse(func(r int) {
		dialog.Destroy()
	})
	dialog.Present()
}
