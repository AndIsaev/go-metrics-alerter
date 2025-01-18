package file

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestFileManager_CreateDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file, _ := os.CreateTemp(tempDir, "tempfile")
	defer file.Close()
	fm := &FileManager{file: file}

	t.Run("succes", func(t *testing.T) {
		newDirPath := tempDir + "/newdir"
		if err := fm.CreateDir(newDirPath); err != nil {
			t.Errorf("Ожидалось, что директория создана успешно, но возникла ошибка: %v", err)
		}
		if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
			t.Errorf("Директория %s должна существовать, но её нет", newDirPath)
		}
	})
	t.Run("dir already exist", func(t *testing.T) {
		newDirPath := tempDir + "/newdir"
		if err := fm.CreateDir(newDirPath); err != nil {
			t.Errorf("Ожидалось, что метод завершится успешно, несмотря на существование директории, но возникла ошибка: %v", err)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		if err := fm.CreateDir("/root/forbidden"); err == nil {
			t.Errorf("Ожидалась ошибка при создании директории в защищённой области, но её не было")
		}
	})
}

func TestFileManager_Overwrite(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name())

	fm := &FileManager{
		file:     file,
		producer: json.NewEncoder(file),
		consumer: json.NewDecoder(file),
	}

	metrics := []common.Metrics{
		{
			ID:    "metric1",
			MType: common.Gauge,
			Value: linkFloat64(42.0),
		},
	}
	t.Run("success overwrite", func(t *testing.T) {
		if err := fm.Overwrite(metrics); err != nil {
			t.Fatalf("Не удалось записать начальные данные в файл: %v", err)
		}

		readMetrics, err := fm.ReadFile()
		assert.NoError(t, err, os.ErrClosed)
		assert.Equal(t, metrics, readMetrics)
	})

	t.Run("error open file", func(t *testing.T) {
		os.Remove(fm.file.Name())
		err := fm.Overwrite(metrics)
		assert.Error(t, err, fmt.Errorf("no such file or directory"))
	})
}

func TestFileManager_ReadFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name())

	metrics := []common.Metrics{
		{
			ID:    "metric1",
			MType: common.Gauge,
			Value: linkFloat64(42.0),
		},
	}
	t.Run("success read file", func(t *testing.T) {
		fm := &FileManager{
			file:     file,
			producer: json.NewEncoder(file),
			consumer: json.NewDecoder(file),
		}

		if _, err := os.OpenFile(fm.file.Name(), os.O_WRONLY, 0666); err != nil {
			t.Fatalf("Не удалось открыть файл: %v", err)
		}
		if err := fm.Overwrite(metrics); err != nil {
			t.Fatalf("Не удалось записать начальные данные в файл: %v", err)
		}

		readMetrics, err := fm.ReadFile()
		assert.NoError(t, err, os.ErrClosed)
		assert.Equal(t, metrics, readMetrics)
	})
	t.Run("error read file", func(t *testing.T) {
		var expected []common.Metrics
		fm := &FileManager{
			file:     file,
			producer: json.NewEncoder(file),
			consumer: json.NewDecoder(file),
		}

		if err := fm.producer.Encode(`{"not":"a list"`); err != nil {
			t.Fatalf("Не удалось записать начальные данные в файл: %v", err)
		}
		result, err := fm.ReadFile()
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestFileManager_Close(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name())

	fm := &FileManager{
		file:     file,
		producer: json.NewEncoder(file),
		consumer: json.NewDecoder(file),
	}

	t.Run("success close file", func(t *testing.T) {
		err := fm.Close()
		assert.NoError(t, err)
	})

	t.Run("doesn't exists file", func(t *testing.T) {
		fm := &FileManager{
			file:     file,
			producer: json.NewEncoder(file),
			consumer: json.NewDecoder(file),
		}

		os.Remove(fm.file.Name())
		err := fm.Close()
		assert.ErrorAs(t, err, &os.ErrClosed)
	})
}
