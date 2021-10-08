package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ErrEmptyValue = errors.New("empty config value")

const _mapnikConfigPathENV = "MAPNIK_CONFIG_PATH"

func MapnikStylesheet() (string, error) {
	replacers := []struct {
		target string
		value  string
	}{
		{target: "{{DB_HOST}}", value: viper.GetString(DbHost)},
		{target: "{{DB_PORT}}", value: viper.GetString(DbPort)},
		{target: "{{DB_USER}}", value: viper.GetString(DbUser)},
		{target: "{{DB_PASSWORD}}", value: viper.GetString(DbPassword)},
		{target: "{{DB_NAME}}", value: viper.GetString(DbName)},
	}

	mapnikConfigPath := os.Getenv(_mapnikConfigPathENV)
	if mapnikConfigPath == "" {
		return "", fmt.Errorf("empty mapnik config path: %w", ErrEmptyValue)
	}

	data, err := os.ReadFile(mapnikConfigPath)
	if err != nil {
		return "", fmt.Errorf("read mapnik config: %w", err)
	}

	dataString := string(data)

	for _, r := range replacers {
		dataString = strings.ReplaceAll(dataString, r.target, r.value)
	}

	if err := os.WriteFile(mapnikConfigPath, []byte(dataString), os.ModePerm); err != nil {
		return "", fmt.Errorf("rewrite mapnik config: %w", err)
	}

	return mapnikConfigPath, nil
}

func GetFileCacheDir() (string, error) {
	const op = "GetFileCacheDir"

	value := viper.GetString(FileCacheDir)
	if value == "" {
		return "", fmt.Errorf("%s: %w", op, ErrEmptyValue)
	}

	return value, nil
}

func GetMemoryCacheMaxSize() (int, error) {
	const op = "GetMemoryCacheMaxSize"

	value := viper.GetInt(MemoryCacheMaxSize)
	if value == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyValue)
	}

	return value, nil
}

func GetMemoryCacheLifeWindow() (time.Duration, error) {
	const op = "GetMemoryCacheLifeWindow"

	value := viper.GetDuration(MemoryCacheLifeWindow)
	if value == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyValue)
	}

	return value, nil
}

func GetRendererQueueCapacity() (int, error) {
	const op = "GetRendererQueueCapacity"

	value := viper.GetInt(RendererQueueCapacity)
	if value == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyValue)
	}

	return value, nil
}

func GetRendererWorkersCount() (int, error) {
	const op = "GetRendererWorkersCount"

	value := viper.GetInt(RendererWorkersCount)
	if value < 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyValue)
	}

	return value, nil
}

func InitWorkersCountSetter(cb func(count int) error) {
	viper.OnConfigChange(func(in fsnotify.Event) {
		_ = cb(viper.GetInt(RendererWorkersCount))
	})
}
