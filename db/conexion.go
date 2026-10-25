package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// Conexion estructura que emula tu clase C#
type Conexion struct {
	db *sql.DB
}

var (
	instance *Conexion
	once     sync.Once
	connError error
)

// GetInstance implementa el patrón Singleton (reemplaza al lock de C#)
func GetInstance() (*Conexion, error) {
	once.Do(func() {
		// Configuración (podés usar variables de entorno aquí)
		usuario := "root"
		clave := "total"
		servidor := "localhost"
		base := "finca"

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", usuario, clave, servidor, base)
		
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			connError = fmt.Errorf("error de configuración: %v", err)
			return
		}

		// Verificar la conexión inicial
		if err := db.Ping(); err != nil {
			connError = fmt.Errorf("no se pudo conectar al servidor %s: %v", servidor, err)
			return
		}

		instance = &Conexion{db: db}
	})
	if connError != nil {
		// Reseteamos el 'once' para que el usuario pueda reintentar si arregla la DB
		// (Opcional, pero recomendado en apps de escritorio)
		if instance == nil { once = sync.Once{} } 
		return nil, connError
	}
	return instance, nil
}

// Transaccion ejecuta un comando que no devuelve filas (INSERT, UPDATE, DELETE)
func (c *Conexion) Transaccion(sqlStr string) error {
	_, err := c.db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("error en comando sql: %s - %v", sqlStr, err)
	}
	return nil
}

// MultipleTransaccion maneja un slice de queries dentro de una transacción real
func (c *Conexion) MultipleTransaccion(consultas []string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	for _, query := range consultas {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("error en comando sql: %s - %v", query, err)
		}
	}
	return tx.Commit()
}

// MuestraDatosDT emula el DataTable devolviendo un slice de mapas
func (c *Conexion) MuestraDatosDT(consulta string) ([]map[string]interface{}, error) {
	rows, err := c.db.Query(consulta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		// Magia para leer columnas dinámicamente
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		result = append(result, rowMap)
	}
	return result, nil
}

// Desconectar cierra la base de datos
func (c *Conexion) Desconectar() {
	if c.db != nil {
		c.db.Close()
	}
}

// GetFechaServidor devuelve la fecha y hora actual del servidor MySQL
func (c *Conexion) GetFechaServidor() (string, error) {
	var fechaHora string
	// Usamos DATE_FORMAT para recibirlo ya como un string amigable
	query := "SELECT DATE_FORMAT(NOW(), '%d/%m/%Y %H:%i') as fecha"
	
	err := c.db.QueryRow(query).Scan(&fechaHora)
	if err != nil {
		return "", fmt.Errorf("error al obtener hora del servidor: %v", err)
	}
	return fechaHora, nil
}
