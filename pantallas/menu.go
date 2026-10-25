package pantallas

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ConfigurarMenu(app *gtk.Application, win *gtk.ApplicationWindow, contenedor *gtk.Box, usuarioUM int) {
	menubar := gio.NewMenu()
	mnuEgreso := gio.NewMenu()
	mnuEgreso.Append("Cooperativa", "win.cooperativa")


	menubar.AppendSubmenu("Egreso", mnuEgreso)

	accionElectrica := gio.NewSimpleAction("cooperativa", nil)
	accionElectrica.ConnectActivate(func(parameter *glib.Variant) {
		limpiarContenedor(contenedor)
		NewCooperativa(contenedor, &win.Window, usuarioUM)
	})
	win.AddAction(accionElectrica)

	mnuHerramientas := gio.NewMenu()

	mnuHerramientas.Append("Usuarios", "win.usuarios")
	mnuHerramientas.Append("Cuentas", "win.cuentas")
	menubar.AppendSubmenu("Herramientas", mnuHerramientas)

	accionUsuarios := gio.NewSimpleAction("usuarios", nil)
	accionUsuarios.ConnectActivate(func(_ *glib.Variant) {
		limpiarContenedor(contenedor)
		Usuarios(contenedor, &win.Window, usuarioUM) // Pasamos el contenedor y el ID del usuario para que sepa dónde mostrar su contenido
	})
	win.AddAction(accionUsuarios)

	accionCuentas := gio.NewSimpleAction("cuentas", nil)
	accionCuentas.ConnectActivate(func(_ *glib.Variant) {
		limpiarContenedor(contenedor)
		NewCuentas(contenedor, &win.Window, usuarioUM) // Pasamos el contenedor y el ID del usuario para que sepa dónde mostrar su contenido
	})
	win.AddAction(accionCuentas)


	app.SetMenubar(menubar)
	
	// Mostrar la barra de menú en la ventana
	win.SetShowMenubar(true)
}

func limpiarContenedor(contenedor *gtk.Box) {
	for child := contenedor.FirstChild(); child != nil; child = contenedor.FirstChild() {
		contenedor.Remove(child)
	}
}