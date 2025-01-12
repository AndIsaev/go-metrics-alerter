# cmd/profile

В данной директории будет содержаться код для нагрузочного тестирования роута для создания/обновления метрик

в папке profiles лежа два файла ```base.pprof``` и ```result.pprof```
последние оптимизации выдали следующий результат

```text
File: server
Type: inuse_space
Time: Jan 12, 2025 at 5:17pm (MSK)
Duration: 60.01s, Total samples = 24036.05kB 
Showing nodes accounting for -84.34kB, 0.35% of 24036.05kB total
Dropped 1 node (cum <= 120.18kB)
      flat  flat%   sum%        cum   cum%
-11536.86kB 48.00% 48.00% -1536.30kB  6.39%  github.com/AndIsaev/go-metrics-alerter/internal/storage/inmemory.(*MemStorage).UpsertByValue
10000.56kB 41.61%  6.39% 10000.56kB 41.61%  github.com/AndIsaev/go-metrics-alerter/internal/storage/inmemory.(*MemStorage).Create
 7220.69kB 30.04% 23.65%  7769.53kB 32.32%  compress/flate.NewWriter (inline)
-6318.10kB 26.29%  2.64%   902.59kB  3.76%  compress/gzip.(*Writer).Write
-4096.22kB 17.04% 19.68%  -512.03kB  2.13%  net/http.readRequest
 3584.19kB 14.91%  4.77%  3584.19kB 14.91%  net/textproto.(*Reader).ReadLine (inline)
  548.84kB  2.28%  2.48%   548.84kB  2.28%  compress/flate.(*compressor).initDeflate (inline)
  512.56kB  2.13%  0.35%   512.56kB  2.13%  sync.(*Pool).pinSlow
         0     0%  0.35%   512.56kB  2.13%  bufio.(*Writer).Flush
         0     0%  0.35%   902.59kB  3.76%  compress/gzip.(*Writer).Close
         0     0%  0.35%  -633.71kB  2.64%  github.com/AndIsaev/go-metrics-alerter/internal/logger.RequestLogger.func1
         0     0%  0.35%  -633.71kB  2.64%  github.com/AndIsaev/go-metrics-alerter/internal/logger.ResponseLogger.func1
         0     0%  0.35% -1536.30kB  6.39%  github.com/AndIsaev/go-metrics-alerter/internal/service/server.(*Methods).UpdateMetricByValue
         0     0%  0.35%  7769.53kB 32.32%  github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware.(*compressWriter).Close
         0     0%  0.35%  -633.71kB  2.64%  github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware.GzipMiddleware.func1
         0     0%  0.35% -6866.95kB 28.57%  github.com/AndIsaev/go-metrics-alerter/internal/service/server/middleware.GzipMiddleware.func1.deferwrap1
         0     0%  0.35%  -633.71kB  2.64%  github.com/go-chi/chi.(*ChainHandler).ServeHTTP
         0     0%  0.35%  -633.71kB  2.64%  github.com/go-chi/chi.(*Mux).ServeHTTP
         0     0%  0.35%  -633.71kB  2.64%  github.com/go-chi/chi.(*Mux).routeHTTP
         0     0%  0.35%  -633.71kB  2.64%  github.com/go-chi/chi/middleware.Recoverer.func1
         0     0%  0.35%  -633.71kB  2.64%  github.com/go-chi/chi/middleware.StripSlashes.func1
         0     0%  0.35%  -633.71kB  2.64%  main.(*ServerApp).initRouter.SetHeader.func5.1
         0     0%  0.35% -1536.30kB  6.39%  main.(*ServerApp).initRouter.func4.(*Handler).SetMetricHandler.2
         0     0%  0.35%   512.56kB  2.13%  net/http.(*chunkWriter).Write
         0     0%  0.35%   512.56kB  2.13%  net/http.(*chunkWriter).writeHeader
         0     0%  0.35%  -512.03kB  2.13%  net/http.(*conn).readRequest
         0     0%  0.35%  -633.18kB  2.63%  net/http.(*conn).serve
         0     0%  0.35%   512.56kB  2.13%  net/http.(*response).finishRequest
         0     0%  0.35%  -633.71kB  2.64%  net/http.HandlerFunc.ServeHTTP
         0     0%  0.35%   512.56kB  2.13%  net/http.Header.WriteSubset (inline)
         0     0%  0.35%   512.56kB  2.13%  net/http.Header.sortedKeyValues
         0     0%  0.35%   512.56kB  2.13%  net/http.Header.writeSubset
         0     0%  0.35%  -633.71kB  2.64%  net/http.serverHandler.ServeHTTP
         0     0%  0.35%   512.56kB  2.13%  sync.(*Pool).Get
         0     0%  0.35%   512.56kB  2.13%  sync.(*Pool).pin

```
