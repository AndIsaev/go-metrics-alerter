package main

import (
	"fmt"
	"math/rand"
	"os"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randomString генерирует случайную строку заданной длины
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

// generateRandomMetrics генерирует срез строк с рандомными значениями метрик
func generateRandomMetrics(count int) []string {
	metrics := make([]string, count)
	for i := 0; i < count; i++ {
		if rand.Intn(2) == 0 {
			// Generate counter
			name := "testCounter" + randomString(5)
			value := rand.Int63n(1000)
			metrics[i] = fmt.Sprintf("counter,%s,%d", name, value)
		} else {
			// Generate gauge
			name := "testGauge" + randomString(5)
			value := rand.Float64() * 100
			metrics[i] = fmt.Sprintf("gauge,%s,%.2f", name, value)
		}
	}
	return metrics
}

func writeToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
