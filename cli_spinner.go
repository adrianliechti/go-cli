package cli

import (
	"fmt"
	"sync"
	"time"
)

// Spinner frames (braille dots)
var spinnerFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

func Run(title string, fn func() error) error {
	var fnErr error
	var wg sync.WaitGroup
	done := make(chan struct{})

	// Start the action in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fnErr = fn()
		close(done)
	}()

	// Hide cursor
	hideCursor()
	defer showCursor()

	// Spinner loop
	frame := 0
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			clearLine()
			fmt.Println(themeSuccess("✓") + " " + themeText(title))
			wg.Wait()
			return fnErr
		case <-ticker.C:
			clearLine()
			spinner := themeHighlight(string(spinnerFrames[frame]))
			fmt.Print(spinner + " " + themeText(title))
			frame = (frame + 1) % len(spinnerFrames)
		}
	}
}

func MustRun(title string, fn func() error) error {
	err := Run(title, fn)

	if err != nil {
		Fatal(err)
	}

	return err
}
