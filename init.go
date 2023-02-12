package logy

func init() {
	RegisterHandler("console", NewConsoleHandler())
	RegisterHandler("file", NewFileHandler())
}
