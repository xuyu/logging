logging
=======

logging library in golang base on log pkg


features
--------

*support logging level
*support file logger
*support rotation logger(by filename)


example
-------

```go
import "github.com/xuyu/logging"
```

simple usage:

```go
logging.SetLevel(INFO)
logging.SetPrefix("Prefix")
Debug("%d, %s", 1, "OK")
Error("%d, %s", 4, "OK")
```

simple file logger:

```go
l, err := logging.NewFileLogger("/tmp/file.log")
if err != nil {
	panic(err)
}
logging.SetDefaultLogger(l)
...
```

rotation logger:

```go
l, err := logging.NewRotationLogger("/tmp/rotation.log", "060102-15")
if err != nil {
	panic(err)
}
logging.SetDefaultLogger(l)
...
```