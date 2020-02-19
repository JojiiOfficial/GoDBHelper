package godbhelper

//ErrHookFunc error hook 1. string is the query; 2. string is the hookPrefix if set
type ErrHookFunc func(error, string, string)

//ErrHookOptions options for ErrorHooks
type ErrHookOptions struct {
	ReturnNilOnErr bool
}

func (dbhelper *DBhelper) handleErrHook(err error, content string) error {
	if dbhelper.ErrHookFunc == nil || err == nil {
		return err
	}

	nextPrefix := ""
	//Use nextPrefix if not empty
	if dbhelper.NextLogPrefix != nil {
		nextPrefix = *dbhelper.NextLogPrefix
		//Reset nextPrefix
		dbhelper.NextLogPrefix = nil
	}

	//Call the correct hook
	if dbhelper.NextErrHookFunc == nil {
		dbhelper.ErrHookFunc(err, content, nextPrefix)
	} else {
		dbhelper.NextErrHookFunc(err, content, nextPrefix)
		dbhelper.NextErrHookFunc = nil
	}

	//Use the correct options
	options := dbhelper.NextErrHookOption
	if options == nil {
		options = dbhelper.ErrHookOptions
	} else {
		dbhelper.NextErrHookOption = nil
	}

	if options != nil && options.ReturnNilOnErr {
		return nil
	}

	return err
}
