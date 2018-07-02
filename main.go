package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientset "github.com/ishustava/k8s-credhub-controller/pkg/client/clientset/versioned"
	informers "github.com/ishustava/k8s-credhub-controller/pkg/client/informers/externalversions"
	"github.com/ishustava/k8s-credhub-controller/pkg/signals"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub"
)

var (
	credhubURL          string
	credhubClient       string
	credhubClientSecret string

	kubeconfig string
	masterURL  string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	credhubControllerClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	credhubClient, err := credhub.New(
		credhubURL,
		credhub.SkipTLSValidation(true),
		credhub.Auth(
			auth.UaaClientCredentials(credhubClient, credhubClientSecret)),
	)
	if err != nil {
		glog.Fatalf("Error creating credhub client: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	credhubSecretInformerFactory := informers.NewSharedInformerFactory(credhubControllerClient, time.Second*30)

	controller := NewController(
		kubeClient,
		credhubControllerClient,
		credhubClient,
		kubeInformerFactory.Core().V1().Secrets(),
		credhubSecretInformerFactory.Pivotal().V1().CredhubSecrets())

	go kubeInformerFactory.Start(stopCh)
	go credhubSecretInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API serverOma. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&credhubURL, "credhuburl", "", "The URL of the CredHub server, e.g. https://my-credhub.com:8844")
	flag.StringVar(&credhubClient, "credhubclient", "", "UAA client name that is authorized to talk to CredHub")
	flag.StringVar(&credhubClientSecret, "credhubclientsecret", "", "Secret of the UAA client that is authorized to talk to CredHub")
}
