package command_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/seq"
)

// Test basic one-argument form: seq LAST
func TestSeq_OneArg(t *testing.T) {
	result := run.Quick(command.Seq(5.0))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3", "4", "5"})
}

// Test two-argument form: seq FIRST LAST
func TestSeq_TwoArgs(t *testing.T) {
	result := run.Quick(command.Seq(2.0, 5.0))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"2", "3", "4", "5"})
}

// Test three-argument form: seq FIRST INCREMENT LAST
func TestSeq_ThreeArgs(t *testing.T) {
	result := run.Quick(command.Seq(1.0, 2.0, 10.0))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "3", "5", "7", "9"})
}

// Test negative increment
func TestSeq_NegativeIncrement(t *testing.T) {
	result := run.Quick(command.Seq(10.0, -2.0, 1.0))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"10", "8", "6", "4", "2"})
}

// Test custom separator
func TestSeq_CustomSeparator(t *testing.T) {
	result := run.Quick(command.Seq(1.0, 5.0, command.Separator(",")))
	assertion.NoError(t, result.Err)
	output := strings.Join(result.Stdout, "")
	assertion.Equal(t, output, "1,2,3,4,5", "comma separated")
}

// Test equal width
func TestSeq_EqualWidth(t *testing.T) {
	result := run.Quick(command.Seq(8.0, 12.0, command.EqualWidth))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"08", "09", "10", "11", "12"})
}

// Test format
func TestSeq_Format(t *testing.T) {
	result := run.Quick(command.Seq(1.0, 3.0, command.Format("%.2f")))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1.00", "2.00", "3.00"})
}

// Test with custom separator and format
func TestSeq_SeparatorAndFormat(t *testing.T) {
	result := run.Quick(command.Seq(1.0, 3.0, command.Separator(" "), command.Format("%.1f")))
	assertion.NoError(t, result.Err)
	output := strings.Join(result.Stdout, "")
	assertion.Equal(t, output, "1.0 2.0 3.0", "space separated with format")
}

// Test equal width with format
func TestSeq_EqualWidthWithFormat(t *testing.T) {
	result := run.Quick(command.Seq(8.0, 12.0, command.EqualWidth, command.Format("%g")))
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
}

// Test seq 0 (from 1 to 0, which is invalid range so outputs nothing with positive increment)
func TestSeq_Zero(t *testing.T) {
	result := run.Quick(command.Seq(0.0))
	assertion.NoError(t, result.Err)
	// When last (0) < first (1) with positive increment, no output
	assertion.Empty(t, result.Stdout)
}

// Test single value
func TestSeq_SingleValue(t *testing.T) {
	result := run.Quick(command.Seq(1.0))
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1"})
}

// Test decimals
func TestSeq_Decimals(t *testing.T) {
	result := run.Quick(command.Seq(1.5, 0.5, 3.5))
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
}

// Test no arguments
func TestSeq_NoArgs(t *testing.T) {
	result := run.Quick(command.Seq())
	assertion.ErrorContains(t, result.Err, "invalid")
}

// Test too many arguments
func TestSeq_TooManyArgs(t *testing.T) {
	result := run.Quick(command.Seq(1.0, 2.0, 3.0, 4.0))
	assertion.ErrorContains(t, result.Err, "invalid")
}

// Test output error
func TestSeq_OutputError(t *testing.T) {
	result := run.Command(command.Seq(1.0, 3.0)).
		WithStdoutError(errors.New("write failed")).
		Run()
	assertion.ErrorContains(t, result.Err, "write failed")
}

// Test output error on separator
func TestSeq_OutputError_Separator(t *testing.T) {
	result := run.Command(command.Seq(1.0, 10.0, command.Separator(","))).
		WithStdoutError(errors.New("write failed")).
		Run()
	assertion.ErrorContains(t, result.Err, "write failed")
}

// Test output error on final newline
func TestSeq_OutputError_FinalNewline(t *testing.T) {
	result := run.Command(command.Seq(1.0, 2.0)).
		WithStdoutError(errors.New("write failed")).
		Run()
	assertion.ErrorContains(t, result.Err, "write failed")
}

// Test various combinations
func TestSeq_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		args     []float64
		expected []string
	}{
		{"1 to 3", []float64{3}, []string{"1", "2", "3"}},
		{"5 to 8", []float64{5, 8}, []string{"5", "6", "7", "8"}},
		{"1 by 3 to 10", []float64{1, 3, 10}, []string{"1", "4", "7", "10"}},
		{"10 by -1 to 8", []float64{10, -1, 8}, []string{"10", "9", "8"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]any, len(tt.args))
			for i, v := range tt.args {
				args[i] = v
			}
			result := run.Quick(command.Seq(args...))
			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

