package constants

import "time"

const (
	// DefaultKubeConfig default kube config for admin
	DefaultKubeConfig = "/etc/kubernetes/admin.conf"

	// The namespace that prometheus cluster exists
	DefaultNameSpace = "kube-system"

	// The name that prometheus pod be created
	PrometheusName = "prometheus-k8s-0"

	// The name that alertmanager pod be created
	AlertmanagerName = "alertmanager-main-0"

	// DefaultLoopPeriod defines how long prometheus plugin detects `unknown` status pods
	DefaultLoopPeriod = "2m"

	NodeLost = "NodeLost"

	// APICallRetryInterval defines how long k8s should wait before retrying a failed API operation
	APICallRetryInterval = 500 * time.Millisecond

	// DeletePodTimeout specifies how long k8s should wait for applying the label and taint on the master before timing out
	DeletePodTimeout = 10 * time.Second

	// UpdatePodTimeout specifies how long k8s should wait for applying the label and taint on the master before timing out
	UpdatePodTimeout = 10 * time.Second
)
