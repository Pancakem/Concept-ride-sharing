package common

func CheckError(err error) {
	if err != nil {
		Log.Println(err)
	}
}
