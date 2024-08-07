package support

type DataAccessError interface {
	Error() string
}

type PersistenceErrorTranslator interface {
	TranslateExceptionIfPossible(err error) DataAccessError
}
