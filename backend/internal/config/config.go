package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config хранит всю конфигурацию приложения.
// Теги `mapstructure` нужны, чтобы viper мог корректно
// сопоставить поля из YAML с полями структуры.
type Config struct {
	HTTP struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"http"`
	Postgres struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"postgres"`
}

// New создает новый экземпляр Config, читая данные из config/config.yaml.
// Также он настроен на переопределение значений через переменные окружения.
func New() (*Config, error) {
	vp := viper.New()
	var cfg Config

	// 1. Установка пути и имени файла конфигурации.
	vp.AddConfigPath("./config") // Путь к директории с конфигом
	vp.SetConfigName("config")   // Имя файла без расширения
	vp.SetConfigType("yaml")

	// 2. Чтение файла конфигурации.
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	// 3. Настройка возможности переопределения через переменные окружения.
	// Например, HTTP_PORT переопределит cfg.HTTP.Port
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.AutomaticEnv()

	// 4. Анмаршалинг данных в структуру Config.
	if err := vp.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
