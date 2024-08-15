package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func craftDSNs() (string, string, error) {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbName == "" || dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" {
		return "", "", fmt.Errorf("Variables de entorno requeridas no están configuradas")
	}
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	return connection, dsn, nil
}

func ConnectToDB() error {
	connection, dsn, err := craftDSNs()
	if err != nil {
		return fmt.Errorf("Error al construir cadenas de conexión: %w", err)
	}

	// Conexión inicial para crear la base de datos si no existe
	tempDB, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Error al abrir la conexión a MySQL: %w", err)
	}

	// Crear la base de datos si no existe
	createDatabaseCommand := "CREATE DATABASE IF NOT EXISTS " + os.Getenv("DB_NAME") + ";"
	if err := tempDB.Exec(createDatabaseCommand).Error; err != nil {
		return fmt.Errorf("Error al crear la base de datos: %w", err)
	}

	// Cerrar la conexión temporal
	sqlDB, err := tempDB.DB()
	if err != nil {
		return fmt.Errorf("Error al obtener el objeto SQL DB: %w", err)
	}
	sqlDB.Close()

	// Conectar a la base de datos específica
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Error al conectar a la base de datos: %w", err)
	}
	DB = db
	return nil
}
