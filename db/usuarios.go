package db

import (
	"finca/modelos"
	"fmt"
	"time"
)

// DBUsuarios estructura que representa al usuario
type DBUsuarios struct {
	con             *Conexion
	IdUsuario       int
	Usuario         string
	Nombre          string
	Activo          bool
	Clave           string
	FechaExpiracion time.Time
	UsuarioUM       int
}

// NewDBUsuarios constructor
func NewDBUsuarios() (*DBUsuarios, error) {
	con, err := GetInstance()
	if err != nil {
		return nil, err
	}
	return &DBUsuarios{
		con: con,
	}, nil
}

// ListaUsuario devuelve un slice de ItemData (reemplaza al ArrayList)
func (db *DBUsuarios) ListaUsuario(filtro string, nuevo bool) ([]modelos.ItemData, error) {
	var lista []modelos.ItemData
	if nuevo {
		lista = append(lista, modelos.ItemData{ID: 0, Nombre: "  Nuevo"})
	}

	query := "SELECT id, nombre FROM usuarios WHERE nombre LIKE '%" + filtro + "%';"
	dt, err := db.con.MuestraDatosDT(query)
	if err != nil {
		return nil, err
	}

	for _, row := range dt {
		id, _ := row["id"].(int64) // El driver de Go suele devolver int64 para IDs
		nombre, _ := row["nombre"].(string)
		lista = append(lista, modelos.ItemData{ID: int(id), Nombre: nombre})
	}
	return lista, nil
}

// BuscarUsuario busca y carga los datos en la instancia actual
func (db *DBUsuarios) BuscarUsuario(idU int) error {
	sq := fmt.Sprintf("SELECT usuario, nombre, activo, clave, fechaExpiracion FROM usuarios WHERE id = %d", idU)
	dt, err := db.con.MuestraDatosDT(sq)
	if err != nil || len(dt) == 0 {
		return err
	}

	row := dt[0]
	db.IdUsuario = idU
	db.Usuario = row["usuario"].(string)
	db.Nombre = row["nombre"].(string)
	
	// Manejo del booleano (MySQL guarda 0 o 1)
	activo, _ := row["activo"].(int64)
	db.Activo = activo != 0
	
	db.Clave = row["clave"].(string)
	// Recuperar fecha de expiración de forma segura
	if fechaVal, ok := row["fechaExpiracion"]; ok {
    switch v := fechaVal.(type) {
    case time.Time:
      db.FechaExpiracion = v
    case string:
      // Si por alguna razón viene como string, lo parseamos
      db.FechaExpiracion, _ = time.Parse("2006-01-02", v)
    }
	}
	
	return nil
}

// AgregarUsuario inserta un nuevo registro
func (db *DBUsuarios) AgregarUsuario() error {
	activoInt := 0
	if db.Activo {
		activoInt = 1
	}
	
	fechaStr := db.FechaExpiracion.Format("2006-01-02")
	sq := fmt.Sprintf("INSERT INTO usuarios(usuario, nombre, activo, clave, fechaExpiracion, usuarioUM) VALUES ('%s', '%s', %d, MD5('%s'), '%s', %d)",
		db.Usuario, db.Nombre, activoInt, db.Clave, fechaStr, db.UsuarioUM)
	return db.con.Transaccion(sq)
}

// ModificarUsuario actualiza un registro existente
func (db *DBUsuarios) ModificarUsuario() error {
	activoInt := 0
	if db.Activo {
		activoInt = 1
	}
	fechaStr := db.FechaExpiracion.Format("2006-01-02")
	sq := fmt.Sprintf("UPDATE usuarios SET usuario='%s', nombre='%s', activo=%d, clave=MD5('%s'), fechaExpiracion='%s', usuarioUM=%d WHERE id=%d",
		db.Usuario, db.Nombre, activoInt, db.Clave, fechaStr, db.UsuarioUM, db.IdUsuario)

	return db.con.Transaccion(sq)
}

// ActualizaClave actualiza solo la contraseña y expiración
func (db *DBUsuarios) ActualizaClave(usuario string, clave string, nuevaFecha time.Time) error {
	fechaStr := nuevaFecha.Format("2006-01-02")
	sq := fmt.Sprintf("UPDATE usuarios SET clave=MD5('%s'), fechaExpiracion='%s' WHERE usuario='%s'",
		clave, fechaStr, usuario)
	
	return db.con.Transaccion(sq)
}

// LoginUsuario verifica credenciales y devuelve el MD5 de la clave ingresada
func (db *DBUsuarios) LoginUsuario(usuario string, clave string) (string, error) {
	var claveEncriptada string
	sq := fmt.Sprintf("SELECT id, usuario, nombre, activo, clave, fechaExpiracion, MD5('%s') as enc FROM usuarios WHERE usuario='%s' AND activo=1",	clave, usuario)

	dt, err := db.con.MuestraDatosDT(sq)
	if err != nil {
		return "", err
	}
	if len(dt) == 0 {
		return "", fmt.Errorf("usuario no encontrado o inactivo")
	}

	row := dt[0]
	id, _ := row["id"].(int64)
	db.IdUsuario = int(id)
	db.Usuario = row["usuario"].(string)
	db.Nombre = row["nombre"].(string)
	
	activo, _ := row["activo"].(int64)
	db.Activo = activo != 0
	
	db.Clave = row["clave"].(string)
	// Recuperar fecha de expiración de forma segura
	if fechaVal, ok := row["fechaExpiracion"]; ok {
    switch v := fechaVal.(type) {
    case time.Time:
      db.FechaExpiracion = v
    case string:
      // Si por alguna razón viene como string, lo parseamos
      db.FechaExpiracion, _ = time.Parse("2006-01-02", v)
    }
	}
	// Recuperar el MD5 calculado (el campo "enc")
	if encVal, ok := row["enc"]; ok {
    claveEncriptada = encVal.(string)
	}

	return claveEncriptada, nil
}
