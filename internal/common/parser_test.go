package common

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromJSON(t *testing.T) {
	// Создаём временный JSON файл для тестирования
	jsonContent := `{
	  "key1": "value1",
	  "key2": 2,
	  "duration": "2s"
	 }`

	// Создаём временный файл
	tempFile, err := os.CreateTemp("", "test_config_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	_, err = tempFile.Write([]byte(jsonContent))
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	// Определяем ключи для обработки
	keys := []string{"duration"} // Пример ключа, который обрабатывается в ParseAndSetDuration

	// Вызываем функцию
	resultJson, err := LoadConfigFromJSON(tempFile.Name(), keys)
	assert.NoError(t, err)

	// Проверяем результаты
	var resultData map[string]interface{}
	err = json.Unmarshal(resultJson, &resultData)
	assert.NoError(t, err)

	duration, _ := time.ParseDuration("2s")
	expectedValues := map[string]interface{}{
		"key1":     "value1",
		"key2":     float64(2),
		"duration": float64(duration),
	}

	assert.Equal(t, expectedValues, resultData)
}

func TestLoadConfigFromJSONError(t *testing.T) {
	// Создаём временный JSON файл для тестирования
	jsonContent := `{
	  "key1": "value1",
	  "key2": 2,
	  "duration": "ddsds"
	 }`

	// Создаём временный файл
	tempFile, err := os.CreateTemp("", "test_config_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	_, err = tempFile.Write([]byte(jsonContent))
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	// Определяем ключи для обработки
	keys := []string{"duration"} // Пример ключа, который обрабатывается в ParseAndSetDuration

	// Вызываем функцию
	_, err = LoadConfigFromJSON(tempFile.Name(), keys)
	assert.Error(t, err, "invalid duration format for 'duration': time: invalid duration \"ddsds\"")
}

func TestParseAndSetDuration(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "valid duration string",
			key:     "duration",
			value:   "2s",
			want:    2 * time.Second,
			wantErr: false,
		},
		{
			name:    "invalid duration string",
			key:     "duration",
			value:   "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-string value",
			key:     "duration",
			value:   123, // некорректное значение типа
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]any)
			err := ParseAndSetDuration(tt.key, tt.value, result)

			if tt.wantErr {
				assert.Error(t, err, fmt.Sprintf("expected error for test '%s'", tt.name))
			} else {
				assert.NoError(t, err, fmt.Sprintf("unexpected error for test '%s'", tt.name))
				actual, ok := result[tt.key].(time.Duration)
				assert.True(t, ok, "value should be of type time.Duration")
				assert.Equal(t, tt.want, actual, "duration values should be equal")
			}
		})
	}
}
