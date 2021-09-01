package v1alpha1_test

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
)

func TestClusterValidateUpdateKubernetesVersionImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			KubernetesVersion:         v1alpha1.Kube119,
			ExternalEtcdConfiguration: &v1alpha1.ExternalEtcdConfiguration{Count: 3},
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count: 3, Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.KubernetesVersion = v1alpha1.Kube120

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationEqual(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count:           3,
				Endpoint:        &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
				MachineGroupRef: &v1alpha1.Ref{Name: "test", Kind: "MachineConfig"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		Count:           3,
		Endpoint:        &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
		MachineGroupRef: &v1alpha1.Ref{Name: "test", Kind: "MachineConfig"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count:           3,
				Endpoint:        &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
				MachineGroupRef: &v1alpha1.Ref{Name: "test", Kind: "MachineConfig"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		Count:           10,
		Endpoint:        &v1alpha1.Endpoint{Host: "1.1.1.1/2"},
		MachineGroupRef: &v1alpha1.Ref{Name: "test2", Kind: "SecondMachineConfig"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationOldEndpointImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/2"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationOldEndpointNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Endpoint: nil,
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationNewEndpointNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		Endpoint: nil,
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationOldMachineGroupRefImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				MachineGroupRef: &v1alpha1.Ref{Name: "test1", Kind: "MachineConfig"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		MachineGroupRef: &v1alpha1.Ref{Name: "test2", Kind: "MachineConfig"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationOldMachineGroupRefNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				MachineGroupRef: nil,
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		MachineGroupRef: &v1alpha1.Ref{Name: "test", Kind: "MachineConfig"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateControlPlaneConfigurationNewMachineGroupRefNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				MachineGroupRef: &v1alpha1.Ref{Name: "test", Kind: "MachineConfig"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ControlPlaneConfiguration = v1alpha1.ControlPlaneConfiguration{
		MachineGroupRef: nil,
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateDatacenterRefImmutableEqual(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "test", Kind: "DatacenterConfig",
			},
		},
	}
	c := cOld.DeepCopy()

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateDatacenterRefImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "test", Kind: "DatacenterConfig",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.DatacenterRef = v1alpha1.Ref{Name: "test2", Kind: "SecondDatacenterConfig"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateDatacenterRefImmutableName(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "test", Kind: "DatacenterConfig",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.DatacenterRef = v1alpha1.Ref{Name: "test2", Kind: "DatacenterConfig"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateDatacenterRefNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "test", Kind: "DatacenterConfig",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.DatacenterRef = v1alpha1.Ref{}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateExternalEtcdReplicasImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ExternalEtcdConfiguration: &v1alpha1.ExternalEtcdConfiguration{Count: 3},
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count: 3, Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ExternalEtcdConfiguration.Count = 5

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateDataCenterRefNameImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "oldBadDatacetner",
				Kind: v1alpha1.VSphereDatacenterKind,
			},
			ExternalEtcdConfiguration: &v1alpha1.ExternalEtcdConfiguration{Count: 3},
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count: 3, Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.DatacenterRef.Name = "FancyNewDataCenter"

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateDataCenterRefKindImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			DatacenterRef: v1alpha1.Ref{
				Name: "oldBadDatacetner",
				Kind: v1alpha1.VSphereDatacenterKind,
			},
			ExternalEtcdConfiguration: &v1alpha1.ExternalEtcdConfiguration{Count: 3},
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count: 3, Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.DatacenterRef.Name = v1alpha1.DockerDatacenterKind

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateClusterNetworkEqualOrder(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ClusterNetwork: v1alpha1.ClusterNetwork{
				Pods: v1alpha1.Pods{
					CidrBlocks: []string{"1.2.3.4/5", "1.2.3.4/6"},
				},
				Services: v1alpha1.Services{
					CidrBlocks: []string{"1.2.3.4/7", "1.2.3.4/8"},
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ClusterNetwork = v1alpha1.ClusterNetwork{
		Pods: v1alpha1.Pods{
			CidrBlocks: []string{"1.2.3.4/6", "1.2.3.4/5"},
		},
		Services: v1alpha1.Services{
			CidrBlocks: []string{"1.2.3.4/8", "1.2.3.4/7"},
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateClusterNetworkImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ClusterNetwork: v1alpha1.ClusterNetwork{
				Pods: v1alpha1.Pods{
					CidrBlocks: []string{"1.2.3.4/5", "1.2.3.4/6"},
				},
				Services: v1alpha1.Services{
					CidrBlocks: []string{"1.2.3.4/7", "1.2.3.4/8"},
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ClusterNetwork = v1alpha1.ClusterNetwork{
		Pods: v1alpha1.Pods{
			CidrBlocks: []string{"1.2.3.4/5"},
		},
		Services: v1alpha1.Services{
			CidrBlocks: []string{"1.2.3.4/9", "1.2.3.4/10"},
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateClusterNetworkOldEmptyImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ClusterNetwork: v1alpha1.ClusterNetwork{},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ClusterNetwork = v1alpha1.ClusterNetwork{
		Pods: v1alpha1.Pods{
			CidrBlocks: []string{"1.2.3.4/5"},
		},
		Services: v1alpha1.Services{
			CidrBlocks: []string{"1.2.3.4/9", "1.2.3.4/10"},
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateClusterNetworkNewEmptyImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ClusterNetwork: v1alpha1.ClusterNetwork{
				Pods: v1alpha1.Pods{
					CidrBlocks: []string{"1.2.3.4/5", "1.2.3.4/6"},
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ClusterNetwork = v1alpha1.ClusterNetwork{
		Pods: v1alpha1.Pods{},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateProxyConfigurationEqualOrder(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ProxyConfiguration: &v1alpha1.ProxyConfiguration{
				HttpsProxy: "httpsproxy",
				NoProxy: []string{
					"noproxy1",
					"noproxy2",
				},
			},
		},
	}

	c := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ProxyConfiguration: &v1alpha1.ProxyConfiguration{
				HttpProxy:  "",
				HttpsProxy: "httpsproxy",
				NoProxy: []string{
					"noproxy2",
					"noproxy1",
				},
			},
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateProxyConfigurationImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ProxyConfiguration: &v1alpha1.ProxyConfiguration{
				HttpProxy:  "httpproxy1",
				HttpsProxy: "httpsproxy1",
				NoProxy:    []string{"noproxy1"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ProxyConfiguration = &v1alpha1.ProxyConfiguration{
		HttpProxy:  "httpproxy2",
		HttpsProxy: "httpsproxy2",
		NoProxy:    []string{"noproxy1", "noproxy2"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateProxyConfigurationNoProxyImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ProxyConfiguration: &v1alpha1.ProxyConfiguration{
				HttpProxy:  "httpproxy",
				HttpsProxy: "httpsproxy",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ProxyConfiguration = &v1alpha1.ProxyConfiguration{
		HttpProxy:  "httpproxy",
		HttpsProxy: "httpsproxy",
		NoProxy:    []string{"noproxy"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateProxyConfigurationOldNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{ProxyConfiguration: nil},
	}
	c := cOld.DeepCopy()
	c.Spec.ProxyConfiguration = &v1alpha1.ProxyConfiguration{
		HttpProxy:  "httpproxy",
		HttpsProxy: "httpsproxy",
		NoProxy:    []string{"noproxy"},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateProxyConfigurationNewNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			ProxyConfiguration: &v1alpha1.ProxyConfiguration{
				HttpProxy:  "httpproxy",
				HttpsProxy: "httpsproxy",
				NoProxy:    []string{"noproxy"},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.ProxyConfiguration = nil
	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefImmutableNilEqual(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			GitOpsRef: nil,
		},
	}
	c := cOld.DeepCopy()

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateGitOpsRefImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			GitOpsRef: &v1alpha1.Ref{
				Name: "test1", Kind: "GitOpsConfig1",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.GitOpsRef = &v1alpha1.Ref{Name: "test2", Kind: "GitOpsConfig2"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefImmutableName(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			GitOpsRef: &v1alpha1.Ref{
				Name: "test1", Kind: "GitOpsConfig",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.GitOpsRef = &v1alpha1.Ref{Name: "test2", Kind: "GitOpsConfig"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefImmutableKind(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			GitOpsRef: &v1alpha1.Ref{
				Name: "test", Kind: "GitOpsConfig1",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.GitOpsRef = &v1alpha1.Ref{Name: "test", Kind: "GitOpsConfig2"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefOldNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{GitOpsRef: nil},
	}
	c := cOld.DeepCopy()
	c.Spec.GitOpsRef = &v1alpha1.Ref{Name: "test", Kind: "GitOpsConfig"}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefNewNilImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			GitOpsRef: &v1alpha1.Ref{
				Name: "test", Kind: "GitOpsConfig",
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.GitOpsRef = nil

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateIdentityProviderRefsImmutableEqualOrder(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{
				{
					Kind: "identity",
					Name: "name1",
				},
				{
					Kind: "identity",
					Name: "name2",
				},
			},
		},
	}
	c := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{
				{
					Kind: "identity",
					Name: "name2",
				},
				{
					Kind: "identity",
					Name: "name1",
				},
			},
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateIdentityProviderRefsImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{
				{
					Kind: "identity1",
					Name: "name1",
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.IdentityProviderRefs = []v1alpha1.Ref{
		{
			Kind: "identity1",
			Name: "name1",
		},
		{
			Kind: "identity2",
			Name: "name2",
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateIdentityProviderRefsImmutableName(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{
				{
					Kind: "identity",
					Name: "name1",
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.IdentityProviderRefs = []v1alpha1.Ref{
		{
			Kind: "identity",
			Name: "name2",
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateIdentityProviderRefsImmutableKind(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{
				{
					Kind: "identity1",
					Name: "name",
				},
			},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.IdentityProviderRefs = []v1alpha1.Ref{
		{
			Kind: "identity2",
			Name: "name",
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateGitOpsRefOldEmptyImmutable(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			IdentityProviderRefs: []v1alpha1.Ref{},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.IdentityProviderRefs = []v1alpha1.Ref{
		{
			Kind: "identity",
			Name: "name",
		},
	}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateWithPausedAnnotation(t *testing.T) {
	cOld := &v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string, 1),
		},
		Spec: v1alpha1.ClusterSpec{
			KubernetesVersion: v1alpha1.Kube119,
		},
	}
	cOld.PauseReconcile()
	c := cOld.DeepCopy()
	c.Spec.KubernetesVersion = v1alpha1.Kube120

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}

func TestClusterValidateUpdateInvalidType(t *testing.T) {
	cOld := &v1alpha1.VSphereDatacenterConfig{}
	c := &v1alpha1.Cluster{}

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).NotTo(Succeed())
}

func TestClusterValidateUpdateSuccess(t *testing.T) {
	workerConfiguration := append([]v1alpha1.WorkerNodeGroupConfiguration{}, v1alpha1.WorkerNodeGroupConfiguration{Count: 5})
	cOld := &v1alpha1.Cluster{
		Spec: v1alpha1.ClusterSpec{
			WorkerNodeGroupConfigurations: workerConfiguration,
			KubernetesVersion:             v1alpha1.Kube119,
			ControlPlaneConfiguration: v1alpha1.ControlPlaneConfiguration{
				Count: 3, Endpoint: &v1alpha1.Endpoint{Host: "1.1.1.1/1"},
			},
			ExternalEtcdConfiguration: &v1alpha1.ExternalEtcdConfiguration{Count: 3},
		},
	}
	c := cOld.DeepCopy()
	c.Spec.WorkerNodeGroupConfigurations[0].Count = 10

	g := NewWithT(t)
	g.Expect(c.ValidateUpdate(cOld)).To(Succeed())
}