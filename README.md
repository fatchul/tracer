# Momolog
Log library was created by evermos to simplify the need for logs on our service.
Momolog also has a middleware logger that you can use.
Momolog has a wrapped library from the [Zerolog](https://github.com/rs/zerolog) library.

## Features
1. Log
2. Middleware logging
3. Distributed Tracing (using [opentelemetry](https://opentelemetry.io/))
4. Exporter data tracing (OTLP, Jaeger, Prometheus)

## How to Use
```go
package main

import (
	"context"
	"github.com/fatchul/tracer"
)

func main() {
	log := momolog.New()

	type CustomData struct {
		A int
	}
	value := CustomData{A: 12}
	
	log.Fatal(context.Background()).
		Field("key", value).
		Layer("something went wrong")
}
```

Another [examples](https://github.com/fatchul/tracer/tree/main/examples) if you want to check.

## License
[TBD]

## Contributing
[TBD]