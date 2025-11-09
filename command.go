package command

import (
	"context"
	"fmt"
	"io"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[float64, flags]

func Seq(parameters ...any) gloo.Command {
	return command(gloo.Initialize[float64, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Parse arguments: seq [first [increment]] last
		first, increment, last := 1.0, 1.0, 0.0

		switch len(p.Positional) {
		case 1:
			last = p.Positional[0]
		case 2:
			first = p.Positional[0]
			last = p.Positional[1]
		case 3:
			first = p.Positional[0]
			increment = p.Positional[1]
			last = p.Positional[2]
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
		hasOutput := false
		for n := first; (increment > 0 && n <= last) || (increment < 0 && n >= last); n += increment {
			if !isFirst {
				_, err := io.WriteString(stdout, separator)
				if err != nil {
					return err
				}
			}
			isFirst = false
			hasOutput = true

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

		// Add final newline if using newline separator and we output something
		if separator == "\n" && hasOutput {
			_, err := io.WriteString(stdout, "\n")
			if err != nil {
				return err
			}
		}

		return nil
	}
}
