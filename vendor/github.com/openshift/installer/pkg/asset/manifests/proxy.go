package manifests

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

var proxyCfgFilename = filepath.Join(manifestDir, "cluster-proxy-01-config.yaml")

// Proxy generates the cluster-proxy-*.yml files.
type Proxy struct {
	FileList []*asset.File
	Config   *configv1.Proxy
}

var _ asset.WritableAsset = (*Proxy)(nil)

// Name returns a human-friendly name for the asset.
func (*Proxy) Name() string {
	return "Proxy Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Proxy) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&Networking{},
	}
}

// Generate generates the Proxy config and its CRD.
func (p *Proxy) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	network := &Networking{}
	dependencies.Get(installConfig, network)

	p.Config = &configv1.Proxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Proxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
	}

	if installConfig.Config.Proxy != nil {
		p.Config.Spec = configv1.ProxySpec{
			HTTPProxy:  installConfig.Config.Proxy.HTTPProxy,
			HTTPSProxy: installConfig.Config.Proxy.HTTPSProxy,
			NoProxy:    installConfig.Config.Proxy.NoProxy,
		}
	}

	if p.Config.Spec.HTTPProxy != "" || p.Config.Spec.HTTPSProxy != "" {
		noProxy, err := createNoProxy(installConfig, network)
		if err != nil {
			return err
		}
		p.Config.Status = configv1.ProxyStatus{
			HTTPProxy:  installConfig.Config.Proxy.HTTPProxy,
			HTTPSProxy: installConfig.Config.Proxy.HTTPSProxy,
			NoProxy:    noProxy,
		}
	}

	configData, err := yaml.Marshal(p.Config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", p.Name())
	}

	p.FileList = []*asset.File{
		{
			Filename: proxyCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

// createNoProxy combines user-provided & platform-specific values to create a comma-separated
// list of unique NO_PROXY values. Platform values are: serviceCIDR, podCIDR, localhost,
// 127.0.0.1, api.clusterdomain, api-int.clusterdomain, etcd-idx.clusterdomain
func createNoProxy(installConfig *installconfig.InstallConfig, network *Networking) (string, error) {
	apiServerURL, err := url.Parse(getAPIServerURL(installConfig.Config))
	if err != nil {
		return "", errors.New("failed parsing API server when creating Proxy manifest")
	}
	internalAPIServer, err := url.Parse(getInternalAPIServerURL(installConfig.Config))
	if err != nil {
		return "", errors.New("failed parsing internal API server when creating Proxy manifest")
	}

	set := sets.NewString(
		"127.0.0.1",
		"localhost",
		network.Config.Spec.ServiceNetwork[0],
		apiServerURL.Hostname(),
		internalAPIServer.Hostname(),
	)

	for i := int64(0); i < *installConfig.Config.ControlPlane.Replicas; i++ {
		etcdHost := fmt.Sprintf("etcd-%d.%s", i, installConfig.Config.ClusterDomain())
		set.Insert(etcdHost)
	}

	for _, clusterNetwork := range network.Config.Spec.ClusterNetwork {
		set.Insert(clusterNetwork.CIDR)
	}

	for _, userValue := range strings.Split(installConfig.Config.Proxy.NoProxy, ",") {
		set.Insert(userValue)
	}

	return strings.Join(set.List(), ","), nil
}

// Files returns the files generated by the asset.
func (p *Proxy) Files() []*asset.File {
	return p.FileList
}

// Load loads the already-rendered files back from disk.
func (p *Proxy) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
