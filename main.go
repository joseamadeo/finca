package main

import (
	_ "embed"
	"finca/db"
	"finca/pantallas"
	"fmt"
	"os"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var cssData string

var usuarioActualID int
var usuarioActualNombre string
var lblAyuda *gtk.Label
var lblReloj *gtk.Label

func main() {
	// Solo en Windows para limpiar la consola de logs de drivers
  if os.Getenv("DEBUG_MODE") != "true" {
    os.Stderr, _ = os.Open(os.DevNull)
  }
	// Silencia los logs de GDK (Vulkan, Gl, etc.)
	os.Setenv("GDK_DEBUG", "")
	os.Setenv("GLIB_DEBUG", "")
	// Opcional: si querés ocultar todos los mensajes de GLib/GTK
	os.Setenv("G_MESSAGES_DEBUG", "none")
	// Forzar renderizado por software (desactiva Vulkan/GL)
  os.Setenv("GSK_RENDERER", "cairo")


	// 1. Crear la aplicación
	app := gtk.NewApplication("com.jas.finca", gio.ApplicationFlagsNone)

	// 2. Definir qué sucede cuando se activa la app
	app.ConnectActivate(func() {
		activate(app)
	})

	// 3. Ejecutar la aplicación
	if code := app.Run(os.Args); code > 0 {
		fmt.Printf("La aplicación terminó con código de salida: %d\n", code)
		os.Exit(code)
	}
}

func loadCss() func() { 
	return func() {
		cssProvider := gtk.NewCSSProvider()
		cssProvider.LoadFromString(cssData) 
		gtk.StyleContextAddProviderForDisplay(
			gdk.DisplayGetDefault(),
			cssProvider,
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
		)
	}
}

func activate(app *gtk.Application) {
	lblAyuda = gtk.NewLabel("Listo")
	lblReloj = gtk.NewLabel("🕒 --:--")

	display := gdk.DisplayGetDefault()
	if display != nil {
		// Configuramos el nombre del tema
		settings := gtk.SettingsGetForDisplay(display)
		settings.SetObjectProperty("gtk-icon-theme-name", "Adwaita")

		// Forzamos a GTK a buscar en tu carpeta local 'share/icons'
		iconTheme := gtk.IconThemeGetForDisplay(display)
		iconTheme.AddSearchPath("share\\icons")
	}else{
		fmt.Println("**** No se pudo obtener el display para configurar los iconos.")
	}

	// Mostrar la ventana de login
	pantallas.ShowLogin(app ,func(id int, nombre string) {
		usuarioActualID = id
		usuarioActualNombre = nombre
		loadCss()() // Cargar el CSS después del login exitoso
		// 1 - Crear la ventana principal
		win := gtk.NewApplicationWindow(app)

		win.SetTitle("Sistema de Gestión Finca")
		win.SetDefaultSize(800, 700)
		win.SetShowMenubar(true)

		vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		win.SetChild(vbox)

		// Aquí iría el contenido principal de la ventana (por ejemplo, un panel o grilla)
		contenidoPrincipal := gtk.NewBox(gtk.OrientationVertical, 0)
		contenidoPrincipal.SetVExpand(true)
		vbox.Append(contenidoPrincipal)
		barraObj := pantallas.CrearBarraEstado(usuarioActualID, usuarioActualNombre)
		vbox.Append(barraObj.Contenedor)
		pantallas.ConfigurarMenu(app, win, contenidoPrincipal, usuarioActualID)

		// Actualizar el reloj cada 1 minuto
		go func() {
			con, err := db.GetInstance()
			if err != nil {
				fmt.Println("Error al obtener conexión para reloj:", err)
				return
			}
			for {
				hora, err := con.GetFechaServidor()
					if err == nil {
					glib.IdleAdd(func() {
						barraObj.LblReloj.SetText("🕒 " + hora)
					})
				}
				time.Sleep(1 * time.Minute)
			}
		}()

		// 6 - Presenta la ventana
		win.Present()
	})
}
