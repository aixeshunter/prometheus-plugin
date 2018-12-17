package k8s

import (
	"fmt"
	"time"

	"github.com/aixeshunter/prometheus-plugin/pkg/constants"
	"github.com/aixeshunter/prometheus-plugin/pkg/utils"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// PluginManager
type PluginManager struct {
	Namespace string
	client    kubernetes.Interface
	Period    time.Duration
	PodList   []string
}

// NewPluginManager
func NewPluginManager(
	client kubernetes.Interface,
	namespace string,
	prometheusName string,
	alertmanagerName string,
	period string) *PluginManager {

	var t time.Duration
	p, err := utils.GetTimeDurationStringToSeconds(period)
	if err != nil {
		t = 1 * time.Minute
	}
	t = time.Duration(p) * time.Second

	return &PluginManager{
		Namespace: namespace,
		client:    client,
		Period:    t,
		PodList:   []string{prometheusName, alertmanagerName},
	}
}

// Run loops implements delete unavailable pods of prometheus cluster.
func (m *PluginManager) Run(stopCh chan struct{}) error {
	go wait.Until(m.PrometheusDetect, m.Period, stopCh)
	<-stopCh
	glog.Info("label manager cron job exit.")
	return nil
}

func (m *PluginManager) PrometheusDetect() {
	for _, p := range m.PodList {
		pod, err := GetPodbyName(m.client, metav1.GetOptions{}, p, m.Namespace)
		if err != nil {
			glog.Errorf("get the pod %s failed with %v", p, err)
			continue
		}

		done := make(chan bool)
		var gracePeriodSeconds int64 = 0
		if isPodUnknown(pod) {
			go DeletePod(m.client, &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds}, p, m.Namespace, done)

			select {
			case <-done:
				glog.Infof("Delete the pod %s return.", p)
			case <-time.After(constants.DeletePodTimeout):
				glog.Infof("Need to patch the pod %s with matadata finalizers null.", p)
				PatchPod(m.client, p, m.Namespace)
			}
		}
	}
}

// Judge the pod in unknown status
func isPodUnknown(pod *corev1.Pod) bool {
	if pod.Status.Reason == constants.NodeLost || corev1.PodUnknown == pod.Status.Phase {
		return true
	}

	return false
}

// Delete the pod with name and namespace
func DeletePod(client kubernetes.Interface, opts *metav1.DeleteOptions, name string, namespace string, done chan bool) {
	if err := client.CoreV1().Pods(namespace).Delete(name, opts); err != nil {
		glog.Errorf("force delete the pod %s failed with %v", name, err)
	}

	done <- true
}

// Update the pod with name and namespace
func PatchPod(client kubernetes.Interface, name string, namespace string) error {
	addControllerPatch := fmt.Sprint(`{"metadata":{"finalizers":null}}`)
	return wait.Poll(constants.APICallRetryInterval, constants.UpdatePodTimeout, func() (bool, error) {
		if _, err := client.CoreV1().Pods(namespace).Patch(name, types.StrategicMergePatchType, []byte(addControllerPatch)); err != nil {
			glog.Errorf("update the pod %s failed with %v", name, err)
			return false, nil
		}

		return true, nil
	})
}

// Get the pod by name
func GetPodbyName(client kubernetes.Interface, opts metav1.GetOptions, name string, namespace string) (*corev1.Pod, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(name, opts)
	if err != nil {
		return nil, fmt.Errorf("get the pod %s error, %v", name, err)
	}

	return pod, nil
}
