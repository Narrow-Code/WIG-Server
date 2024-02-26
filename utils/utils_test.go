package utils

import (
	"os"
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t  *testing.T) {
	// Set up test environment
	os.Setenv("TOKEN_SECRET", "mysecret")

	// Test Case
	token := GenerateToken("testuser", "testhash")

	// Assertion
	assert.NotEmpty(t, token, "Generated token should not be empty")
}

func TestCallerFunctionName(t *testing.T) {
	// Test case 1: Test getting the name of the immediate caller function
	callerName := CallerFunctionName(1)
	expected := "utils.TestCallerFunctionName"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(1) = %s; want %s", callerName, expected)
	}

	// Test case 2: Test getting the name of a function higher in the call stack
	callerName = CallerFunctionName(2)
	expected = "testing.tRunner"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(2) = %s; want %s", callerName, expected)
	}
}

func TestCallerFunctionNameNegative(t *testing.T) {
	callerName := CallerFunctionName(1000)
	expected := "unknown"
	if !strings.Contains(callerName, expected) {
		t.Errorf("CallerFunctionName(-1) = %s; want %s", callerName, expected)
	}
}
