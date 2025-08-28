package main

import (
	"fmt"
	"syscall/js"

	"github.com/hugheaves/scadformat/internal/formatter"
)

func main() {
	fmt.Println("Initializing wasm")
	// Export the function
	js.Global().Set("formatString", formatStringWrapper())
	// Infinite wait so the runtime does not close down
	<-make(chan struct{})
}

func formatStringWrapper() js.Func {
	// Turn call into a promise https://stackoverflow.com/a/67441946
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// Get the input string from the outer function
		input := args[0].String()
		handler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				output, err := formatString(input)
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())

					reject.Invoke(errorObject)
				} else {
					resolve.Invoke(js.ValueOf(output))
				}
			}()

			return nil
		})
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// Format a string using scadformat
func formatString(input string) (string, error) {
	f := formatter.NewFormatter("input.scad")
	output, err := f.FormatBytes([]byte(input))
	return string(output), err
}
