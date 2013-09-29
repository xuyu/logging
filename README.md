logging
=======

logging library in golang base on log pkg


features
--------

* support logging level

* support file handler

* support rotation handler(by filename)

* support multi handlers


example
-------

```go
import "github.com/xuyu/logging"
```

stdout handler:

```go
logging.EnableDefaultStdout().SetLevel(INFO)
logging.Debug("%d, %s", 1, "OK")
logging.Error("%d, %s", 4, "OK")
```

simple file handler:

```go
l, err := logging.NewFileHandler("/tmp/file.log")
if err != nil {
	panic(err)
}
logging.AddHandler("file", l)
...
```

rotation handler:

```go
l, err := logging.NewRotationHandler("/tmp/rotation.log", "060102-15")
if err != nil {
	panic(err)
}
logging.AddHandler("rotation", l)
...
```

multi handler:

```go
...
logging.EnableDefaultStdout().SetLevel(INFO)
logging.AddHandler("file", l1)
logging.AddHandler("rotation", l2)
...
```