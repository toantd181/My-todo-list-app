package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Environment string
	Port        string
	DataDir     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DataPath string
}

func Load() *Config {
	DataDir := getDataDirectory()

	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("PORT", "8080"),
			DataDir:     dataDir,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "todo_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			DataPath: filepath.Join(dataDir, "todos.db"),
		},
	}
}

func getEnv(ket, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDataDirectory() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "TodoApp")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "TodoApp")

	case "linux":
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, "TodoApp")
		}
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".local", "share", "TodoApp")

	case "android":
		return "/data/data/com.yourapp.todo/files"

	default:
		return "./data"
	}

}
