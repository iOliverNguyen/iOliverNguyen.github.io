package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall/js"
	"testing"

	"o.o/sample/xtesting"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("wasmExecute", js.FuncOf(wasmExecute))
	<-done
}

func copyStdout(out *string, execFunc func()) error {
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = w

	execFunc() // execute the function to get its output

	resultCh := make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, r)
		resultCh <- b.String()
	}()

	w.Close()
	os.Stdout = old // restoring the real stdout
	*out = <-resultCh
	return nil
}

var jsInput string
var jsInputLength int

func wasmExecute(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return "ERROR: invalid number of arguments"
	}
	var err error
	jsInput = args[0].String()
	jsInputLength, err = strconv.Atoi(args[1].String())
	if err != nil {
		return "ERROR: length is invalid"
	}

	tests := []testing.InternalTest{
		{Name: "TestExample", F: TestExample},
	}
	_ = tests
	testingM := testing.MainStart(xtesting.TestDeps{}, tests, nil, nil, nil)
	testingM.Run()
	logs := ""
	err = copyStdout(&logs, func() {
		testingM.Run()
	})
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return "LOGS: " + logs
}

func TestExample(t *testing.T) {
	expected, actual := jsInputLength, len(jsInput)
	fmt.Printf("TEST: len(%q) == %d\n", jsInput, jsInputLength)

	if len(jsInput) != jsInputLength {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}
