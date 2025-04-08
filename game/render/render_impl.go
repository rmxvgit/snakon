package render

import (
	"bufio"
	"os"
)

var CLR []byte = []byte{27, 91, 72, 27, 91, 50, 74}

func SetupRender(w, h int) *Renderer {
	renderer := &Renderer{
		Width:        w,
		Height:       h,
		cursor_x:     0,
		cursor_y:     0,
		clr_on_flush: true,
		scr_buf:      make([][]byte, h),
		writer:       *bufio.NewWriter(os.Stdin),
	}

	for i := range renderer.scr_buf {
		renderer.scr_buf[i] = make([]byte, w)
		for ii := range renderer.scr_buf[i] {
			renderer.scr_buf[i][ii] = 32
		}
	}
	return renderer
}

func (renderer *Renderer) Flush() {
	if renderer.clr_on_flush {
		renderer.writer.Write(CLR)
	}
	for i := range renderer.scr_buf {
		renderer.writer.Write(renderer.scr_buf[i])
		renderer.writer.WriteByte('\n')
	}
	renderer.writer.Flush()
}

func (renderer *Renderer) WriteChar(x, y int, c byte) {
	if x < 0 || x >= renderer.Width || y < 0 || y >= renderer.Height {
		return
	}
	renderer.scr_buf[y][x] = c
	renderer.cursor_x = x
	renderer.cursor_y = y
}

func (renderer *Renderer) PlaceCursor(x, y int) {
	if x < 0 || x >= renderer.Width || y < 0 || y >= renderer.Height {
		return
	}
	renderer.cursor_x = x
	renderer.cursor_y = y
}

func (renderer *Renderer) WriteAtCursor(data []byte) {
}

func (renderer *Renderer) CleanBuffer() {
	for i := range renderer.scr_buf {
		for ii := range renderer.scr_buf[i] {
			renderer.scr_buf[i][ii] = 32
		}
	}
}

func (renderer *Renderer) Clear() {
	renderer.writer.Write(CLR)
}
