package kubernetesImpl


import (
	"time"
	"log"
	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	r "sshGateway/internalRouter"
)





func  InformerCache(){

	
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 30*time.Minute)


	podInformer := factory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		
		
		
		
		AddFunc: func(obj interface{}){
			
		},




		UpdateFunc: func(oldObj, newObj interface{}){ //this thing has to notify the pod on restart
			oldPod := oldObj.(*corev1.Pod)
        	newPod := newObj.(*corev1.Pod)

			if oldPod.Status.Phase == "Pending" && newPod.Status.Phase == "Running" { //unblocks the downstream through the golang demuxer
				log.Println("Pod started running")


				
				entry, exists := r.GlobalUnblock.GlobalDemuxMap[newPod.Name]
				if exists {
					*entry.StatusChan <- true
				}else{
					log.Println("No entry for pod")
					*entry.StatusChan <- false
				}


			}

			if oldPod.Status.Phase == newPod.Status.Phase {
				log.Println("Pod status unchanged, is an error")
			}
		},
		
		
		
		
		DeleteFunc: func(obj interface{}){
			pod := obj.(*corev1.Pod)

			if _, exists := r.GlobalUnblock.GlobalDemuxMap[pod.Name]; exists {
				*r.GlobalUnblock.GlobalDemuxMap[pod.Name].StatusChan <- false
			}
		},
	})

}