# logger
--
    import "."


## Usage

#### func  NewLogger

```go
func NewLogger(field string, level logrus.Level) *logrus.Entry
```
NewLogger Create a Logrus logger with the nested-logrus-formatter format

#### type Log

```go
type Log struct {
	LogLevel logrus.Level `yaml:"logLevel"`
}
```

Log yalm structure
