package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	"k8s.io/client-go/rest"
)

var MEMBER_CLUSTER_CLIENT_ID = "<MEMBER_CLUSTER_CLIENT_ID>"

func main() {

	clientID := azidentity.ClientID(MEMBER_CLUSTER_CLIENT_ID)
	opts := &azidentity.ManagedIdentityCredentialOptions{ID: clientID}
	managed, err := azidentity.NewManagedIdentityCredential(opts)
	if err != nil {
		fmt.Printf("error creating the managed identity. err: %v", err.Error())
	}

	token, err := managed.GetToken(context.TODO(), policy.TokenRequestOptions{
		Scopes: []string{"6dae42f8-4368-4678-94ff-3960e28e3630"},
	})
	if err != nil {
		fmt.Printf("error getting the token. err: %v", err.Error())
	}

	cf := rest.Config{
		BearerToken: token.Token,
		Host:        "https://hubcluster-demo-043560-9a369bac.hcp.westus.azmk8s.io:443",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(&cf)
	if err != nil {
		fmt.Printf("error creating clientset. err: %v", err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("error listing pods. err: %v", err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		namespace := "default"
		pod := "demo-msi"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}
		fmt.Print("Nothing so far....")
		time.Sleep(1 * time.Second)
	}
}
