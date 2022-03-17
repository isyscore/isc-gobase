package iconv

func Convert(input []byte, output []byte, fromEncoding string, toEncoding string) (bytesRead int, bytesWritten int, err error) {
	// create a temporary converter
	converter, err := NewConverter(fromEncoding, toEncoding)
	if err == nil {
		// call converter's Convert
		bytesRead, bytesWritten, err = converter.Convert(input, output)

		if err == nil {
			var shiftBytesWritten int
			// call Convert with a nil input to generate any end shift sequences
			_, shiftBytesWritten, err = converter.Convert(nil, output[bytesWritten:])
			// add shift bytes to total bytes
			bytesWritten += shiftBytesWritten
		}
		// close the converter
		_ = converter.Close()
	}
	return
}

// ConvertString All in one ConvertString method, rather than requiring the construction of an iconv.Converter
func ConvertString(input string, fromEncoding string, toEncoding string) (output string, err error) {
	// create a temporary converter
	converter, err := NewConverter(fromEncoding, toEncoding)
	if err == nil {
		// convert the string
		output, err = converter.ConvertString(input)
		// close the converter
		_ = converter.Close()
	}
	return
}

func GB2312ToUTF8String(in string) (string, error) {
	return ConvertString(in, "GB2312", "UTF-8")
}

func GB2312ToUTF8(in []byte) ([]byte, error) {
	output := make([]byte, len(in)*2)
	_, outputLen, err := Convert(in, output, "GB2312", "UTF-8")
	if err != nil {
		return nil, err
	}
	return output[0:outputLen], nil
}

func GBKToUTF8(in []byte) ([]byte, error) {
	output := make([]byte, len(in)*2)
	_, outputLen, err := Convert(in, output, "GBK", "UTF-8")
	if err != nil {
		return nil, err
	}
	return output[0:outputLen], nil
}
