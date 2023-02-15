package logy

func init() {
	loadConfigFromEnv()

	RegisterHandler("console", NewConsoleHandler())
	RegisterHandler("file", NewFileHandler())
}
