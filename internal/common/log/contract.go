//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock
package log

const (
	kindString = iota
	kindError
)

type Field struct {
	kind  int64
	key   string
	value any
}

type Logger interface {
	Error(msg string, fields ...Field)
}

func String(key, value string) Field {
	return Field{
		kind:  kindString,
		key:   key,
		value: value,
	}
}

func Error(err error) Field {
	return Field{
		kind:  kindError,
		key:   "error",
		value: err,
	}
}
