package assuan

type ResponseType string

const (
	Ok      ResponseType = "OK"
	Error                = "ERR"
	Status               = "S"
	Comment              = "#"
	Data                 = "D"
	Inquire              = "INQUIRE"
)

type Response struct {
	Type ResponseType
	// OK, Comment
	Comment string
	// Error
	Code        string
	Description string
	// Status
	Status string
	// Inquire, Status
	Keyword string
	// Inquire
	Parameters string
	// Data
	Data []byte
}
