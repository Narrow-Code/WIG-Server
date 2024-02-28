package utils

import (
	"os"
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func GenerateTokenTestPossitive(t  *testing.T) {
	os.Setenv("TOKEN_SECRET", "mysecret")
	token := GenerateToken("testuser", "testhash")

	assert.NotEmpty(t, token, "Generated token should not be empty")
}

func CallerFunctionNameTestPossitive(t *testing.T) {
	// Get the name of the immediate caller function
	callerName := CallerFunctionName(0)
	expected := "utils.CallerFunctionName"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(0) = %s; want %s", callerName, expected)
	}

	// Get the name of the immediate caller function
	callerName = CallerFunctionName(1)
	expected = "utils.TestCallerFunctionName"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(1) = %s; want %s", callerName, expected)
	}

	// Get the name of a function higher in the call stack
	callerName = CallerFunctionName(2)
	expected = "testing.tRunner"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(2) = %s; want %s", callerName, expected)
	}
}

func CallerFunctionNameTestNegativeHighCallback(t *testing.T) {
	callerName := CallerFunctionName(1000)
	expected := "unknown"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(1000) = %s; want %s", callerName, expected)
	}
}
