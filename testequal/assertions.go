//go:build !solution

package testequal

import (
	"fmt"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if areEqual(expected, actual) {
		return true
	}

	formatErrorMessage(t, "not equal:", expected, actual, msgAndArgs...)
	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if !areEqual(expected, actual) {
		return true
	}

	formatErrorMessage(t, "should not be equal:", expected, actual, msgAndArgs...)
	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if !AssertEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if !AssertNotEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// areEqual compares expected and actual values based on their types
func areEqual(expected, actual interface{}) bool {
	// Handle nil
	if expected == nil || actual == nil {
		return expected == actual
	}

	// Compare based on type
	switch e := expected.(type) {
	case int:
		if a, ok := actual.(int); ok {
			return e == a
		}
	case int8:
		if a, ok := actual.(int8); ok {
			return e == a
		}
	case int16:
		if a, ok := actual.(int16); ok {
			return e == a
		}
	case int32:
		if a, ok := actual.(int32); ok {
			return e == a
		}
	case int64:
		if a, ok := actual.(int64); ok {
			return e == a
		}
	case uint:
		if a, ok := actual.(uint); ok {
			return e == a
		}
	case uint8:
		if a, ok := actual.(uint8); ok {
			return e == a
		}
	case uint16:
		if a, ok := actual.(uint16); ok {
			return e == a
		}
	case uint32:
		if a, ok := actual.(uint32); ok {
			return e == a
		}
	case uint64:
		if a, ok := actual.(uint64); ok {
			return e == a
		}
	case string:
		if a, ok := actual.(string); ok {
			return e == a
		}
	case []int:
		if a, ok := actual.([]int); ok {
			return compareIntSlices(e, a)
		}
	case []byte:
		if a, ok := actual.([]byte); ok {
			return compareByteSlices(e, a)
		}
	case map[string]string:
		if a, ok := actual.(map[string]string); ok {
			return compareStringMaps(e, a)
		}
	}

	return false
}

func compareIntSlices(a, b []int) bool {
	// Different nil status means not equal
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareByteSlices(a, b []byte) bool {
	// Different nil status means not equal
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareStringMaps(a, b map[string]string) bool {
	// Different nil status means not equal
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

// formatErrorMessage formats and prints the error message
func formatErrorMessage(t T, prefix string, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	msg := formatMessage(msgAndArgs...)

	if msg != "" {
		t.Errorf("%s\n\texpected: %#v\n\tactual  : %#v\n\tmessage : %s",
			prefix, expected, actual, msg)
	} else {
		t.Errorf("%s\n\texpected: %#v\n\tactual  : %#v",
			prefix, expected, actual)
	}
}

// formatMessage formats the message using fmt.Sprintf if arguments are provided
func formatMessage(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}

	msg := msgAndArgs[0]
	args := msgAndArgs[1:]

	if len(args) == 0 {
		if str, ok := msg.(string); ok {
			return str
		}
		return fmt.Sprint(msg)
	}

	return fmt.Sprintf(msg.(string), args...)
}
