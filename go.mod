module github.com/rmenn/bootkube

go 1.13

require (
	github.com/coreos/etcd v3.3.12+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kubernetes-incubator/bootkube v0.0.0-00010101000000-000000000000
	github.com/pborman/uuid v0.0.0-20150603214016-ca53cad383ca
	github.com/spf13/cobra v0.0.0-20170515075120-4cdb38c072b8
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	golang.org/x/sys v0.0.0-20190801041406-cbf593c0f2f3
	google.golang.org/grpc v1.19.0
	k8s.io/api v0.0.0-20190202010724-74b699b93c15
	k8s.io/apiextensions-apiserver v0.0.0-20190202013456-d4288ab64945
	k8s.io/apimachinery v0.0.0-20190117220443-572dfc7bdfcb
	k8s.io/client-go v0.0.0-20190202011228-6e4752048fde
	k8s.io/klog v0.2.0
)

replace github.com/kubernetes-incubator/bootkube => github.com/rmenn/bootkube v0.14.1-0.20191021103846-f8cd89c
