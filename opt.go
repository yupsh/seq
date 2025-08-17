package command

type Separator string
type Format string

type EqualWidthFlag bool

const (
	EqualWidth   EqualWidthFlag = true
	NoEqualWidth EqualWidthFlag = false
)

type flags struct {
	Separator  Separator
	Format     Format
	EqualWidth EqualWidthFlag
}

func (s Separator) Configure(flags *flags)      { flags.Separator = s }
func (f Format) Configure(flags *flags)         { flags.Format = f }
func (e EqualWidthFlag) Configure(flags *flags) { flags.EqualWidth = e }
