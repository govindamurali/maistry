package maistry

type ILogger interface {
	Error(message string, err error, args map[string]interface{})
	Trace(message string, args map[string]interface{})
}
