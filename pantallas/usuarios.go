package pantallas

import (
	"finca/db"
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
)

var dbusuario, _ = db.NewDBUsuarios()

func Usuarios(contenedor *gtk.Box, win *gtk.Window, usuarioUM int) {
	// 1. Crear el Paned (divisor horizontal)
	panelDivisor := gtk.NewPaned(gtk.OrientationHorizontal)
	panelDivisor.SetMarginTop(10)
	panelDivisor.SetMarginBottom(10)
	panelDivisor.SetMarginStart(10)
	panelDivisor.SetMarginEnd(10)
	panelDivisor.SetWideHandle(true)

	// 1. Contenedor vertical para Buscador + Lista
	vboxIzquierda := gtk.NewBox(gtk.OrientationVertical, 5)

	// 2. El Buscador (arriba)
	txtBuscar := gtk.NewEntry()
	txtBuscar.SetPlaceholderText("Buscar...")
	txtBuscar.SetMarginBottom(5)
	vboxIzquierda.Append(txtBuscar)

	// 3. El Listado con Scroll
	lstUsuarios := gtk.NewListBox()
	scrollList := gtk.NewScrolledWindow()
	scrollList.SetChild(lstUsuarios)
	scrollList.SetVExpand(true) // La lista ocupa el resto del alto disponible
	vboxIzquierda.Append(scrollList)

	txtBuscar.Connect("changed", func() {
		actualizarListaUsuarios(lstUsuarios, txtBuscar.Text())
	})

	// --- LADO DERECHO: Formulario ---
	formGrid := gtk.NewGrid()
	formGrid.SetColumnSpacing(10)
	formGrid.SetRowSpacing(8)
	formGrid.SetMarginStart(10)

	// Campos Simples
	lblId := gtk.NewLabel("id")
	lblIdVal := gtk.NewLabel("") // LblId
	formGrid.Attach(lblId, 0, 0, 1, 1)
	formGrid.Attach(lblIdVal, 1, 0, 1, 1)

	txtUsuario := gtk.NewEntry()
	formGrid.Attach(gtk.NewLabel("Usuario"), 0, 1, 1, 1)
	formGrid.Attach(txtUsuario, 1, 1, 1, 1)

	// El Nombre es de una persona, asi que usa el widget NombrePropio
	// (filtra caracteres, formato nombre propio y validación visual)
	txtNombre := entradas.NewNombrePropio("Nombre y Apellido", false, 60)
	formGrid.Attach(gtk.NewLabel("Nombre"), 0, 2, 1, 1)
	formGrid.Attach(txtNombre, 1, 2, 1, 1)

	// Checkbox Activo
	chkActivo := gtk.NewCheckButtonWithLabel("Activo")
	formGrid.Attach(chkActivo, 1, 3, 1, 1)

	// Passwords
	txtPass1 := gtk.NewPasswordEntry()
	formGrid.Attach(gtk.NewLabel("Clave"), 0, 4, 1, 1)
	formGrid.Attach(txtPass1, 1, 4, 1, 1)

	txtPass2 := gtk.NewPasswordEntry()
	formGrid.Attach(gtk.NewLabel("Repetir Clave"), 0, 5, 1, 1)
	formGrid.Attach(txtPass2, 1, 5, 1, 1)
	
	lstUsuarios.Connect("row-activated", func(row *gtk.ListBoxRow) {
		// Se obtiente el Id del usuario seleccionado
		idUsuario, _ := strconv.Atoi(lstUsuarios.SelectedRow().Name())
		if idUsuario == 0 {
			// Nuevo usuario: Limpiar campos
			lblIdVal.SetText("0")
			txtUsuario.SetText("")
			txtNombre.SetTexto("")
			chkActivo.SetActive(true)
			txtPass1.SetText("")
			txtPass2.SetText("")
			
			return
		}
		// Cargar datos del usuario seleccionado
		err := dbusuario.BuscarUsuario(idUsuario)
		if err != nil {
			// create a message dialog attached to window 'ww'
			dialog := gtk.NewMessageDialog(
				win,
				gtk.DialogModal,
				gtk.MessageError,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-error")
			dialog.SetObjectProperty("text", "Error leyendo usuario")
			dialog.SetObjectProperty("secondary-text", "No se encontró el usuario seleccionado.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
			return
		}
		lblIdVal.SetText(strconv.Itoa(dbusuario.IdUsuario))
		txtUsuario.SetText(dbusuario.Usuario)
		txtNombre.SetTexto(dbusuario.Nombre)
		chkActivo.SetActive(dbusuario.Activo)
		txtPass1.SetText("") // Por seguridad, no mostrar la clave
		txtPass2.SetText("")
	})

	// --- BOTONES INFERIORES ---
	boxButones := gtk.NewBox(gtk.OrientationHorizontal, 10)
	boxButones.SetHAlign(gtk.AlignCenter)
	
	btnAceptar := gtk.NewButtonWithLabel("Aceptar")
	btnAceptar.ConnectClicked(func() {
		// Validar campos
		if txtUsuario.Text() == "" || txtNombre.GetTexto() == "" {
			dialog := gtk.NewMessageDialog(
				win,
				gtk.DialogModal,
				gtk.MessageWarning,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-warning")
			dialog.SetObjectProperty("text", "Campos incompletos")
			dialog.SetObjectProperty("secondary-text", "Por favor, complete los campos Usuario y Nombre.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
			return
		}
		if txtPass1.Text() != txtPass2.Text() {
			dialog := gtk.NewMessageDialog(
				win,
				gtk.DialogModal,
				gtk.MessageWarning,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-warning")
			dialog.SetObjectProperty("text", "Claves no coinciden")
			dialog.SetObjectProperty("secondary-text", "Las claves ingresadas no coinciden.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
			return
		}

		// Guardar en DB
		usuario, err := db.NewDBUsuarios()
		if err != nil {
			dialog := gtk.NewMessageDialog(
				win,
				gtk.DialogModal,
				gtk.MessageError,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-error")
			dialog.SetObjectProperty("text", "Error creando usuario")
			dialog.SetObjectProperty("secondary-text", "Ocurrió un error al crear el usuario. Intente nuevamente.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
			return
		}

		usuario.IdUsuario =	dbusuario.IdUsuario
		usuario.Usuario = txtUsuario.Text()
		usuario.Nombre = txtNombre.GetTexto()
		usuario.Activo = chkActivo.Active() // Usar el valor del checkbox
		usuario.Clave = txtPass1.Text()
		usuario.UsuarioUM = usuarioUM // Aquí podrías usar el ID del usuario logueado, si tienes esa info
		usuario.FechaExpiracion = time.Now().AddDate(0, 1, 0) // Por ejemplo, expira en 3 meses
		if usuario.IdUsuario == 0 {
			err = usuario.AgregarUsuario()
		} else {
			err = usuario.ModificarUsuario()
		}	
		if err != nil {
			dialog := gtk.NewMessageDialog(
				win,
				gtk.DialogModal,
				gtk.MessageError,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-error")
			dialog.SetObjectProperty("text", "Error guardando usuario")
			dialog.SetObjectProperty("secondary-text", "Ocurrió un error al guardar el usuario. Intente nuevamente.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
			return
		}

		// Refrescar lista y limpiar formulario
		actualizarListaUsuarios(lstUsuarios, txtBuscar.Text())
		lblIdVal.SetText("0")
		txtUsuario.SetText("")
		txtNombre.SetTexto("")
		txtPass1.SetText("")
		txtPass2.SetText("")

	})

	btnCerrar := gtk.NewButtonWithLabel("Cerrar")
	btnCerrar.ConnectClicked(func() {
		// Limpiar el contenido del formulario
		contenedor.Remove(panelDivisor)
	})
	
	boxButones.Append(btnAceptar)
	boxButones.Append(btnCerrar)

	formGrid.Attach(boxButones, 0, 6, 2, 1)

	panelDivisor.SetStartChild(vboxIzquierda)
	panelDivisor.SetEndChild(formGrid)

	// 2. Configurar la proporción al cambiar el tamaño de la ventana
	win.Connect("resize", func(width int, height int) {
    posicion := int(float64(width) * 0.40)
    panelDivisor.SetPosition(posicion)
	})

	contenedor.Append(panelDivisor)
}

// Función para refrescar los datos en el ListBox
func actualizarListaUsuarios(lst *gtk.ListBox, filtro string) {
	// 1. Limpiar la lista actual
	for {
		hijos := lst.FirstChild()
		if hijos == nil {
			break
		}
		lst.Remove(hijos)
	}

	// 2. Obtener datos de tu DB (puedes pasar el filtro a tu función si la editas, 
	// o filtrar el slice resultante aquí)
	usuarios, err := dbusuario.ListaUsuario(filtro, true)
	if err != nil {
		println("Error al cargar usuarios:", err.Error())
		return
	}

	// 3. Crear los widgets para cada fila
	for _, u := range usuarios {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(u.Nombre) // Suponiendo que ItemData tiene campo Nombre
		label.SetHAlign(gtk.AlignStart)
		label.SetMarginStart(10)
		label.SetMarginTop(5)
		label.SetMarginBottom(5)
		row.SetName(fmt.Sprint(u.ID))		
		row.SetChild(label)
	
		lst.Append(row)
	}
}
