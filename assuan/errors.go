package assuan

import "fmt"

type AssuanError struct {
	Code        string
	Description string
}

func (a AssuanError) Error() string {
	return fmt.Sprintf("%s (%s)", a.Description, a.Code)
}
