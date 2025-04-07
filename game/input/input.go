package input

import "sync"

type GameInput struct {
	Keyboard *KeyboardInput
}

type KeyboardInput struct {
	key_mutex sync.Mutex
	key       rune
}
