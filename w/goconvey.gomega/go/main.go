package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	// "syscall/js"
	"testing"

	"github.com/onsi/gomega"

	"o.o/sample/xconvey"
	"o.o/sample/xtesting"
)

func main() {
	wasmExecute()
}

func copyStdout(out *string, execFunc func()) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
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
}

func wasmExecute() {
	tests := []testing.InternalTest{
		{Name: "GoConvey", F: TestGoconvey},
	}
	testingM := testing.MainStart(xtesting.TestDeps{}, tests, nil, nil, nil)

	logs, code := "", 0
	copyStdout(&logs, func() {
		code = testingM.Run()
	})
	fmt.Println("---\n", code, logs)
	fmt.Println("---\n")
}

func TestGoconvey(t *testing.T) {
	Ω := xconvey.ConveyGomegaExpect //nolint:govet

	xconvey.Convey("Convey", t, func() {
		s := "Start"
		defer func() { fmt.Println(s) }()

		xconvey.FConvey("Executing test 1", func() {
			s += " -> test1"
			Ω(s).To(gomega.Equal("Start -> test1"))

			xconvey.Convey("Test 1.1:", func() {
				s += " -> test1.1"
				Ω(s).To(gomega.Equal("Start -> test1 -> test1.1"))
			})
			xconvey.Convey("Test 1.2:", func() {
				s += " -> test1.2"
				Ω(s).To(gomega.Equal("Start -> test1 -> test1.2"))
			})
		})
		xconvey.Convey("Execute test 2", func() {
			s += " -> test2"
			Ω(s).To(gomega.Equal("Start -> test2"))

			xconvey.Convey("Test 2.1:", func() {
				s += " -> test2.1"
			})
			xconvey.Convey("Test 2.2:", func() {
				s += " -> test2.2"
				Ω(s).Should(gomega.Equal("Start -> test2 -> test2.2"))

				// this line will be fail for seeing the failure message
				Ω(s).Should(gomega.Equal("Start -> TEST FAIL"))
			})
		})
	})
}
