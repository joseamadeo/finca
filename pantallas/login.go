package pantallas

import (
	"finca/db"
	"finca/utils"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ShowLogin(app *gtk.Application, onSuccess func(idUsuario int, nombreUsuario string)) {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("Login")
	win.SetDefaultSize(300, 200)
	box := gtk.NewBox(gtk.OrientationVertical, 10)
	box.SetMarginTop(20)
	box.SetMarginBottom(20)
	box.SetMarginStart(20)
	box.SetMarginEnd(20)

	user := gtk.NewEntry()
	user.SetPlaceholderText("Usuario")
	user.SetText("amadeo") // Para pruebas, puedes eliminar esta línea en producción

	pass := gtk.NewPasswordEntry()
	pass.SetText("total") // Para pruebas, puedes eliminar esta línea en producción

	btn := gtk.NewButtonWithLabel("Login")
	btn.ConnectClicked(func() {
		dbUsuarios, err := db.NewDBUsuarios()
		if err != nil {
			// Si falló la conexión, mostramos el diálogo y salimos
      utils.MostrarErrorFatal(&win.Window, "Error de Sistema", "Fallo de conexión: " + err.Error())
			return
		}
		claveIngresadaMD5, err := dbUsuarios.LoginUsuario(user.Text(), pass.Text())
		// Validamos el usuario
		if err == nil  && dbUsuarios.Clave == claveIngresadaMD5 {
			idLogueado := dbUsuarios.IdUsuario
			nombreLogueado := dbUsuarios.Nombre
			win.Close() // Cerrar la ventana de login
			onSuccess(idLogueado, nombreLogueado)  // Llamar a la función de éxito
		} else {
			user.SetText("")
			pass.SetText("")
			// create a message dialog attached to window 'ww'
			dialog := gtk.NewMessageDialog(
				&win.Window,
				gtk.DialogModal,
				gtk.MessageError,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-error")
			dialog.SetObjectProperty("text", "Error de acceso")
			dialog.SetObjectProperty("secondary-text", "Usuario o contraseña incorrectos.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
		}
	})

	box.Append(user)
	box.Append(pass)
	box.Append(btn)
	win.SetChild(box)

	win.Present()
}