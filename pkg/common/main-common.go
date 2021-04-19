package common

// FailOnError will panic if err
func FailOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type FailOnErrorFunc func(err error)
