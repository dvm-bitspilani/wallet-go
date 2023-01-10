package exceptions

// TODO:	Find a more elegant way to do this

type Exception struct {
	Message string
}

func (e *Exception) Error() string {
	return e.Message
}

type UserDisabledException struct {
	Exception
}

type InvalidOptionException struct {
	Exception
}

type ForbiddenTransactionException struct {
	Exception
}

type AuthorityException struct {
	Exception
}

type VendorClosedException struct {
	Exception
}

type AvaliabilityException struct {
	Exception
}

type IncorrectVendorException struct {
	Exception
}

type NotAllowedException struct {
	Exception
}

type MaxStatusException struct {
	Exception
}

type InvalidOrderException struct {
	Exception
}

type InvalidChecksumException struct {
	Exception
}

type TransactionAlreadyEncashedException struct {
	Exception
}
