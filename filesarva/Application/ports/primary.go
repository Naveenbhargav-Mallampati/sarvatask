package ports

type WebApiPort interface {
	ProcessFile(filepath string)
}
