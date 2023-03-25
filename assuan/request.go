package assuan

import "bytes"

type Request func() []byte

func RequestGeneric(command string, parameters []byte) Request {
	var msg []byte
	if len(parameters) == 0 {
		msg = []byte(command)
	} else {
		buf := bytes.NewBufferString(command)
		buf.Write([]byte{' '})
		buf.Write(parameters)
		msg = buf.Bytes()
	}

	return func() []byte {
		return msg
	}
}

var (
	RequestBye   = RequestGeneric("BYE", nil)
	RequestReset = RequestGeneric("RESET", nil)
	RequestEnd   = RequestGeneric("END", nil)
	RequestHelp  = RequestGeneric("HELP", nil)
	RequestQuit  = RequestGeneric("QUIT", nil)
	RequestNOP   = RequestGeneric("NOP", nil)
)

func RequestOption(name, value string) []byte {
	buf := bytes.NewBufferString("OPTION")
	buf.Write([]byte{' '})
	buf.WriteString(name)
	if value != "" {
		buf.WriteRune('=')
		buf.WriteString(value)
	}
	return buf.Bytes()
}

/*
Sends raw data to the server.
There must be exactly one space after the ’D’.
The values for ’%’, CR and LF must be percent escaped.
These are encoded as %25, %0D and %0A, respectively.
Only uppercase letters should be used in the hexadecimal representation.
Other characters may be percent escaped for easier debugging.
All Data lines are considered one data stream up to the OK or ERR response.
Status and Inquiry Responses may be mixed with the Data lines.
*/
func RequestData(data []byte) []byte {
	buf := bytes.NewBufferString("D")
	buf.Write([]byte{' '})
	buf.Write(data)
	return buf.Bytes()
}
