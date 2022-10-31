package expression

type Expression interface {
	Name() string
	CheckFieldType(filedType string) bool
	Parse(exp string) error
}
