package main

import (
	"k8s.io/client-go/tools/clientcmd"
	"flag"
	clientset "k8s-website-customize-controller/pkg/client/clientset/versioned"
	informers "k8s-website-customize-controller/pkg/client/informers/externalversions"
	"time"
	"k8s-website-customize-controller/pkg/controller"
	"k8s-website-customize-controller/pkg/signals"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	println("main")

	// 处理信号量
	stopCh := signals.SetupSignalHandler()

	// 处理入参
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		panic(err)
	}

	//如果需要处理其它资源，使用该Client
	/*kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}*/

	// 创建websiteClient
	websiteClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	// 创建InformerFactory
	webSiteInformerFactory := informers.NewSharedInformerFactory(websiteClient, time.Second*30) //默认的Resync同步时间  Resync同样会触发update,所以在update方法中判断一下old&new是否一致

	//创建controller
	ctl := controller.NewController()
	//创建handler
	handler := controller.EventHandler{Controller: *ctl}

	//将handler绑定到informer,这样资源的增删改均会调用handler中的OnAdd Ondelete OnUpdate方法
	webSiteInformerFactory.Kevin().V1().Websites().Informer().AddEventHandler(&handler)

	//启动informer
	webSiteInformerFactory.Start(stopCh) //此方法实际上执行的是informer.Informer().Run(stopCh)

	//controller开始处理消息
	go ctl.Run(2, stopCh)

	//通过informer获取liste
	time.Sleep(time.Second * 10)
	lister := webSiteInformerFactory.Kevin().V1().Websites().Lister();
	list, _ := lister.List(labels.NewSelector())
	println("lister")
	for e := range list {
		println(">>>", list[e].ResourceVersion)
	}

	web, _ := lister.Websites("default").Get("website")

	println("web:", web.Kind, web.ResourceVersion)
	<-stopCh

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "./resources/kubeconfig", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
