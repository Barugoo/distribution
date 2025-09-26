# your-lib

Библиотека на Go для предсказуемого дробления произвольных чисел на заранее заданные доли.

[![Go Reference](https://pkg.go.dev/badge/github.com/barugoo/distribution.svg)](https://pkg.go.dev/github.com/barugoo/distribution)
[![CI](https://github.com/barugoo/distribution/actions/workflows/ci.yml/badge.svg)](https://github.com/barugoo/distribution/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/barugoo/distribution)](https://goreportcard.com/report/github.com/barugoo/distribution)

---

## Установка

```bash
go get github.com/barugoo/distribution@latest
```

---

## Быстрый старт

```go
package main

import (
	"fmt"
	"github.com/barugoo/distribution"
)

func main() {
	c := yourlib.New(yourlib.WithPrefix("demo: "))
	out := c.Do("hello")
	fmt.Println(out)
}
```

---

## Почему your-lib?

* Простое, предсказуемое API
* Минимум зависимостей
* Семантическое версионирование и стабильность API после v1

---

## Документация

* [API на pkg.go.dev](https://pkg.go.dev/github.com/barugoo/distribution)
* Примеры использования: папка [`/examples`](./examples)
* Вопросы и обсуждения: [Discussions](https://github.com/barugoo/distribution/discussions)

---

## Безопасность

См. [SECURITY.md](./SECURITY.md).

---

## Вклад

См. [CONTRIBUTING.md](./CONTRIBUTING.md).
Будем рады любому участию! 🙌

```
