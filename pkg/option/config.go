package option

type Config struct {
	Http    *Http    `mapstructure:"http"`
	Cluster *Cluster `mapstructure:"cluster"`
	K8s     *K8s     `mapstructure:"k8s"`
}

type Http struct {
	Addr string `mapstructure:"addr"`
}

type Cluster struct {
	Name string `mapstructure:"cluster_name"`
}

type K8s struct {
	KubeConfPath string `mapstructure:"kube_conf_path"`
}
