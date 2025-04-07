package input

import (
	"os"
	"sync"
)

func (keyboard *KeyboardInput) KeyboardRead() {
	stdin := os.Stdin
	buffer := make([]byte, 4)
	for {
		stdin.Read(buffer)
		keyboard.key_mutex.Lock()
		keyboard.key = []rune(string(buffer))[0]
		keyboard.key_mutex.Unlock()
	}
}

func SetupGameInput() (gm_input *GameInput) {
	gm_input = &GameInput{}

	gm_input.Keyboard = SetupKeyboardInput()

	return
}

func SetupKeyboardInput() (keyboard *KeyboardInput) {
	keyboard = &KeyboardInput{}

	keyboard.key_mutex = sync.Mutex{}
	go keyboard.KeyboardRead()

	return
}

func (keyboard *KeyboardInput) GetLastKey() rune {
	keyboard.key_mutex.Lock()
	key := keyboard.key
	keyboard.key_mutex.Unlock()
	return key
}
