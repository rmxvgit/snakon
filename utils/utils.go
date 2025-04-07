package utils

func INeedThat(a any) {
	return
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
