package render

import (
	"bufio"
	"os"
)

func SetupRender() *Renderer {
	renderer := &Renderer{
		Scr_buffer: bufio.NewWriter(os.Stdin),
	}
	return renderer
}
