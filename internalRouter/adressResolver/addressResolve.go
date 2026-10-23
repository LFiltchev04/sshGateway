package addressResolver

import (
	ir "sshGateway/internalRouter"
	"log"
	"time"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

)

type DownstreamResolve interface {
	Resolve(req ir.RequestedResource) string
	StartResource() error
}

type KubernetesResolver struct{
	client kubernetes.Interface
	namespace string
	informer cache.SharedIndexInformer

	podDefinition corev1.Pod
}

//returns bool if it can resolve 
func (k KubernetesResolver) Resolve(req ir.RequestedResource) bool{
	
	
	//fqdn := req.EnvName + req.EnvVersion + "." + req.Uname + ".svc.cluster.local"

	ctx := context.TODO()
	pv, err := k.client.CoreV1().PersistentVolumes().Get(ctx, req.EnvName, metav1.GetOptions{})

	if err != nil {
		log.Print("couldnt find the PV for ", req.EnvName)
		return false
	}

	if pv == nil {
		log.Printf("The PV pointer was null for %s", req.EnvName)
		return false
	}

	log.Printf("PV name: %s", pv.Name)

	podNameStr := req.EnvName + "-" + req.EnvVersion + "-" + string(time.Now().Unix())

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
							MountPath: "/",
						},
					},
				},

			},
			Volumes: []corev1.Volume{
				{
					Name: req.EnvName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: req.EnvName, // has to match with the username
						},
					},
				},
			},
		},
	}

	k.podDefinition = pod

	return true //default exist
}

//just starts a kubernetes pod, watching is handled separately
func (k KubernetesResolver) StartResource(req *ir.RequestedResource) error {
	log.Printf("Starting resource for %s", req.EnvName)
	
	_, err := k.client.CoreV1().Pods(req.Uname).Create(context.TODO(), &k.podDefinition, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

type StaticResolver struct{
	allowed map[string]string
}
func (s StaticResolver) Resolve(req ir.RequestedResource) string {
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
