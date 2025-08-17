package opt

// Custom types for parameters
type Separator string
type Format string

// Boolean flag types with constants
type EqualWidthFlag bool
const (
	EqualWidth   EqualWidthFlag = true
	NoEqualWidth EqualWidthFlag = false
)

// Flags represents the configuration options for the seq command
type Flags struct {
	Separator   Separator      // Output separator (default: newline)
	Format      Format         // Printf-style format string
	EqualWidth  EqualWidthFlag // Pad numbers to equal width with leading zeros
}

// Configure methods for the opt system
func (s Separator) Configure(flags *Flags)     { flags.Separator = s }
func (f Format) Configure(flags *Flags)        { flags.Format = f }
func (e EqualWidthFlag) Configure(flags *Flags) { flags.EqualWidth = e }
