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

	re "sshGateway/internalRouter/adressResolver"
)


//context struct for access to the shared informer cache
type KubContext struct {
	ClientSet *kubernetes.Clientset
	Factory informers.SharedInformerFactory
}

func (kc *KubContext)InformerCache(kubeconfig string, url string, unblockChan chan<- bool){
	
	
	config, err := clientcmd.BuildConfigFromFlags(url, kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	kc.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}


	kc.Factory = informers.NewSharedInformerFactory(kc.ClientSet, 15*time.Minute)

	podInformer := kc.Factory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			
		
		AddFunc: func(obj interface{}){
			
		},




		UpdateFunc: func(oldObj, newObj interface{}){ //this thing has to notify the pod on restart
			//oldPod := oldObj.(*corev1.Pod)
        	newPod := newObj.(*corev1.Pod)

			if newPod.Status.Phase == "Running" { //unblocks the downstream through the golang demuxer
				log.Println("Pod started running")


				


				for _, condition := range newPod.Status.Conditions {
        			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
            			
					entry, exists := re.GlobalUnblock.Get(newPod.Name)
					if exists {
					println("exists for pod:", newPod.Name)
					*entry.StatusChan <- true
				}else{
					log.Println("No entry for pod")
				}


        			}
    			}


				


			}

			
		},
		
		
		
		
		DeleteFunc: func(obj interface{}){
			pod := obj.(*corev1.Pod)

			if entry, exists := re.GlobalUnblock.Get(pod.Name); exists {
				*entry.StatusChan <- false
			}
		},
	})

	stopCh := make(chan struct{})
	kc.Factory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) {
		unblockChan <- false
		return
	}
	unblockChan <- true

}