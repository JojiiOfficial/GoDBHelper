package godbhelper

//ErrHookFunc error hook
type ErrHookFunc func(error)

//ErrHookOptions options for ErrorHooks
type ErrHookOptions struct {
	ReturnNilOnErr bool
}

func (dbhelper *DBhelper) handleErrHook(err error) error {
	if dbhelper.ErrHookFunc == nil || err == nil {
		return err
	}

	dbhelper.ErrHookFunc(err)
	options := dbhelper.NextHookErrOption
	if options == nil {
		options = dbhelper.ErrHookOptions
	}

	if options != nil && options.ReturnNilOnErr {
		return nil
	}

	return err
}
