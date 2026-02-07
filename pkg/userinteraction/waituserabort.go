package userinteraction

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitUserAbortf(format string, args ...any) {
	WaitUserAbort(fmt.Sprintf(format, args...))
}

// WaitUserAbort displays a message and waits until the user presses CTRL+C (SIGINT).
func WaitUserAbort(msg string) {
	// 1. Show the message from the parameter
	fmt.Println(msg)
	fmt.Println("Press CTRL+C to abort...")

	// 2. Set up a channel to listen for OS signals.
	// We use an unbuffered channel to block until a signal is received.
	sigChan := make(chan os.Signal, 1)

	// 3. Register the channel to receive notifications for specific signals.
	// We are specifically listening for syscall.SIGINT (Interrupt, caused by CTRL+C)
	// and syscall.SIGTERM (Termination, often sent by job control systems).
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 4. Wait until a signal is sent to the channel.
	// The program execution will halt here until the user presses CTRL+C.
	<-sigChan

	// 5. Cleanup the signal notification registration.
	signal.Stop(sigChan)

	// Print a confirmation message
	fmt.Println("\nUser aborted execution by pressing CTRL+C.")

	return
}
