package app

import (
	"flag"
	"io"
	"os"

	_ "github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Run creates ans executes new ansible-label-plugin command.
func Run() error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Set("logtostderr", "true")

	// We do not want these flags to show up in --help
	// These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("google-json-key")
	pflag.CommandLine.MarkHidden("log-flush-frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log-backtrace-at")
	pflag.CommandLine.MarkHidden("log-dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")

	c := NewPrometheusPluginCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}

// NewPrometheusPluginCommand return cobra.Command to run prometheus-plugin command.
func NewPrometheusPluginCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "prometheus-plugin",
		Short: "Prometheus plugin on the k8s",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Parse()
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(NewCmdRun(out))

	return cmds
}
