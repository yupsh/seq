package command

import (
	"context"
	"fmt"
	"io"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[string, flags]

func Seq(parameters ...any) yup.Command {
	return command(yup.Initialize[string, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Parse arguments: seq [first [increment]] last
		var first, increment, last float64

		switch len(p.Positional) {
		case 1:
			// seq last
			first = 1
			increment = 1
			fmt.Sscanf(p.Positional[0], "%f", &last)
		case 2:
			// seq first last
			increment = 1
			fmt.Sscanf(p.Positional[0], "%f", &first)
			fmt.Sscanf(p.Positional[1], "%f", &last)
		case 3:
			// seq first increment last
			fmt.Sscanf(p.Positional[0], "%f", &first)
			fmt.Sscanf(p.Positional[1], "%f", &increment)
			fmt.Sscanf(p.Positional[2], "%f", &last)
		default:
			return fmt.Errorf("seq: invalid number of arguments")
		}

		// Determine separator
		separator := "\n"
		if p.Flags.Separator != "" {
			separator = string(p.Flags.Separator)
		}

		// Determine format
		format := "%g"
		if p.Flags.Format != "" {
			format = string(p.Flags.Format)
		}

		// Calculate width for equal-width output
		width := 0
		if bool(p.Flags.EqualWidth) {
			// Simple width calculation
			lastStr := fmt.Sprintf(format, last)
			width = len(lastStr)
		}

		// Generate sequence
		isFirst := true
		for n := first; (increment > 0 && n <= last) || (increment < 0 && n >= last); n += increment {
			if !isFirst {
				_, err := io.WriteString(stdout, separator)
				if err != nil {
					return err
				}
			}
			isFirst = false

			var output string
			if bool(p.Flags.EqualWidth) && width > 0 {
				output = fmt.Sprintf("%0*g", width, n)
			} else {
				output = fmt.Sprintf(format, n)
			}

			_, err := io.WriteString(stdout, output)
			if err != nil {
				return err
			}
		}

		// Add final newline if using newline separator
		if separator == "\n" {
			_, err := io.WriteString(stdout, "\n")
			if err != nil {
				return err
			}
		}

		return nil
	}
}
