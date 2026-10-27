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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
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




type ClientCtx struct { // bool persistently kills an upstream session
	StatusChan *chan bool
}




type DownstreamResolve interface {
	Resolve(req RequestedResource) (ResolvableEntry, error)
	StartResource(req ResolvableEntry) (string,error)
}

type KubernetesResolver struct{
	client kubernetes.Interface
	namespace string
	informer cache.SharedIndexInformer
}


//each resolve operation yields this thing, certain info about how a certain resource is to be started

type ResolvableEntry interface {
	WatchStart(req RequestedResource) (ResolvableEntry,error)
	ResourceOK() (bool, error)
	GetResourceReq() RequestedResource
}

type K8sResolvable struct {
	req RequestedResource
	unblockOp ClientCtx
}

func (k *K8sResolvable) GetResourceReq() RequestedResource{
	return k.req
}

func (k *K8sResolvable) WatchStart(req RequestedResource) (ResolvableEntry,error) {

	signalChan := make(chan bool)
	r := &K8sResolvable{req: req, unblockOp: ClientCtx{StatusChan: &signalChan}}

	
	GlobalUnblock.Set(req.EnvName, r.unblockOp)
	return r, nil
}

func (k *K8sResolvable) ResourceOK() (bool, error) {
	
	//has to wait to be signalled by the informer cache to send something on the unblock channel
	
		tmrC := time.NewTimer(15 * time.Second)
		select{
			case <-*k.unblockOp.StatusChan:
				println("unblocked")
				GlobalUnblock.Delete(k.req.EnvName + "-" + k.req.EnvVersion)
				tmrC.Stop()
				return true, nil
	
			case <-tmrC.C:
				println("timeout")
				GlobalUnblock.Delete(k.req.EnvName + "-" + k.req.EnvVersion)
				return false, fmt.Errorf("timeout")
		}

}


func NewKubernetesResolver(client kubernetes.Interface) *KubernetesResolver {
	return &KubernetesResolver{client: client}
}

//returns bool if it can resolve 
func (k *KubernetesResolver) Resolve(req RequestedResource) (ResolvableEntry, error){
	
	

	if k == nil || k.client == nil {
		log.Print("Kubernetes resolver has no client")
		return nil, fmt.Errorf("no client")
	}

	ctx := context.TODO()
	pvc, err := k.client.CoreV1().PersistentVolumeClaims(req.EnvName).Get(ctx, req.EnvName, metav1.GetOptions{})

	if err != nil {
		log.Printf("couldnt find the PV for %s with err %v", req.EnvName, err.Error())
		return nil, err
	}

	if pvc == nil {
		log.Printf("The PVC pointer was null for %s", req.EnvName)
		return nil, fmt.Errorf("pvc nil")
	}

	log.Printf("PV name: %s", pvc.Name)

	podNameStr := req.EnvName + "-" + req.EnvVersion

	svcName := "svc-" + req.EnvName
	imageName := "docker.io/llf2004/" + req.EnvName + ":" + req.EnvVersion

	log.Printf("Pod name: %s", podNameStr)

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podNameStr,
			Namespace: req.EnvName,
			Labels: map[string]string{
				"app": "workspace",
				"env": req.EnvName,
				"version": req.EnvVersion,
				"ssh": "true",
			},
		},

		Spec: corev1.PodSpec{
			Hostname: podNameStr,
			Subdomain: svcName,
			Containers: []corev1.Container{
				{
					Name: "workspace",
					Image: imageName,
					//VolumeMounts: []corev1.VolumeMount{
					//	{
					//		Name: req.EnvName,
					//		MountPath: "/mnt", //stub test, nothing more, remove for the proper link to the backend
					//	},
					//},
					StartupProbe: &corev1.Probe{
						InitialDelaySeconds: 1,
						PeriodSeconds: 1,
						FailureThreshold: 14,
						ProbeHandler: corev1.ProbeHandler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(22),
							},
						},
					},
				},
				

			},
			//Volumes: []corev1.Volume{
			//	{
			//		Name: req.EnvName,
			//		VolumeSource: corev1.VolumeSource{
			//			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			//				ClaimName: req.EnvName, // has to match with the env name
			//		},
			//	},
			//},
		},

	}

	

	log.Printf("Starting resource for %s", req.EnvName)
	
	unblockChan := make(chan bool)
	unblockOp := ClientCtx{
		StatusChan: &unblockChan,
	}

	res := K8sResolvable{}
	res.unblockOp = unblockOp
	res.req = req

	println("adding to demux")
	GlobalUnblock.Set(podNameStr, unblockOp) //allows for the client to watch for unblock operations


	p, err := k.client.CoreV1().Pods(req.EnvName).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {

		if apierrors.IsAlreadyExists(err) {

			GlobalUnblock.Delete(podNameStr) // makes sure not to leak memory
		}else{
			log.Printf("Error creating pod: %v", err)
			return nil, err
		}

	}

	if p.Status.Phase == corev1.PodRunning {
		log.Printf("Pod is running, signal completion downstream, do not register global demux map")
		log.Print("Pod exists, no controller event expected")
		*res.unblockOp.StatusChan <- true //short circuits the waiter function, when it runs it will start right away
	}
	
	return &res, nil //default exist
}

//just starts a kubernetes pod, watching is handled separately
func (k *KubernetesResolver) StartResource(req ResolvableEntry) (string, error) {
	
	svcName := "svc-" + req.GetResourceReq().EnvName
	podNameStr := req.GetResourceReq().EnvName + "-" + req.GetResourceReq().EnvVersion
	namespace := req.GetResourceReq().EnvName

	fqdn := podNameStr + "." + svcName + "." + namespace + ".svc.cluster.local"
	
	//has to wait to be signalled by the informer cache to send something on the unblock channel
	println("fqdn is: ", fqdn)

	state, err := req.ResourceOK()
	if err != nil {
		return "", err
	}
	
	time.Sleep(5 * time.Second)
	if state{
		return fqdn, nil
	}else{
		return "", fmt.Errorf("Resource watch failed")
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
