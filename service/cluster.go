package service

import (
	clustertype "admin-panel/cluster_type"
	"admin-panel/logger"
	"admin-panel/utils"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var ClusterOnce sync.Once

type LpCluster struct {
	ClientSet *kubernetes.Clientset
}

var LpClusterInstance *LpCluster

const (
	ContainersReady string = "ContainersReady"
	PodInitialized  string = "Initialized"
	PodReady        string = "Ready"
	PodScheduled    string = "PodScheduled"
)

const (
	ConditionTrue    string = "True"
	ConditionFalse   string = "False"
	ConditionUnknown string = "Unknown"
)

func NewLpCluster() *LpCluster {
	ClusterOnce.Do(func() {
		log.Println("init LpCluster🟥")
		lpCluster := &LpCluster{}
		lpCluster.InitClient()
		LpClusterInstance = lpCluster
	})
	return LpClusterInstance
}

func (lpc *LpCluster) InitClient() error {
	var clientset *kubernetes.Clientset
	_, err := os.ReadFile("./env.json")
	clusterType := "in"
	if err == nil {
		clusterType = "out"
	}
	if clusterType == "in" {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println("init cluster failed", "🦠")
			return err
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Println("init cluster failed", "🦠")
			return err
		}
	} else {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return errors.New("cannot init ClientCmd correctly")
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return errors.New("cannot init clientset correctly")
		}
	}
	lpc.ClientSet = clientset
	return nil
}

// list all client
func (lpc *LpCluster) ListClientList(namespace string) (retList []*clustertype.LpClientPodItem, err error) {
	retList = make([]*clustertype.LpClientPodItem, 0)
	useNamespace := apiv1.NamespaceDefault
	if namespace != "" {
		useNamespace = namespace
	}
	deploymentsClient := lpc.ClientSet.AppsV1().Deployments(useNamespace)

	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		errors.WithMessage(err, "K8S list Error:")
		return
	}
	for _, d := range list.Items {
		if !strings.Contains(d.Name, "obridge-chain-client-") {
			continue
		}
		retList = append(retList, &clustertype.LpClientPodItem{
			Name:   d.Name,
			Status: struct{ AvailableReplicas int32 }{AvailableReplicas: d.Status.AvailableReplicas},
		})
	}
	return
}

// list all client
func (lpc *LpCluster) ListClientServiceList(namespace string) (retList []*clustertype.LpClientServiceItem, err error) {
	retList = make([]*clustertype.LpClientServiceItem, 0)
	useNamespace := apiv1.NamespaceDefault
	if namespace != "" {
		useNamespace = namespace
	}
	servicesClient := lpc.ClientSet.CoreV1().Services(useNamespace)

	list, err := servicesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		errors.WithMessage(err, "K8S list Error:")
		return
	}
	for _, d := range list.Items {
		if !strings.Contains(d.Name, "obridge-chain-client-") {
			continue
		}
		ports := make([]struct {
			Protocol string
			Port     int
		}, 0)
		for _, portItemVal := range d.Spec.Ports {
			ports = append(ports, struct {
				Protocol string
				Port     int
			}{Protocol: string(portItemVal.Protocol), Port: int(portItemVal.Port)})
		}
		retList = append(retList, &clustertype.LpClientServiceItem{
			Name:          d.Name,
			ConditionsLen: len(d.Status.Conditions),
			Ports:         ports,
		})
	}
	return
}
func (lpc *LpCluster) RestartPod(namespace string, findName string, name string) (err error) {
	if lpc.ClientSet == nil {
		err = errors.New("cluster link not init correctly")
		return
	}
	log.Println("prepare restart pod, and wait ready", namespace, findName, name)
	if namespace == "" {
		logger.Cluster.Warnf("used default namespace :%s", apiv1.NamespaceDefault)
		namespace = apiv1.NamespaceDefault
	}
	podsClient := lpc.ClientSet.CoreV1().Pods(namespace)
	list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		errors.WithMessage(err, "K8S list Error:")
		return
	}
	if len(list.Items) <= 0 {
		errors.WithMessage(utils.GetNoEmptyError(err), "cannot find pod to restart")
		return
	}
	beDelPodLit := []string{}
	for _, item := range list.Items {
		if strings.Contains(item.Name, findName) && strings.Contains(item.Name, name) {
			beDelPodLit = append(beDelPodLit, item.Name)
		}
	}
	for _, podName := range beDelPodLit {
		delpodErr := podsClient.Delete(context.TODO(), podName, metav1.DeleteOptions{})
		if delpodErr != nil {
			log.Println(fmt.Sprintf("delete pod %s error", podName), delpodErr)
		}
	}
	log.Println("deleting pod is", beDelPodLit)
	err = lpc.WaitDeploymentAvailable(namespace, findName, name, beDelPodLit)
	if err != nil {
		log.Println("service cannot wait ready in limit time", err)
		return
	}
	log.Println("service already ready")
	return
}
func (lpc *LpCluster) WaitDeploymentAvailable(namespace string, findName string, name string, delList []string) error {
	var delSet map[string]bool = make(map[string]bool, 0)
	for _, item := range delList {
		delSet[item] = true
	}
	if lpc.ClientSet == nil {
		err := errors.New("cluster link not init correctly")
		return err
	}
	podsClient := lpc.ClientSet.CoreV1().Pods(namespace)
	retryer := utils.RetryerNew().SetOption(&utils.RepetOption{
		Interval: 1000,
		MaxCount: 60,
	})

	err := retryer.Repet(func() error {
		var ready bool = false
		list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		newPosds := []string{}
		for _, item := range list.Items {
			if strings.Contains(item.Name, findName) && strings.Contains(item.Name, name) {
				_, ok := delSet[item.Name]
				if !ok {
					newPosds = append(newPosds, item.Name)
				}
			}
		}
		for _, podName := range newPosds {
			pod, err := podsClient.Get(context.TODO(), podName, metav1.GetOptions{})
			log.Println(err)
			status := lpc.GetPodStatus(pod)
			if status == "Running" {
				ready = true
			}
			time.Sleep(time.Second * 1)
		}
		if ready {
			return nil
		}
		return errors.New("not ready temporary")
	})
	if err != nil {
		return err
	}
	return nil
}

func (lpc *LpCluster) GetPodStatus(pod *v1.Pod) string {
	for _, cond := range pod.Status.Conditions {
		if string(cond.Type) == ContainersReady {
			if string(cond.Status) != ConditionTrue {
				return "Unavailable"
			}
		} else if string(cond.Type) == PodInitialized && string(cond.Status) != ConditionTrue {
			return "Initializing"
		} else if string(cond.Type) == PodReady {
			if string(cond.Status) != ConditionTrue {
				return "Unavailable"
			}
			for _, containerState := range pod.Status.ContainerStatuses {
				if !containerState.Ready {
					return "Unavailable"
				}
			}
		} else if string(cond.Type) == PodScheduled && string(cond.Status) != ConditionTrue {
			return "Scheduling"
		}
	}
	return string(pod.Status.Phase)
}
