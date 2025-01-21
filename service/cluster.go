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
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
		log.Println("clusterType", clusterType)
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", "/.kube/config", "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "/.kube/config", "absolute path to the kubeconfig file")
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
func (lpc *LpCluster) DescPodEnv(namespace string, podName string) (envList []struct {
	Name  string
	Value string
}, err error) {
	log.Printf("Starting to get pod env vars - namespace: %s, podName: %s", namespace, podName)

	useNamespace := apiv1.NamespaceDefault
	if namespace != "" {
		useNamespace = namespace
	}
	log.Printf("Using namespace: %s", useNamespace)

	deploymentsClient := lpc.ClientSet.AppsV1().Deployments(useNamespace)
	deployment, err := deploymentsClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting deployment: %v", err)
		return nil, err
	}

	log.Printf("Successfully got deployment: %s", deployment.Name)

	// Use map to handle duplicates
	envMap := make(map[string]string)

	containerCount := len(deployment.Spec.Template.Spec.Containers)
	log.Printf("Found %d containers in deployment", containerCount)

	if containerCount <= 0 {
		log.Printf("No containers found in deployment %s/%s", namespace, podName)
		return nil, fmt.Errorf("no containers found in deployment")
	}

	// Loop through all containers
	for i, container := range deployment.Spec.Template.Spec.Containers {
		log.Printf("Processing container %d: %s", i+1, container.Name)

		envCount := len(container.Env)
		log.Printf("Container %s has %d environment variables", container.Name, envCount)

		// Later containers will override earlier ones for duplicate keys
		for _, envItem := range container.Env {
			oldValue, exists := envMap[envItem.Name]
			if exists {
				log.Printf("Overriding env var %s: old value=%s, new value=%s",
					envItem.Name, oldValue, envItem.Value)
			}
			envMap[envItem.Name] = envItem.Value
		}
	}

	log.Printf("Total unique environment variables found: %d", len(envMap))

	// Convert map back to slice
	envList = make([]struct {
		Name  string
		Value string
	}, 0, len(envMap))

	for name, value := range envMap {
		envList = append(envList, struct {
			Name  string
			Value string
		}{
			Name:  name,
			Value: value,
		})
		log.Printf("Added env var to final list - Name: %s, Value: %s", name, value)
	}

	log.Printf("Successfully processed all environment variables")
	return envList, nil
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
		if !strings.Contains(d.Name, "chain-client-") {
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
		if !strings.Contains(d.Name, "chain-client-") {
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

func (lpc *LpCluster) PathEnv(namespace string, deploymentName string, containersName string, key string, value string) (ret bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deploy, getErr := lpc.ClientSet.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if getErr != nil {
		err = getErr
		return
	}
	if deploy.Name != "" {
		log.Println(deploy.Name)
		ctxApply, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer ctxCancel()
		patch := []byte(fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","env":[{"name":"%s","value":"%s"}]}]}}}}`, containersName, key, value))
		_, pathErr := lpc.ClientSet.AppsV1().Deployments(namespace).Patch(ctxApply, deploymentName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		log.Println(pathErr)
		if pathErr != nil {
			err = pathErr
			return
		}
		ret = true
		return
	}

	ret = false
	err = errors.WithMessage(utils.GetNoEmptyError(err), "PathEnv err")
	return
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
