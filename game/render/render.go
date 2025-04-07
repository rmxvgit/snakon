package render

import "bufio"

type Renderer struct {
	Width      int
	Height     int
	Scr_buffer *bufio.Writer
}
