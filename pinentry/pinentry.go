package pinentry

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/jmhobbs/pinentry-client/assuan"
)

type PinEntry struct {
	bin    string
	cmd    *exec.Cmd
	client *assuan.Client
	queue  []assuan.Request
}

func (p *PinEntry) setStr(command, value string) *PinEntry {
	p.queue = append(p.queue, assuan.RequestGeneric(command, []byte(value)))
	return p
}

func (p *PinEntry) SetTimeout(seconds int) *PinEntry {
	p.queue = append(p.queue, assuan.RequestGeneric("SETTIMEOUT", []byte(strconv.Itoa(seconds))))
	return p
}

func (p *PinEntry) SetPrompt(prompt string) *PinEntry {
	return p.setStr("SETPROMPT", prompt)
}

func (p *PinEntry) SetDescription(description string) *PinEntry {
	return p.setStr("SETDESC", description)
}

func (p *PinEntry) SetTitle(title string) *PinEntry {
	return p.setStr("SETTITLE", title)
}

func (p *PinEntry) SetButtonOk(value string) *PinEntry {
	return p.setStr("SETOK", value)
}

func (p *PinEntry) SetButonCancel(value string) *PinEntry {
	return p.setStr("SETCANCEL", value)
}

func (p *PinEntry) SetButonNotOk(value string) *PinEntry {
	return p.setStr("SETNOTOK", value)
}

func (p *PinEntry) Close() {
	if p.cmd != nil {
		p.client.Close()
		p.cmd.Wait()
	}
}

func New(bin string) *PinEntry {
	return &PinEntry{bin: bin}
}

func (p *PinEntry) spawn() error {
	cmd := exec.Command(p.bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	p.cmd = cmd
	p.client = assuan.New(stdout, stdin)

	resp, err := p.client.Read()
	if err != nil {
		return err
	}
	if resp.Type != assuan.Ok {
		return fmt.Errorf("server did not say hello")
	}

	return nil
}

func (p *PinEntry) executeQueue() error {
	if p.cmd == nil {
		err := p.spawn()
		if err != nil {
			return err
		}
	}

	// run through the queue of pending commands
	for _, msg := range p.queue {
		err := p.client.Write(msg)
		if err != nil {
			return err
		}
		// we expect 1-to-1 writes to reads for most commands
		response, err := p.client.Read()
		if err != nil {
			return err
		}
		if response.Type != assuan.Ok {
			return fmt.Errorf("unexpected response")
		}
	}

	return nil
}

func (p *PinEntry) GetPIN() (string, error) {
	err := p.executeQueue()
	if err != nil {
		return "", err
	}

	err = p.client.Write(assuan.RequestGeneric("GETPIN", []byte{}))
	if err != nil {
		return "", err
	}
	response, err := p.client.Read()
	if err != nil {
		return "", err
	}

	if response.Type == assuan.Error {
		return "", NewPinentryError(response)
	}
	if response.Type != assuan.Data {
		return "", fmt.Errorf("unexpected response")
	}

	return string(response.Data), nil
}

func (p *PinEntry) Confirm() (bool, error) {
	err := p.executeQueue()
	if err != nil {
		return false, err
	}

	err = p.client.Write(assuan.RequestGeneric("CONFIRM", []byte{}))
	if err != nil {
		return false, err
	}
	response, err := p.client.Read()
	if err != nil {
		return false, err
	}

	if response.Type == assuan.Error {
		if response.Code == ErrorCodeNotConfirmed {
			return false, nil
		}
		return false, NewPinentryError(response)
	}

	if response.Type != assuan.Ok {
		return false, fmt.Errorf("unexpected response")
	}

	return true, nil
}
