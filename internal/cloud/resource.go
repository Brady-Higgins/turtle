package cloud

import (
	"context"
	"encoding/json"
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

// TODO: add a check if a turtle server instance is already active

func InitTf() *TfClient {
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
	// TODO: set state file
	return &TfClient{tf: tf}
}

func (t *TfClient) CreateCloudResources() (string, error) {
	err := t.tf.Apply(context.Background())
	if err != nil {
		return "", err
	}
	resp, err := t.tf.Output(context.Background())
	if err != nil {
		return "", err
	}
	var ip string
	err = json.Unmarshal(resp["public_ip"].Value, &ip)
	if err != nil {
		return "", err
	}
	return ip, err
}

func (t *TfClient) DestroyCloudResources() error {
	err := t.tf.Destroy(context.Background())
	return err
}
