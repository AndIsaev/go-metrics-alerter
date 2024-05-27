package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize() error {
	cfg := zap.NewProductionConfig()

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter

	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (r *loggingResponseWriter) Done() {
	// if the `w.WriteHeader` wasn't called, set status code to 200 OK
	if r.responseData.status == 0 {
		r.responseData.status = http.StatusOK
	}
}

// ResponseLogger — middleware-логер ответов HTTP-запросов.
func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{ResponseWriter: w, responseData: responseData}
		next.ServeHTTP(&lw, r)
		lw.Done()

		Log.Info("response",
			zap.Int("status", lw.responseData.status),
			zap.Int("size", lw.responseData.size),
		)
	})
}
