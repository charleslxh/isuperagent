package error

type InvalidMiddlewareFactoryArgumentError struct{}

func (e *InvalidMiddlewareFactoryArgumentError) Error() string {
	return "invalid middleware factory arguments"
}
