package db

import "fmt"

type DBCooperativa struct {
	con *Conexion
	/*
		IdUsuario       int
		Usuario         string
		Nombre          string
		Activo          bool
		Clave           string
		FechaExpiracion time.Time
		UsuarioUM       int
	*/
}

type Cuenta struct {
	Numero string
	Nombre string
}

// NewDBUsuarios constructor
func NewDBCooperativo() (*DBCooperativa, error) {
	con, err := GetInstance()
	if err != nil {
		return nil, err
	}
	return &DBCooperativa{
		con: con,
	}, nil
}

func (db *DBCooperativa) ObtenerCuentaContable() ([]Cuenta, error) {
	consulta := `SELECT numero, nombre FROM cuentas 
	          WHERE numero LIKE "1-04-%" AND imputable = 1 
	          ORDER BY numero`
	filas, err := db.con.MuestraDatosDT(consulta)
	if err != nil {
		return nil, err
	}
	var cuentas []Cuenta
	for _, fila := range filas {
		cuenta := Cuenta{
			Numero:        fmt.Sprintf("%v", fila["numero"]),
			Nombre:        fmt.Sprintf("%v", fila["nombre"]),
		}
		cuentas = append(cuentas, cuenta)
	}
	return cuentas, nil
}