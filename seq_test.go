package seq

import (
	"context"
	"strings"
	"testing"
	"time"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/seq/opt"
)

func TestSeqBasic(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "single argument - count from 1 to 5",
			args:     []any{"5"},
			expected: "1\n2\n3\n4\n5\n",
		},
		{
			name:     "two arguments - range from 3 to 7",
			args:     []any{"3", "7"},
			expected: "3\n4\n5\n6\n7\n",
		},
		{
			name:     "three arguments - range with increment",
			args:     []any{"1", "2", "9"},
			expected: "1\n3\n5\n7\n9\n",
		},
		{
			name:     "negative increment - descending",
			args:     []any{"10", "-2", "2"},
			expected: "10\n8\n6\n4\n2\n",
		},
		{
			name:     "decimal numbers",
			args:     []any{"1.5", "0.5", "3"},
			expected: "1.5\n2\n2.5\n3\n",
		},
		{
			name:     "negative start",
			args:     []any{"-3", "1"},
			expected: "-3\n-2\n-1\n0\n1\n",
		},
		{
			name:     "empty sequence - start > end with positive increment",
			args:     []any{"5", "3"},
			expected: "",
		},
		{
			name:     "empty sequence - start < end with negative increment",
			args:     []any{"3", "-1", "5"},
			expected: "",
		},
		{
			name:     "single value - exact match",
			args:     []any{"5", "5"},
			expected: "5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(tt.args...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSeqSeparator(t *testing.T) {
	tests := []struct {
		name      string
		args      []any
		separator opt.Separator
		expected  string
	}{
		{
			name:      "comma separator",
			args:      []any{"1", "3"},
			separator: opt.Separator(","),
			expected:  "1,2,3\n",
		},
		{
			name:      "space separator",
			args:      []any{"1", "3"},
			separator: opt.Separator(" "),
			expected:  "1 2 3\n",
		},
		{
			name:      "tab separator",
			args:      []any{"1", "3"},
			separator: opt.Separator("\t"),
			expected:  "1\t2\t3\n",
		},
		{
			name:      "custom separator",
			args:      []any{"1", "3"},
			separator: opt.Separator(" | "),
			expected:  "1 | 2 | 3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(append(tt.args, tt.separator)...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSeqFormat(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		format   opt.Format
		expected string
	}{
		{
			name:     "integer format",
			args:     []any{"1.5", "0.5", "3"},
			format:   opt.Format("%.0f"),
			expected: "2\n2\n2\n3\n",
		},
		{
			name:     "two decimal places",
			args:     []any{"1", "0.1", "1.3"},
			format:   opt.Format("%.2f"),
			expected: "1.00\n1.10\n1.20\n1.30\n",
		},
		{
			name:     "scientific notation",
			args:     []any{"1000", "1000", "3000"},
			format:   opt.Format("%.2e"),
			expected: "1.00e+03\n2.00e+03\n3.00e+03\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(append(tt.args, tt.format)...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSeqEqualWidth(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "equal width with different digit counts",
			args:     []any{"8", "12", opt.EqualWidth},
			expected: "08\n09\n10\n11\n12\n",
		},
		{
			name:     "equal width with negative numbers",
			args:     []any{"98", "102", opt.EqualWidth},
			expected: "098\n099\n100\n101\n102\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(tt.args...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSeqErrors(t *testing.T) {
	tests := []struct {
		name          string
		args          []any
		expectError   bool
		errorContains string
	}{
		{
			name:          "no arguments",
			args:          []any{},
			expectError:   true,
			errorContains: "missing arguments",
		},
		{
			name:          "too many arguments",
			args:          []any{"1", "2", "3", "4"},
			expectError:   true,
			errorContains: "too many arguments",
		},
		{
			name:          "invalid first number",
			args:          []any{"abc", "5"},
			expectError:   true,
			errorContains: "invalid number",
		},
		{
			name:          "invalid second number",
			args:          []any{"1", "abc"},
			expectError:   true,
			errorContains: "invalid number",
		},
		{
			name:          "invalid increment",
			args:          []any{"1", "abc", "5"},
			expectError:   true,
			errorContains: "invalid number",
		},
		{
			name:          "zero increment",
			args:          []any{"1", "0", "5"},
			expectError:   true,
			errorContains: "increment cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(tt.args...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) && !strings.Contains(stderr.String(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got: %v\nStderr: %s", tt.errorContains, err, stderr.String())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v\nStderr: %s", err, stderr.String())
				}
			}
		})
	}
}

