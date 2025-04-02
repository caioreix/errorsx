package errorsx_test

import (
	"testing"

	"github.com/caioreix/errorsx"
	"github.com/stretchr/testify/assert"
)

func TestStackFrame_String(t *testing.T) {
	tt := []struct {
		name     string
		frame    errorsx.StackFrame
		expected string
	}{
		{
			name: "with all fields populated",
			frame: errorsx.StackFrame{
				Function: "main.main",
				File:     "/path/to/main.go",
				Line:     42,
			},
			expected: "main.main\n\t/path/to/main.go:42\n",
		},
		{
			name: "with empty function",
			frame: errorsx.StackFrame{
				Function: "",
				File:     "/path/to/main.go",
				Line:     42,
			},
			expected: "\n\t/path/to/main.go:42\n",
		},
		{
			name: "with empty file",
			frame: errorsx.StackFrame{
				Function: "main.main",
				File:     "",
				Line:     42,
			},
			expected: "main.main\n\t:42\n",
		},
		{
			name: "with zero line",
			frame: errorsx.StackFrame{
				Function: "main.main",
				File:     "/path/to/main.go",
				Line:     0,
			},
			expected: "main.main\n\t/path/to/main.go:0\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.frame.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStack_String(t *testing.T) {
	tt := []struct {
		name     string
		stack    errorsx.Stack
		expected string
	}{
		{
			name:     "empty stack",
			stack:    errorsx.Stack{},
			expected: "",
		},
		{
			name: "single frame",
			stack: errorsx.Stack{
				&errorsx.StackFrame{
					Function: "main.main",
					File:     "/path/to/main.go",
					Line:     42,
				},
			},
			expected: "main.main\n\t/path/to/main.go:42\n",
		},
		{
			name: "multiple frames",
			stack: errorsx.Stack{
				&errorsx.StackFrame{
					Function: "main.main",
					File:     "/path/to/main.go",
					Line:     42,
				},
				&errorsx.StackFrame{
					Function: "main.init",
					File:     "/path/to/main.go",
					Line:     10,
				},
			},
			expected: "main.main\n\t/path/to/main.go:42\nmain.init\n\t/path/to/main.go:10\n",
		},
		{
			name: "with empty frame",
			stack: errorsx.Stack{
				&errorsx.StackFrame{
					Function: "",
					File:     "",
					Line:     0,
				},
			},
			expected: "\n\t:0\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.stack.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}
