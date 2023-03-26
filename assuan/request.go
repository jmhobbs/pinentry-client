package assuan

import "bytes"

type Request func() []byte

func Escape(input []byte) []byte {
	if input == nil {
		return nil
	}
	return bytes.ReplaceAll(
		bytes.ReplaceAll(
			bytes.ReplaceAll(
				input,
				[]byte{'%'},
				[]byte{'%', '2', '5'},
			),
			[]byte{'\n'},
			[]byte{'%', '0', 'A'},
		),
		[]byte{'\r'},
		[]byte{'%', '0', 'D'},
	)
}

func RequestGeneric(command string, parameters []byte) Request {
	var msg []byte
	if len(parameters) == 0 {
		msg = []byte(command)
	} else {
		buf := bytes.NewBufferString(command)
		buf.Write([]byte{' '})
		buf.Write(Escape(parameters))
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

func RequestOption(name, value string) Request {
	buf := bytes.NewBufferString("OPTION")
	buf.Write([]byte{' '})
	buf.WriteString(name)
	if value != "" {
		buf.WriteRune('=')
		buf.WriteString(value)
	}
	return func() []byte {
		return buf.Bytes()
	}
}

func RequestData(data []byte) Request {
	buf := bytes.NewBufferString("D")
	buf.Write([]byte{' '})
	buf.Write(Escape(data))
	return func() []byte {
		return buf.Bytes()
	}
}
