package routes

import (
	"fmt"
	"os"
	"path"

	"github.com/altid/server/internal/message"
)

type InputHandler struct{}

func NewInput() *InputHandler { return &InputHandler{} }

func (*InputHandler) Normal(msg *message.Message) (interface{}, error) {
	fp := path.Join(msg.Service, msg.Buffer, "input")
	i := &input{
		path: fp,
	}

	return i, nil
}

func (*InputHandler) Stat(msg *message.Message) (os.FileInfo, error) {
	return os.Stat(path.Join(msg.Service, msg.Buffer, "input"))
}

type input struct {
	path string
}

// Simple wrapper around an open call
func (i *input) ReadAt(b []byte, off int64) (n int, err error) {
	fp, err := os.OpenFile(i.path, os.O_RDONLY, 0600)
	if err != nil {
		return
	}

	defer fp.Close()
	return fp.ReadAt(b, off)
}

// Open in correct modes
func (i *input) WriteAt(p []byte, off int64) (n int, err error) {
	fp, err := os.OpenFile(i.path, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return
	}

	defer fp.Close()
	return fmt.Fprintf(fp, "%s\n", p)
}

func (i *input) Close() error { return nil }
