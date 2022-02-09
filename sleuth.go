package sleuthkube

import (
	"context"
	"fmt"
	"github.com/MatthewDolan/sleuth-client-kube/config"
	appsv1 "k8s.io/api/apps/v1"
	"time"

	"github.com/MatthewDolan/sleuth-client-go"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func Run(stopCh <-chan struct{}, configPathOpt []string) error {
	appConfig, err := config.Load(configPathOpt)
	if err != nil {
		return err
	}

	kubernetesConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	kubernetesClient, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return err
	}

	sleuthClient := sleuth.NewClient(appConfig.Sleuth.OrganizationSlug, appConfig.Sleuth.APIKey)

	factory := informers.NewSharedInformerFactory(kubernetesClient, 15*time.Second)
	factory.Apps().V1().ReplicaSets().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: addFunc(appConfig, sleuthClient),
	})
	factory.Start(stopCh)
	<-stopCh

	return nil
}

func addFunc(config *config.App, client *sleuth.Client) func(obj interface{}) {
	return func(obj interface{}) {
		replicaSet, ok := obj.(*appsv1.ReplicaSet)
		if !ok {
			panic(fmt.Sprintf("unexpected type %T of obj", obj))
		}

		annotations := replicaSet.Annotations

		if err := client.RegisterDeploy(
			context.Background(),
			getDeploymentSlug(&config.Kubernetes.Annotations, annotations),
			config.Kubernetes.Environment,
			getSHA(&config.Kubernetes.Annotations, annotations),
			getDeployedAt(&config.Kubernetes.Annotations, annotations),
		); err != nil {
			panic(err)
		}
	}
}

func getDeploymentSlug(config *config.KubernetesAnnotations, annotations map[string]string) string {
	return getAnnotation(config.DeploymentSlugKey, annotations)
}

func getSHA(config *config.KubernetesAnnotations, annotations map[string]string) string {
	return getAnnotation(config.SHAKey, annotations)
}

func getDeployedAt(config *config.KubernetesAnnotations, annotations map[string]string) time.Time {
	deployedAtStr := getAnnotation(config.DeployedAtKey, annotations)
	if deployedAtStr == "" {
		return time.Now()
	}

	deployedAt, err := time.Parse(time.RFC3339, deployedAtStr)
	if err != nil {
		panic(err)
	}

	return deployedAt
}

func getAnnotation(key string, annotations map[string]string) string {
	if key == "" {
		return ""
	}

	value, found := annotations[key]
	if !found {
		return ""
	}

	return value
}
