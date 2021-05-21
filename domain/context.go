package domain

func InitialContext() *Context {
	return &Context{
		CurrentTransaction: nil,
	}
}

func NewContext(tx PersistenceProvider) *Context {
	return &Context{
		CurrentTransaction: tx,
	}
}

// Context Enables transaction re-use across domain method calls
type Context struct {
	CurrentTransaction PersistenceProvider
}

// Transaction Create a new transaction or get a reference to the current transaction
func (c *Context)Transaction(callback func(*Context) error) error {
	if c.CurrentTransaction == nil {
		return DB.Transaction(func (tx PersistenceProvider) error {
			newCtx := Context{
				CurrentTransaction: tx,
			}
			return callback(&newCtx)
		})
	} else {
		return callback(c)
	}
}
