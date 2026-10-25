# Finca 🌿

Sistema de gestión de escritorio escrito en **Go**, con interfaz gráfica **GTK 4** y base de datos **MySQL**.

> Usa los widgets reutilizables del paquete `entradas` de **[joseamadeo/appgtk](https://github.com/joseamadeo/appgtk)**.

---

## ✨ Características

### Pantallas

| Menú | Pantalla | Descripción |
|---|---|---|
| **Egreso** | Cooperativa | Alta de egresos de cooperativa: número de factura, importes, fecha de vencimiento y cuentas contables |
| **Herramientas** | Usuarios | Alta, edición y baja de usuarios con contraseña |
| **Herramientas** | Cuentas | Plan de cuentas contables con tabla y alta de cuentas |

### Al iniciar

- **Login** con usuarios de la base de datos
- **Barra de estado** con el usuario logueado y el **reloj del servidor MySQL** (se actualiza cada minuto)
- **CSS embebido** (`style.css` con `go:embed`)

### Widgets reutilizables (de `appgtk/entradas`)

- **`Fecha`** — campo de fecha con calendario emergente, autoformato `DD/MM/AAAA`, validación Mín/Máx y estilos según validación.
- **`Numerico`** — campo numérico con formato de miles y decimales (`#.###,00`), normalización del portapapeles (formato argentino e internacional).
- **`NombrePropio`** — validación de nombres propios.
- **`CuentaContable`** — validación de cuenta contable.
- **`MiTabla[T]`** — tabla genérica para colecciones de cualquier tipo.
- **`CUIT`** — validación de CUIT.

---

## 🛠 Requisitos

**Para compilar:**

- [Go 1.25+](https://go.dev/dl/)
- GTK 4 y GLib de desarrollo con `pkg-config` (necesario por cgo/gotk4).
  En Windows la vía habitual es [MSYS2](https://www.msys2.org/):

  ```bash
  pacman -S mingw-w64-x86_64-gtk4 mingw-w64-x86_64-pkgconf
  ```

- MySQL/MariaDB corriendo en `localhost`

**Para ejecutar:** hace falta el runtime de GTK 4 accesible (DLLs en el `PATH`, por ejemplo las de MSYS2).

---

## ⚙️ Configuración de la base de datos

La conexión está definida en [`db/conexion.go`](db/conexion.go):

```go
usuario := "root"
clave   := "total"
servidor := "localhost"
base    := "finca"
```

La base debe llamarse `finca` y estar accesible por TCP en `localhost:3306`.

> ⚠️ Estos valores están pensados para un entorno de desarrollo local.

---

## ▶️ Compilar y ejecutar

```bash
# Desarrollo (con consola de logs)
go run .

# Con logs de los drivers (stderr)
set DEBUG_MODE=true

# Ejecutable de distribución (sin ventana de consola)
go build -o finca.exe -ldflags="-H=windowsgui" main.go
```

---

## 📁 Estructura del proyecto

```
finca/
├── main.go               # Punto de entrada: app GTK, ventana principal, CSS, reloj
├── style.css             # Hoja de estilos (embebida con go:embed)
├── go.mod / go.sum
├── db/                   # Capa de acceso a datos (MySQL)
│   ├── conexion.go       #   Singleton de conexión + transacciones + consultas
│   ├── cooperativa.go    #   ABM de egresos de cooperativa
│   ├── cuentas.go        #   ABM de cuentas contables
│   └── usuarios.go       #   ABM y autenticación de usuarios
├── modelos/
│   └── modelos.go        # Tipos compartidos
├── pantallas/            # Pantallas GTK 4
│   ├── login.go          #   Login
│   ├── menu.go           #   Menú principal
│   ├── barraEstado.go    #   Barra de estado con reloj
│   ├── cooperativa.go    #   Pantalla de cooperativa
│   ├── cuentas.go        #   Pantalla de cuentas
│   └── usuarios.go       #   Pantalla de usuarios
└── utils/
    └── dialogos.go       # Diálogos de mensaje reutilizables
```

---

## 📦 Dependencias

| Módulo | Uso |
|---|---|
| [`diamondburned/gotk4`](https://github.com/diamondburned/gotk4) | Bindings de GTK 4 |
| [`go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) | Driver de MySQL |
| [`joseamadeo/appgtk`](https://github.com/joseamadeo/appgtk) | Widgets de entrada reutilizables (`entradas/`) |

---

## 📄 Licencia

Proyecto personal. Consultar al autor respecto al uso y la distribución.
