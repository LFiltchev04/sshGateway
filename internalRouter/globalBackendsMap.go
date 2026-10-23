
package internalrouter

import (
	"k8s.io/client-go/kubernetes"
)

type ClientCtx struct { // bool persistently kills an upstream session
	StatusChan *chan bool
}

type UnblockWatcher struct{
	clientset *kubernetes.Clientset
	GlobalDemuxMap map[string]ClientCtx
}

var GlobalUnblock UnblockWatcher