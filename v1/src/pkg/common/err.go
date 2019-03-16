package common

// CheckError tells if the error passed to it is nil
func CheckError(err error) {
	if err != nil {
		Log.Println(err)
	}
}
