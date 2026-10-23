
package internalrouter

import (
	"k8s.io/client-go/kubernetes"
)



type UnblockWatcher struct{
	clientset *kubernetes.Clientset
	GlobalDemuxMap map[string]ClientCtx
}

var GlobalUnblock UnblockWatcher