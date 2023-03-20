package xconvey

import (
	"fmt"
	"testing"

	// Workaround for ginkgo flags used in CI test commands. Without this import, ginkgo flags are not registered and
	// the command "go test -v ./... -ginkgo.v" will fail. But for some other reason, ginkgo can not be imported from
	// both test and non-test packages (error: flag redefined: ginkgo.seed). So test files using this package (xconvey)
	// must NOT import ginkgo.
	_ "github.com/onsi/ginkgo"

	"github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/smartystreets/goconvey/convey"
)

// Convey is the method intended for use when declaring the scopes of
// a specification. Each scope has a description and a func() which may contain
// other calls to Convey(), Reset() or Should-style assertions. Convey calls can
// be nested as far as you see fit.
func Convey(items ...any) {
	conveyGomegaSetup(items...)
	convey.Convey(items...)
}

// SConkey is alias to SkipConvey
func SConkey(items ...any) {
	conveyGomegaSetup(items...)
	convey.SkipConvey(items...)
}

// SkipConvey is analogous to Convey except that the scope is not executed
// (which means that child scopes defined within this scope are not run either).
// The reporter will be notified that this step was skipped.
func SkipConvey(items ...any) {
	conveyGomegaSetup(items...)
	convey.SkipConvey(items...)
}

// FConvey is alias to FocusConvey
func FConvey(items ...any) {
	conveyGomegaSetup(items...)
	convey.FocusConvey(items...)
}

// FocusConvey is has the inverse effect of SkipConvey. If the top-level
// Convey is changed to `FocusConvey`, only nested scopes that are defined
// with FocusConvey will be run. The rest will be ignored completely. This
// is handy when debugging a large suite that runs a misbehaving function
// repeatedly as you can disable all but one of that function
// without swaths of `SkipConvey` calls, just a targeted chain of calls
// to FocusConvey.
func FocusConvey(items ...any) {
	conveyGomegaSetup(items...)
	convey.FocusConvey(items...)
}

// The top level Convey() call has the signature
//
//	Convey("message", *testing.T, func() { ... })
//
// This function setup gomega to work with goconvey.
func conveyGomegaSetup(items ...any) {
	if len(items) >= 2 {
		testT, ok := items[1].(*testing.T)
		if ok {
			gomega.Default = gomega.NewWithT(testT).
				ConfigureWithFailHandler(func(message string, callerSkip ...int) {})
		}
	}
}

func Reset() {
	panic("use ConveyReset instead")
}

// ConveyReset registers a cleanup function to be run after each Convey()
// in the same scope. See the examples package for a simple use case.
func ConveyReset(fn func()) {
	convey.Reset(fn)
}

// Ωx is alias to ConveyGomegaExpect
func Ωx(actual any, extra ...any) gomega.Assertion {
	assertion := gomega.Expect(actual, extra...)
	return conveyGomegaAssertion{actual: actual, assertion: assertion}
}

// ConveyGomegaExpect is adapter for using gomega with goconvey. Inside your tests, consider using alias:
//
//	var Ωx = comega.ConveyGomegaExpect
func ConveyGomegaExpect(actual any, extra ...any) gomega.Assertion {
	assertion := gomega.Expect(actual, extra...)
	return conveyGomegaAssertion{actual: actual, assertion: assertion}
}

type conveyGomegaAssertion struct {
	actual    any
	assertion gomega.Assertion
}

func (a conveyGomegaAssertion) Should(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.To(matcher, optionalDescription...)
}

func (a conveyGomegaAssertion) ShouldNot(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.ToNot(matcher, optionalDescription...)
}

func (a conveyGomegaAssertion) To(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	if len(optionalDescription) > 0 {
		panic("optionalDescription is not supported")
	}
	convey.So(a.actual, func(_ any, _ ...any) string {
		success, err := matcher.Match(a.actual)
		if err != nil {
			return fmt.Sprintf("ERROR: %v", err)
		}
		if !success {
			return matcher.FailureMessage(a.actual)
		}
		return ""
	})
	return true
}

func (a conveyGomegaAssertion) ToNot(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	if len(optionalDescription) > 0 {
		panic("optionalDescription is not supported")
	}
	convey.So(a.actual, func(_ any, _ ...any) string {
		success, err := matcher.Match(a.actual)
		if err != nil {
			return fmt.Sprintf("ERROR: %v", err)
		}
		if success {
			return matcher.NegatedFailureMessage(a.actual)
		}
		return ""
	})
	return true
}

func (a conveyGomegaAssertion) NotTo(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.ToNot(matcher, optionalDescription...)
}

func (a conveyGomegaAssertion) WithOffset(offset int) gomegatypes.Assertion {
	return conveyGomegaAssertion{assertion: a.assertion.WithOffset(offset)}
}

func (a conveyGomegaAssertion) Error() gomegatypes.Assertion {
	return conveyGomegaAssertion{assertion: a.assertion.Error()}
}
