package service

// Resize images that have been uploaded
func Resize(v interface{}) {

}

// returns nil  if format isn't accepted
// func findFormat(data string) error {
// 	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
// 	_, format, err := image.DecodeConfig(reader)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if format != "jpeg" || format != "png" || format != "jpg" {
// 		var err error
// 		return err
// 	}
// 	return nil
// }
