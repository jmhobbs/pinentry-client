package assuan

import (
	"bytes"
	"fmt"
	"io"
)

type Client struct {
	source io.ReadCloser
	sink   io.WriteCloser
}

func New(r io.ReadCloser, w io.WriteCloser) *Client {
	return &Client{source: r, sink: w}
}

func (c *Client) Close() {
	c.Write(RequestBye)
	c.Read()
	c.sink.Close()
	c.source.Close()
}

func (c *Client) Write(req Request) error {
	_, err := c.sink.Write(req())
	if err != nil {
		return err
	}
	_, err = c.sink.Write([]byte{'\n'})
	return err
}

func (c *Client) Read() (Response, error) {
	resp := Response{}

	buf := make([]byte, 1000)
	_, err := c.source.Read(buf)
	if err != nil {
		return resp, err
	}

	// rough and ready response parsing, we ain't fancy
	split := bytes.SplitN(buf[:bytes.IndexRune(buf, '\n')], []byte{' '}, 2)

	switch string(split[0]) {
	case "OK":
		resp.Type = Ok
		resp.Comment = string(remainingParameters(split))
	case "#":
		resp.Type = Comment
		resp.Comment = string(remainingParameters(split))
	case "ERR":
		resp.Type = Error
		resp.Code, resp.Description = innerSplit(split)
	case "D":
		resp.Type = Data
		resp.Data = remainingParameters(split)
	case "S":
		resp.Type = Status
		resp.Status, resp.Keyword = innerSplit(split)
	case "INQUIRE":
		resp.Type = Inquire
		resp.Keyword, resp.Parameters = innerSplit(split)
	default:
		return resp, fmt.Errorf("unknown command: %q", string(split[0]))
	}

	return resp, nil
}

func remainingParameters(in [][]byte) []byte {
	if len(in) == 2 {
		return in[1]
	}
	return []byte{}
}

func innerSplit(in [][]byte) (a, b string) {
	if len(in) == 2 {
		pSplit := bytes.SplitN(in[1], []byte{' '}, 2)
		a = string(pSplit[0])
		if len(pSplit) == 2 {
			b = string(pSplit[1])
		}
	}
	return a, b
}
