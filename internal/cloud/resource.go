package cloud

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type TfClient struct {
	tf *tfexec.Terraform
}

func initTf() *TfClient {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir, _ := os.Getwd()
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	return &TfClient{tf: tf}
}

func (t *TfClient) CreateCloudResources() error {
	err := t.tf.Apply(context.Background())
	return err
}

func (t *TfClient) DestroyCloudResources() error {
	err := t.tf.Destroy(context.Background())
	return err
}
