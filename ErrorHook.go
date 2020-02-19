package godbhelper

//ErrHookFunc error hook 1. string is the query; 2. string is the hookPrefix if set
type ErrHookFunc func(error, string, string)

//ErrHookOptions options for ErrorHooks
type ErrHookOptions struct {
	ReturnNilOnErr bool
	Prefix         string
}

func (dbhelper *DBhelper) handleErrHook(err error, content string) error {
	if dbhelper.ErrHookFunc == nil || err == nil {
		return err
	}

	//Use the correct options
	options := dbhelper.NextErrHookOption
	if options == nil {
		options = dbhelper.ErrHookOptions
	} else {
		dbhelper.NextErrHookOption = nil
	}

	//Call the correct hook
	if dbhelper.NextErrHookFunc == nil {
		dbhelper.ErrHookFunc(err, content, options.Prefix)
	} else {
		dbhelper.NextErrHookFunc(err, content, options.Prefix)
		dbhelper.NextErrHookFunc = nil
	}

	if options != nil && options.ReturnNilOnErr {
		return nil
	}

	return err
}
