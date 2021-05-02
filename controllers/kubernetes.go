package controllers

// import(
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// func getKubeInClusterConfig() (*Config, error){
// 	config, err := rest.InclusterConfig()
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
// 	// }
// 	return config, err
// }

// func getKubeOutClusterConfig() (*Config, error){
// 	kubeconfigPath := "~/kube/.config"
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)	
// 	// }
// 	return config, err
// }