func TestSeqContextCancellation(t *testing.T) {
	// Create a command that would generate a very large sequence
	cmd := Seq("1", "1", "100000")

	// Create a context that will be cancelled quickly
	ctx, cancel := context.WithCancel(context.Background())

	var output strings.Builder
	var stderr strings.Builder

	// Cancel context immediately
	cancel()

	err := cmd.Execute(ctx, nil, &output, &stderr)

	// Should detect cancellation and return error
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "context cancelled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestSeqContextCancellationTimeout(t *testing.T) {
	// Create a command that would generate a very large sequence
	cmd := Seq("1", "1", "100000")

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	var output strings.Builder
	var stderr strings.Builder

	err := cmd.Execute(ctx, nil, &output, &stderr)

	// Should timeout and return error
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestSeqString(t *testing.T) {
	cmd := Seq("1", "2", "10")
	result := cmd.String()
	expected := "seq [1 2 10]"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSeqInterface(t *testing.T) {
	// Verify that Seq command implements yup.Command interface
	var _ yup.Command = Seq("1", "5")
}

func TestSeqCombinedFlags(t *testing.T) {
	// Test combining multiple flags
	cmd := Seq("1", "3", opt.Separator(","), opt.EqualWidth)

	var output strings.Builder
	var stderr strings.Builder

	ctx := context.Background()
	err := cmd.Execute(ctx, nil, &output, &stderr)

	if err != nil {
		t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
	}

	result := output.String()
	expected := "1,2,3\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSeqEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{
			name:     "very small increment",
			args:     []any{"0", "0.1", "0.3"},
			expected: "0\n0.1\n0.2\n0.3\n",
		},
		{
			name:     "large numbers",
			args:     []any{"1000000", "1000002"},
			expected: "1000000\n1000001\n1000002\n",
		},
		{
			name:     "negative range descending",
			args:     []any{"-1", "-1", "-5"},
			expected: "-1\n-2\n-3\n-4\n-5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Seq(tt.args...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, nil, &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func BenchmarkSeqSmall(b *testing.B) {
	cmd := Seq("1", "10")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, nil, &output, &stderr)
	}
}

func BenchmarkSeqLarge(b *testing.B) {
	cmd := Seq("1", "1000")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, nil, &output, &stderr)
	}
}

func BenchmarkSeqWithFormatting(b *testing.B) {
	cmd := Seq("1", "100", opt.Format("%.2f"), opt.Separator(","))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, nil, &output, &stderr)
	}
}

// Example tests for documentation
func ExampleSeq() {
	cmd := Seq("1", "3")
	ctx := context.Background()

	var output strings.Builder
	cmd.Execute(ctx, nil, &output, &strings.Builder{})
	// Output would be: 1\n2\n3\n
}

func ExampleSeq_withIncrement() {
	cmd := Seq("1", "2", "7")
	ctx := context.Background()

	var output strings.Builder
	cmd.Execute(ctx, nil, &output, &strings.Builder{})
	// Output would be: 1\n3\n5\n7\n
}

func ExampleSeq_withSeparator() {
	cmd := Seq("1", "3", opt.Separator(","))
	ctx := context.Background()

	var output strings.Builder
	cmd.Execute(ctx, nil, &output, &strings.Builder{})
	// Output would be: 1,2,3\n
}
