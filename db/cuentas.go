package db

import (
	"fmt"
)

type DBCuentas struct {
	con             *Conexion
	Numero       		string
	IngresoEgreso   bool
	Imputable       bool
	Nombre          string
	UsuarioUM       int
}

func NewDBCuentas() (*DBCuentas, error) {
	con, err := GetInstance()
	if err != nil {
		return nil, err
	}
	return &DBCuentas{
		con: con,
	}, nil
}

func (db *DBCuentas) ListaCuentas() ([]*DBCuentas, error) {
	query := "SELECT numero, ingresoEgreso, imputable, nombre FROM cuentas ORDER BY numero"
	rows, err := db.con.MuestraDatosDT(query)

	if err != nil {
		return nil, err
	}

	var cuentas []*DBCuentas
	for _, fila := range rows {
		cuenta := &DBCuentas{
			Numero:        fmt.Sprintf("%v", fila["numero"]),
			IngresoEgreso: fila["ingresoEgreso"].(int64) == 1 ,
			Imputable:     fila["imputable"].(int64) == 1 ,
			Nombre:        fmt.Sprintf("%v", fila["nombre"]),
		}
		cuentas = append(cuentas, cuenta)
	}
	return cuentas, nil
}

// AgregarCuenta inserta un nuevo registro
func (db *DBCuentas) AgregarCuenta() error {
	sq := fmt.Sprintf("INSERT INTO cuentas(numero, nombre, ingresoEgreso, imputable, usuarioUM) VALUES ('%s', '%s', %t, %t, %d)",
		db.Numero, db.Nombre, db.IngresoEgreso, db.Imputable, db.UsuarioUM)
	return db.con.Transaccion(sq)
}

// ModificarCuenta actualiza un registro existente
func (db *DBCuentas) ModificarCuenta(CuentaModificada string) error {
	sq := fmt.Sprintf("UPDATE cuentas SET numero='%s', nombre='%s', ingresoEgreso=%t, imputable=%t, usuarioUM=%d WHERE numero='%s'",
		db.Numero, db.Nombre, db.IngresoEgreso, db.Imputable, db.UsuarioUM, CuentaModificada)
	return db.con.Transaccion(sq)
}