package seq

import (
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/seq/opt"
)

// Flags represents the configuration options for the seq command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Seq creates a new seq command with the given parameters
// Arguments: [first] [increment] last
func Seq(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if len(c.Positional) == 0 {
		fmt.Fprintln(stderr, "seq: missing arguments")
		return fmt.Errorf("missing arguments")
	}

	var first, increment, last float64
	var err error

	switch len(c.Positional) {
	case 1:
		// seq LAST
		first = 1
		increment = 1
		last, err = strconv.ParseFloat(c.Positional[0], 64)
	case 2:
		// seq FIRST LAST
		first, err = strconv.ParseFloat(c.Positional[0], 64)
		if err == nil {
			increment = 1
			last, err = strconv.ParseFloat(c.Positional[1], 64)
		}
	case 3:
		// seq FIRST INCREMENT LAST
		first, err = strconv.ParseFloat(c.Positional[0], 64)
		if err == nil {
			increment, err = strconv.ParseFloat(c.Positional[1], 64)
			if err == nil {
				last, err = strconv.ParseFloat(c.Positional[2], 64)
			}
		}
	default:
		fmt.Fprintln(stderr, "seq: too many arguments")
		return fmt.Errorf("too many arguments")
	}

	if err != nil {
		fmt.Fprintf(stderr, "seq: invalid number: %v\n", err)
		return err
	}

	if increment == 0 {
		fmt.Fprintln(stderr, "seq: increment cannot be zero")
		return fmt.Errorf("increment cannot be zero")
	}

	return c.generateSequence(ctx, first, increment, last, stdout)
}

func (c command) generateSequence(ctx context.Context, first, increment, last float64, output io.Writer) error {
	separator := string(c.Flags.Separator)
	if separator == "" {
		separator = "\n"
	}

	format := string(c.Flags.Format)
	if format == "" {
		format = c.determineFormat(first, increment, last)
	}

	var values []string
	current := first
	counter := 0

	// Handle both ascending and descending sequences
	if increment > 0 {
		for current <= last {
			// Check for cancellation every 1000 iterations to avoid excessive overhead
			if counter%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}

			values = append(values, c.formatNumber(current, format))
			current += increment
			counter++
		}
	} else {
		for current >= last {
			// Check for cancellation every 1000 iterations to avoid excessive overhead
			if counter%1000 == 0 {
				if err := yup.CheckContextCancellation(ctx); err != nil {
					return err
				}
			}

			values = append(values, c.formatNumber(current, format))
			current += increment
			counter++
		}
	}

	// Check for cancellation before final output
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if len(values) > 0 {
		fmt.Fprint(output, strings.Join(values, separator))
		if separator != "\n" {
			fmt.Fprintln(output) // Always end with newline
		}
	}

	return nil
}

func (c command) determineFormat(first, increment, last float64) string {
	if bool(c.Flags.EqualWidth) {
		// Calculate the maximum number of digits needed
		maxVal := math.Max(math.Abs(first), math.Abs(last))
		digits := len(fmt.Sprintf("%.0f", maxVal))
		return fmt.Sprintf("%%0%d.0f", digits)
	}

	// Check if all numbers are integers
	if first == math.Trunc(first) && increment == math.Trunc(increment) && last == math.Trunc(last) {
		return "%.0f"
	}

	return "%g"
}

func (c command) formatNumber(num float64, format string) string {
	return fmt.Sprintf(format, num)
}

func (c command) String() string {
	return fmt.Sprintf("seq %v", c.Positional)
}
