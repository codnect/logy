package logy

func init() {
	RegisterHandler("console", NewConsoleHandler())
}
