package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	length := 10
	rs := randomString(length)

	assert.Equal(t, length, len(rs), "Random string should have the correct length")

	for _, char := range rs {
		assert.True(t, strings.Contains(letters, string(char)),
			"Random string should only contain alphabetic characters")
	}
}

func TestGenerateRandomMetrics(t *testing.T) {
	count := 5
	metrics := generateRandomMetrics(count)

	assert.Equal(t, count, len(metrics), "Should generate the correct number of metrics")

	for _, metric := range metrics {
		split := strings.Split(metric, ",")
		assert.True(t, len(split) == 3, "Each metric should have three components")
		assert.True(t, split[0] == "counter" || split[0] == "gauge", "Metric type should be counter or gauge")

		name := split[1]
		assert.True(t, strings.HasPrefix(name, "testCounter") || strings.HasPrefix(name, "testGauge"),
			"Metric name should start with the correct prefix")

		if split[0] == "counter" {
			_, err := fmt.Sscanf(split[2], "%d", new(int64))
			assert.NoError(t, err, "Counter value should be an integer")
		} else if split[0] == "gauge" {
			_, err := fmt.Sscanf(split[2], "%f", new(float64))
			assert.NoError(t, err, "Gauge value should be a float")
		}
	}
}

func TestWriteToFile(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	filename := "testfile.txt"

	err := writeToFile(filename, lines)
	assert.NoError(t, err, "writeToFile should not return an error")

	defer os.Remove(filename) // Clean up

	content, err := os.ReadFile(filename)
	assert.NoError(t, err, "ReadFile should not return an error")

	contentLines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.Equal(t, lines, contentLines, "File content should match the written lines")
}
