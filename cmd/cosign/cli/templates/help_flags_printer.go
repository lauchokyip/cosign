package templates

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/go-wordwrap"
	flag "github.com/spf13/pflag"
)

const offset = 10

// HelpFlagPrinter is a printer that
// processes the help flag and print
// it to i/o writer
type HelpFlagPrinter struct {
	wrapLimit uint
	out       io.Writer
}

// NewHelpFlagPrinter will initialize a HelpFlagPrinter given the
// i/o writer
func NewHelpFlagPrinter(out io.Writer, wrapLimit uint) *HelpFlagPrinter {
	return &HelpFlagPrinter{
		wrapLimit: wrapLimit,
		out:       out,
	}
}

// PrintHelpFlag will beautify the help flags and print it out to p.out
func (p *HelpFlagPrinter) PrintHelpFlag(flag *flag.Flag) {
	formatBuf := new(bytes.Buffer)
	writeFlag(formatBuf, flag)

	wrappedStr := formatBuf.String()
	flagAndUsage := strings.Split(formatBuf.String(), "\n")
	flagStr := flagAndUsage[0]

	// if the flag usage is longer than one line, wrap it again
	if len(flagAndUsage) > 1 {
		nextLines := strings.Join(flagAndUsage[1:], " ")
		wrappedUsages := wordwrap.WrapString(nextLines, p.wrapLimit-offset)
		wrappedStr = flagStr + "\n" + wrappedUsages
	}
	appendTabStr := strings.ReplaceAll(wrappedStr, "\n", "\n\t")

	fmt.Fprintf(p.out, appendTabStr+"\n\n")
}

// writeFlag will output the help flag based
// on the format provided by getFlagFormat to i/o writer
func writeFlag(out io.Writer, f *flag.Flag) {
	deprecated := ""
	if f.Deprecated != "" {
		deprecated = fmt.Sprintf(" (DEPRECATED: %s)", f.Deprecated)
	}
	fmt.Fprintf(out, getFlagFormat(f), f.Shorthand, f.Name, f.DefValue, f.Usage, deprecated)
}

// getFlagFormat will output the flag format
func getFlagFormat(f *flag.Flag) string {
	var format string
	format = "--%s=%s:\n%s%s"
	if f.Value.Type() == "string" {
		format = "--%s='%s':\n%s%s"
	}

	if len(f.Shorthand) > 0 {
		format = "    -%s, " + format
	} else {
		format = "    %s" + format
	}

	return format
}
