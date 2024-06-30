package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const URL = "User:Password@tcp(localhost:3306)/admision_miercoles"

var db *sql.DB

// Realizar la conexión a la base de datos.
func Connect() {
	connection, err := sql.Open("mysql", URL)

	if err != nil {
		panic(err)
	}

	fmt.Println("Conexión exitosa a la base de datos.")
	db = connection
}

// Cerrar la conexión a la base de datos.
func Close() {
	db.Close()
	fmt.Println("Conexión cerrada.")
}

// Verificar si la conexión a la base de datos está abierta.
func Ping() {
	err := db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Ping exitoso.")
}

// Verificar si una tabla existe en la base de datos.
func TableExists(tableName string) bool {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)

	// Recibe una consulta SQL y luego de ello argumentos indefinidos.
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}	

	defer rows.Close()

	// Va a devolver un valor booleano si existe va a devolver
	// un true y no existe va a devolver un false.
	// Next recorre la tabla entonces si puede recorrer.
	return rows.Next()
}

// Polimorfismo de Query.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		fmt.Println(err)
	}

	return rows, err
}
// Polimorfismo de Exec
func Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}