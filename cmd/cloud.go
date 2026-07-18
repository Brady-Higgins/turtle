package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/Brady-Higgins/turtle/internal/cloud"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
	"github.com/Brady-Higgins/turtle/internal/config"
	"github.com/Brady-Higgins/turtle/internal/docker"
	"github.com/Brady-Higgins/turtle/internal/service_client"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/spf13/cobra"
)

func activateCloud(ctx context.Context) error {
	t := cloud.InitTf()
	// check status of instance and wait to change dns record till up
	ip, err := t.CreateCloudResources()
	if err != nil {
		return err
	}
	fmt.Println(ip)
	fmt.Println("Cloud Resources created")

	// connect to cloudflare
	c, err := cloudflare_client.New()
	if err != nil {
		return err
	}
	// get old A DNS record
	// TODO: check for identical DNS record for idempotency
	record, err := c.GetDNSRecord(dns.RecordListParamsTypeCNAME, true, ctx)
	//if err != nil {
	//	return err
	//}

	// remove old DNS record
	if record != nil {
		err = c.DeleteDNSRecord(record, ctx)
		if err != nil {
			return err
		}
	}

	// create a new A record
	err = c.NewDNSRecord(ip, ctx)
	if err == nil {
		fmt.Println("Cloud Server DNS Record Created Successfully!")
	}
	return err
}

func deactivateLocal(imgName string, ctx context.Context) error {
	// stop cloudflared service
	s, err := service_client.New()
	if err != nil {
		return err
	}
	err = s.StopService()
	// TODO: ignore already stopped service
	//if err != nil {
	//	return err
	//}
	fmt.Println("Cloudflare Tunnel Service stopped")
	// connect to docker
	d, err := docker.New()
	if err != nil {
		return err
	}

	id := d.GetContainerID(imgName, ctx)
	// deactivate running container
	if id != "" && d.IsContainerRunning(id, ctx) {
		err = d.StopContainer(id, ctx)
		if err != nil {
			return err
		}
		fmt.Println("Docker container stopped")
	}
	return nil
}

func setStateCloud() (string, string, error) {
	instanceType, err := config.FindServerType()
	if err != nil {
		return "", "", err
	}
	hourlyCost := config.ServerToCost(instanceType)
	if hourlyCost == -1 {
		return "", "", errors.New("server type with unknown cost")
	}
	state, err := config.ReadStateFile()
	if state == nil {
		state = &config.State{
			TotalTime:    "0.0",
			TotalSavings: "0.0",
			LastUse:      "",
		}
	}
	recentSavings, _ := config.UpdateStateStop(state, hourlyCost)
	//if err != nil {
	//	return "", "", err
	//}
	err = config.WriteStateFile(state)
	if err != nil {
		return "", "", err
	}
	return recentSavings, state.TotalSavings, nil
}

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "switch to cloud hosting",
	Long:  "switch to cloud hosting",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config.SetupConfig()
		exists, err := config.ConfigExists()
		if err != nil || !exists {
			fmt.Println("Config doesn't exist. Please create it with turtle setup")
			return
		}
		err = config.SetTFEnv()
		if err != nil {
			fmt.Println(err)
			return
		}
		imgName, _ := cmd.Flags().GetString("image")
		if imgName == "" {
			fmt.Println("Image flag required")
			return
		}
		fmt.Println("Starting Cloud Hosting...")
		err = activateCloud(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = deactivateLocal(imgName, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		recentSavings, totalSavings, err := setStateCloud()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Cloud Hosting Started!")
		fmt.Println()
		fmt.Println("Local Hosting Session Summary:")
		fmt.Printf("    Saved $%s hosting locally\n", recentSavings)
		fmt.Printf("    Total savings of $%s\n", totalSavings)
	},
}

func init() {
	cloudCmd.Flags().String("image", "", "docker image name to run")
}
