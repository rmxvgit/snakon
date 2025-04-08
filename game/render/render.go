package render

import (
	"bufio"
)

type Renderer struct {
	Width        int
	Height       int
	cursor_x     int
	cursor_y     int
	clr_on_flush bool
	scr_buf      [][]byte
	writer       bufio.Writer
}
