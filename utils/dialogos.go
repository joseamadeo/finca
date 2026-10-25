package utils

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// MostrarError abre un diálogo modal de error vinculado a una ventana
func MostrarErrorFatal(parent *gtk.Window, titulo, mensaje string) {
	dialog := gtk.NewMessageDialog(
		parent,
		gtk.DialogModal,
		gtk.MessageError,
		gtk.ButtonsOK,
	)
	
	dialog.SetObjectProperty("text", titulo)
	dialog.SetObjectProperty("secondary-text", mensaje)
	
	// En GTK 4, ConnectResponse es la forma de capturar el clic en OK
	dialog.ConnectResponse(func(_ int) {
		dialog.Destroy()
		os.Exit(1) // Salir de la aplicación después de mostrar el error
	})
	
	dialog.Present()
}

// MostrarInformacion para mensajes de éxito o avisos
func MostrarInformacion(parent *gtk.Window, titulo, mensaje string) {
	dialog := gtk.NewMessageDialog(
		parent,
		gtk.DialogModal,
		gtk.MessageInfo,
		gtk.ButtonsOK,
	)
	
	dialog.SetObjectProperty("text", titulo)
	dialog.SetObjectProperty("secondary-text", mensaje)
	
	dialog.ConnectResponse(func(_ int) {
		dialog.Destroy()
	})
	
	dialog.Present()
}
