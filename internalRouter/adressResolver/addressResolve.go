package addressResolver

import (
	//"io"
	"log"
	"sync"
	"time"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

)



type UnblockWatcher struct{
	clientset *kubernetes.Clientset
	mu             sync.RWMutex
	globalDemuxMap map[string]ClientCtx
}

var GlobalUnblock UnblockWatcher

func (u *UnblockWatcher) Set(key string, ctx ClientCtx) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.globalDemuxMap == nil {
		u.globalDemuxMap = make(map[string]ClientCtx)
	}
	u.globalDemuxMap[key] = ctx
}

func (u *UnblockWatcher) Get(key string) (ClientCtx, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	ctx, ok := u.globalDemuxMap[key]
	return ctx, ok
}

func (u *UnblockWatcher) Delete(key string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.globalDemuxMap, key)
}

//the downstream one 
type RequestedResource struct{
	Uname string
	EnvName string
	EnvVersion string
	Resolver DownstreamResolve
}


type ConnEntry struct {
	WantedRes RequestedResource
	backendUp bool
	backendAddr string
}

type ClientCtx struct { // bool persistently kills an upstream session
	StatusChan *chan bool
}







type DownstreamResolve interface {
	Resolve(req RequestedResource) bool
	StartResource(req RequestedResource) error
}

type KubernetesResolver struct{
	client kubernetes.Interface
	namespace string
	informer cache.SharedIndexInformer

	podDefinition corev1.Pod

	UnblockOp ClientCtx 
}

func NewKubernetesResolver(client kubernetes.Interface) *KubernetesResolver {
	return &KubernetesResolver{client: client}
}

//returns bool if it can resolve 
func (k *KubernetesResolver) Resolve(req RequestedResource) bool{
	
	
	//fqdn := req.EnvName + req.EnvVersion + "." + req.Uname + ".svc.cluster.local"

	if k == nil || k.client == nil {
		log.Print("Kubernetes resolver has no client")
		return false
	}

	ctx := context.TODO()
	pvc, err := k.client.CoreV1().PersistentVolumeClaims(req.EnvName).Get(ctx, req.EnvName, metav1.GetOptions{})

	if err != nil {
		log.Print("couldnt find the PV for ", req.EnvName)
		return false
	}

	if pvc == nil {
		log.Printf("The PV pointer was null for %s", req.EnvName)
		return false
	}

	log.Printf("PV name: %s", pvc.Name)

	podNameStr := req.EnvName + "-" + req.EnvVersion

	imageName := "http://endpoint.local:8080/" + req.EnvName + ":" + req.EnvVersion

	log.Printf("Pod name: %s", podNameStr)

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podNameStr,
			Namespace: req.Uname,
			Labels: map[string]string{
				"app": "workspace",
				"env": req.EnvName,
				"version": req.EnvVersion,
				"ssh": "true",
			},
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "workspace",
					Image: imageName,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: req.EnvName,
							MountPath: "/mnt", //stub test, nothing more, remove for the proper link to the backend
						},
					},
				},

			},
			Volumes: []corev1.Volume{
				{
					Name: req.EnvName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: req.EnvName, // has to match with the env name
						},
					},
				},
			},
		},
	}

	k.podDefinition = pod

	unblockChan := make(chan bool)
	k.UnblockOp = ClientCtx{
		StatusChan: &unblockChan,
	}

	println("adding to demux")
	GlobalUnblock.Set(podNameStr, k.UnblockOp) //allows for the client to watch for unblock operations

	return true //default exist
}

//just starts a kubernetes pod, watching is handled separately
func (k *KubernetesResolver) StartResource(req RequestedResource) error {
	log.Printf("Starting resource for %s", req.EnvName)
	
	_, err := k.client.CoreV1().Pods(req.Uname).Create(context.TODO(), &k.podDefinition, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	//has to wait to be signalled by the informer cache to send something on the unblock channel

	
	tmrC := time.NewTimer(15 * time.Second)
	select{
		case <-*k.UnblockOp.StatusChan:
			println("unblocked")

			tmrC.Stop()
			return nil
	
		case <-tmrC.C:
			println("timeout")
			return fmt.Errorf("timeout")
		}
}

type StaticResolver struct{
	allowed map[string]string
}
func (s StaticResolver) Resolve(req RequestedResource) string {
	return req.EnvName + "." + req.Uname
}

func (s *StaticResolver) AddStaticResolveEntry(key, value string) {
	s.allowed[key] = value
}

func (s *StaticResolver) RemoveStaticResolveEntry(key string) {
	delete(s.allowed, key)
}




type dbImpl interface {
	isPresent(key string, envName string) bool
	getRef(key string, envName string) string
}
