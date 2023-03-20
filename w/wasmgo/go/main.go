package main

import (
	"github.com/google/uuid"
	"strconv"
	"strings"
	"syscall/js"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("wasmGenerateUUIDs", js.FuncOf(generateUUIDs))
	<-done
}

func generateUUIDs(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return uuid.NewString()
	}
	count, err := strconv.Atoi(args[0].String())
	if err != nil {
		return "ERROR: invalid number"
	}
	if count > 10 {
		return "ERROR: too many"
	}
	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		uuids[i] = uuid.NewString()
	}
	return strings.Join(uuids, "\n")
}
