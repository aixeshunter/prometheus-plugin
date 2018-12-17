package app

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCmdVersion provides the version information of ansible-label-plugin.
func NewCmdVersion(out io.Writer) *cobra.Command {
	return nil
}
