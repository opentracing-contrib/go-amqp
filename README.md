# go-amqp

[![Build Status]](https://travis-ci.org/opentracing-contrib/go-amqp)

[AMQP] instrumentation in Go

For documentation on the packages, [check godoc].

**The APIs in the various packages are experimental and may change in
the future. You should vendor them to avoid spurious breakage.**

## Packages

Instrumentation is provided for the following packages, with the
following caveats:

- **github.com/streadway/amqp**: Client and server instrumentation. *Only supported
  with Go 1.7 and later.*

[AMQP]: https://github.com/streadway/amqp
[Build Status]: https://travis-ci.org/opentracing-contrib/go-amqp.svg
[check godoc]: https://godoc.org/github.com/opentracing-contrib/go-amqp/amqptracer
