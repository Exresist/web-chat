package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"webChat/internal/model"
)

type (
	Config struct {
		DB        Database `json:"db"`
		API       API      `json:"api"`
		Files     Files    `json:"files"`
		SecretKey string   `json:"secret_key"`
	}
	Database struct {
		URL      string `json:"url"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     string `json:"port"`
		SSLMode  string `json:"ssl_mode"`
		User     string `json:"user"`

		MaxOpenConns int `json:"max_open_conns"`
		MaxIdleConns int `json:"max_idle_conns"`

		MigrateDown bool `json:"migrate_down"`
	}
	API struct {
		Address         string         `json:"address"`
		ReadTimeout     model.Duration `json:"read_timeout"`
		WriteTimeout    model.Duration `json:"write_timeout"`
		ShutdownTimeout model.Duration `json:"shutdown_timeout"`
	}
	Files struct {
		MaxFileSize      int64    `json:"max_file_size"`
		AllowedFileTypes []string `json:"allowed_file_type"`
	}
)

func (db *Database) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode,
	)
}

func NewConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}

	var f io.ReadCloser
	f, err := os.Open(cfgPath)
	if err != nil {
		return cfg, os.ErrNotExist
	}

	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshalling: %s", err)
	}

	cfg.DB.URL = cfg.DB.GetURL()

	cfg.Files.MaxFileSize = 4 * (1 << (10 * 2)) // 4MB
	cfg.Files.AllowedFileTypes = []string{"image/jpeg", "image/jpg", "image/png"}

	return cfg, nil
}
