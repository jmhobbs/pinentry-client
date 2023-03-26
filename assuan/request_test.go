package assuan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Escape(t *testing.T) {
	tests := []struct {
		Name     string
		Input    []byte
		Expected []byte
	}{
		{
			"Line Feed",
			[]byte("Hello\nWorld"),
			[]byte("Hello%0AWorld"),
		},
		{
			"Carriage Return",
			[]byte("HelloWorld\r"),
			[]byte("HelloWorld%0D"),
		},
		{
			"Percent",
			[]byte("Hello%World"),
			[]byte("Hello%25World"),
		},
		{
			"Nil",
			nil,
			nil,
		},
		{
			"Complete",
			[]byte("Hello World\r\nI have 5% battery."),
			[]byte("Hello World%0D%0AI have 5%25 battery."),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, Escape(test.Input), test.Name)
	}
}

func Test_RequestOption(t *testing.T) {
	assert.Equal(t, []byte("OPTION color=blue"), RequestOption("color", "blue")())
}

func Test_RequestData(t *testing.T) {
	assert.Equal(t, []byte("D -100"), RequestData([]byte("-100"))())
}
