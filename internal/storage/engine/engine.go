package engine

type EngineInterface interface {
	Set(string, string) error
	Get(string) (string, bool, error)
	Delete(string) error
	IsSet(string) (bool, error)
	Keys() *[]string
}
