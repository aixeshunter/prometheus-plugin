package app

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"github.com/aixeshunter/prometheus-plugin/pkg/constants"
	"github.com/aixeshunter/prometheus-plugin/pkg/k8s"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// NewCmdRun returns "prometheus-plugin run" command.
func NewCmdRun(out io.Writer) *cobra.Command {
	var kubeConfigFile string
	pp := &PrometheusPluginOptions{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run this command in order to set up the application.",
		Run: func(cmd *cobra.Command, args []string) {

			if err := pp.run(out, kubeConfigFile); err != nil {
				glog.Fatalf("Run error: %v", err)
			}
		},
	}

	cmd.PersistentFlags().StringVar(
		&kubeConfigFile, "kubeconfig", "",
		"The KubeConfig file to use when talking to the cluster. If the flag is not set, it will use inClusterConfig",
	)

	cmd.PersistentFlags().StringVar(
		&pp.Namespace, "namespace", constants.DefaultNameSpace, "The namespace that prometheus cluster exists.",
	)

	cmd.PersistentFlags().StringVar(
		&pp.PrometheusName, "prometheus-name", constants.PrometheusName, "The prometheus pod name.",
	)

	cmd.PersistentFlags().StringVar(
		&pp.AlertmanagerName, "alertmanager-name", constants.AlertmanagerName, "The alertmanager pod name.",
	)

	cmd.PersistentFlags().StringVar(
		&pp.Period, "period", constants.DefaultLoopPeriod, "The application worker period.",
	)

	return cmd
}

// PrometheusPluginOptions prometheus-plugin setup options.
type PrometheusPluginOptions struct {
	Namespace        string
	PrometheusName   string
	AlertmanagerName string
	Period           string
}

func (pp *PrometheusPluginOptions) run(out io.Writer, kubeConfigFile string) error {
	client, err := k8s.NewClient(kubeConfigFile)
	if err != nil {
		return fmt.Errorf("create k8s client error, %v", err)
	}

	shutdownC := make(chan struct{})
	go listenToSystemSignal(shutdownC)

	glog.Info("[prometheus-plugin] Prometheus-plugin starting...")
	m := k8s.NewPluginManager(client, pp.Namespace, pp.PrometheusName, pp.AlertmanagerName, pp.Period)
	go m.Run(shutdownC)

	<-shutdownC
	glog.Info("prometheus-plugin exit.")
	return nil
}

// listenToSystemSignal listen system signal and exit exporter.
func listenToSystemSignal(stopC chan<- struct{}) {
	glog.V(5).Info("Listen to system signal.")

	signalChan := make(chan os.Signal, 1)
	ignoreChan := make(chan os.Signal, 1)

	signal.Notify(ignoreChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		glog.V(3).Infof("Shutdown by system signal: %s", sig)
		stopC <- struct{}{}
	}
}
