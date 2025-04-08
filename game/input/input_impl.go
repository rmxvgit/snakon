package input

import (
	"os"
	"os/exec"
	"snakon/utils"
	"sync"
)

func (keyboard *KeyboardInput) KeyboardRead() {
	stdin := os.Stdin
	buffer := make([]byte, 1)
	for {
		stdin.Read(buffer)
		keyboard.key_mutex.Lock()
		keyboard.key = buffer[0]
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

	cmd := exec.Command("stty", "-F", "/dev/tty", "-icanon", "min", "1")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	utils.PanicOnError(err)

	keyboard.key_mutex = sync.Mutex{}
	go keyboard.KeyboardRead()

	return
}

func (keyboard *KeyboardInput) GetLastKey() byte {
	keyboard.key_mutex.Lock()
	key := keyboard.key
	keyboard.key_mutex.Unlock()
	return key
}

func (keyboard *KeyboardInput) ConsumeLastKey() byte {
	last_key := keyboard.key
	keyboard.key_mutex.Lock()
	keyboard.key = 0
	keyboard.key_mutex.Unlock()
	return last_key
}
