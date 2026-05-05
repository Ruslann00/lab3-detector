# Лабораторна робота №3

## Тема

Діагностика та рефакторинг високопродуктивних систем у Go.

## Мета роботи

Навчитися виявляти приховані дефекти у високонавантажених системах:
витоки пам'яті, race condition, deadlock-подібні проблеми та неефективні алгоритми
за допомогою pprof, go test -race, benchmark та Delve.

---

## Етап 1. Heap Profiling

Для аналізу витоку пам'яті було використано pprof.

Сервіс було запущено командою:

```bash
go run ./cmd/service

Перший heap profile:
go tool pprof -output=profiles/heap1.pb http://localhost:6060/debug/pprof/heap

Другий heap profile через 2 хвилини:
go tool pprof -output=profiles/heap2.pb http://localhost:6060/debug/pprof/heap

Порівняння:
go tool pprof -http=:8080 -base=profiles/heap1.pb.gz profiles/heap2.pb

Виявлена проблема
У функції processImage було знайдено рядок:
LeakCache[key] = make([]byte, 1024*10)

Глобальна map LeakCache постійно збільшувалася, але старі елементи не видалялися.
Це призводило до поступового зростання використання оперативної пам'яті.

Виправлення
Було створено обмежений кеш SafeCache з максимальним розміром.
Також доступ до кешу було захищено через sync.Mutex.

Етап 2. Race Condition
Для пошуку гонитви було використано:
go test -race ./internal/stats

Проблема
Початковий код:
GlobalStats[imageType]++

Операція виконувалася з багатьох горутин без синхронізації.
Це створювало race condition.

Виправлення
Було створено структуру SafeStatsCounter, яка використовує sync.RWMutex.

Для запису використовується:
mu.Lock()

Для читання використовується:
mu.RLock()

Після виправлення команда:
go test -race ./internal/stats
не показує DATA RACE.

Етап 3. CPU Profiling
Для аналізу CPU було використано benchmark та CPU profile.

Benchmark:
go test -bench=. ./internal/processor

CPU profile до оптимізації:
go test -bench=BenchmarkProcessImageBefore -cpuprofile=profiles/cpu_before.pb.gz ./internal/processor

CPU profile після оптимізації:
go test -bench=BenchmarkProcessImageAfter -cpuprofile=profiles/cpu_after.pb.gz ./internal/processor
Проблема

У початковій версії використовувався код:
regexp.MatchString(`^image_data_\d+_timestamp_\d+$`, data)
Це призводило до повторної компіляції регулярного виразу при кожному виклику.

Виправлення

Було використано:
var optimizedRegexp = regexp.MustCompile(`^image_data_\d+_timestamp_\d+$`)
Після цього регулярний вираз компілюється один раз, а не в кожній ітерації.
Benchmark показав прискорення більше ніж у 2 рази.

Етап 4. Remote Debug
Для remote debugging було використано Delve у headless-режимі.

Команда запуску:
dlv debug ./cmd/service --headless --listen=:40000 --api-version=2 --accept-multiclient

Після цього VS Code було підключено до процесу через конфігурацію Attach to Delve Headless.

