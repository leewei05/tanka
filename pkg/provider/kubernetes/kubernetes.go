package kubernetes

import (
	"fmt"

	"github.com/sh0rez/tanka/pkg/provider/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/objx"
)

// Kubernetes provider bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	APIServer string `mapstructure:"apiserver"`
	Namespace string `mapstructure:"namespace"`
}

var client = Kubectl{}

// Init makes the provider ready to be used
func (k *Kubernetes) Init() error {
	client.APIServer = k.APIServer
	return nil
}

// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Reconcile(raw map[string]interface{}) (state interface{}, err error) {
	docs := flattenManifest(raw)
	for _, d := range docs {
		m := objx.New(d)
		m.Set("metadata.namespace", k.Namespace)
	}
	return docs, nil
}

// flattenManifest traverses deeply nested kubernetes manifest and extracts them into a flat map.
func flattenManifest(deep map[string]interface{}) []map[string]interface{} {
	flat := []map[string]interface{}{}

	for n, d := range deep {
		if n == "__ksonnet" {
			continue
		}
		m := objx.New(d)
		if m.Has("apiVersion") && m.Has("kind") {
			flat = append(flat, m)
		} else {
			flat = append(flat, flattenManifest(m)...)
		}
	}
	return flat
}

// Fmt receives the state and reformats it to YAML Documents
func (k *Kubernetes) Fmt(state interface{}) (string, error) {
	return util.ShowYAMLDocs(state.([]map[string]interface{}))
}

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state interface{}) error {
	yaml, err := k.Fmt(state)
	if err != nil {
		return err
	}
	return client.Apply(yaml)
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state interface{}) (string, error) {
	yaml, err := k.Fmt(state)
	if err != nil {
		return "", err
	}
	return client.Diff(yaml)
}

// Cmd shall return a command to be available under `tk provider`
func (k *Kubernetes) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes provider commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented")
		},
	}
}