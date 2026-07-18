package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Brady-Higgins/turtle/internal/cloud"
	"github.com/Brady-Higgins/turtle/internal/cloudflare_client"
	"github.com/Brady-Higgins/turtle/internal/config"
	"github.com/Brady-Higgins/turtle/internal/docker"
	"github.com/Brady-Higgins/turtle/internal/service_client"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/spf13/cobra"
)

func activateLocal(c *cloudflare_client.CloudflareClient, imgName string, tunnelName string, hostName string, ctx context.Context) error {
	// connect to docker
	d, err := docker.New()
	if err != nil {
		return err
	}

	id := d.GetContainerID(imgName, ctx)
	// container already exists for image
	if id != "" {
		err = d.StartContainer(id, ctx)
	} else { // need to build container first
		id, err = d.BuildContainer(imgName, ctx)
		d.StartContainer(id, ctx)
	}
	if err != nil {
		return err
	}
	fmt.Println("Docker container started")
	s, err := service_client.New()
	if err != nil {
		return err
	}
	// start cloudflared service
	fmt.Println("Cloudflare Tunnel Service started")
	err = s.StartService()
	time.Sleep(5 * time.Second)
	// get old A DNS record
	record, err := c.GetDNSRecord(dns.RecordListParamsTypeA, true, ctx)
	//if err != nil {
	//	return err
	//}
	// if A record exists
	if record != nil {
		// remove old DNS record
		err = c.DeleteDNSRecord(record, ctx)
		if err != nil {
			return err
		}
	}
	// create tunnel DNS record. CNAME
	err = cloudflare_client.CreateTunnelDNSRecord(tunnelName, hostName)
	if err == nil {
		fmt.Println("Tunnel DNS Record Created Successfully")
	}
	return err

}

func deactivateCloud(c *cloudflare_client.CloudflareClient, ctx context.Context) error {

	// get new tunnel record
	tunnelRecord, err := c.GetDNSRecord(dns.RecordListParamsTypeCNAME, false, ctx)
	if tunnelRecord != nil {
		// comment it for easy identification later
		err = c.CommentDNSRecord(tunnelRecord, ctx)
		if err != nil {
			return err
		}
	}

	// destroy cloud resources
	t := cloud.InitTf()
	// make this a go routine
	err = t.DestroyCloudResources()
	fmt.Println("Cloud resources destroyed")
	return err
}

func setStateLocal() error {
	instanceType, err := config.FindServerType()
	if err != nil {
		return err
	}
	hourlyCost := config.ServerToCost(instanceType)
	if hourlyCost == -1 {
		return errors.New("server type with unknown cost")
	}
	state, _ := config.ReadStateFile()
	if state == nil {
		state = &config.State{
			TotalTime:    "0.0",
			TotalSavings: "0.0",
			LastUse:      "",
		}
	}
	config.UpdateStateStart(state)
	err = config.WriteStateFile(state)
	return err
}

func checkSelfHosting() bool {
	return false
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "switch to local hosting",
	Long:  "switch to local hosting",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config.SetupConfig()
		exists, err := config.ConfigExists()
		if err != nil || !exists {
			fmt.Println(err)
			fmt.Println("Config doesn't exist. Please create it with turtle setup")
			return
		}
		err = config.SetTFEnv()
		if err != nil {
			fmt.Println(err)
			return
		}
		//imgName := "example-site"
		imgName, err := cmd.Flags().GetString("image")
		if imgName == "" {
			fmt.Println("Image flag required")
			return
		}
		tunnelName := config.GetConfigValue("cloudflare_tunnel_name")
		if tunnelName == "" {
			fmt.Println("Cloudflare tunnel name isn't set in your environment")
			return
		}
		hostName := config.GetConfigValue("host_name")
		if hostName == "" {
			fmt.Println("Host name isn't set in your environment")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Starting Self Hosting...")
		// connect to cloudflare
		c, err := cloudflare_client.New()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = activateLocal(c, imgName, tunnelName, hostName, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = deactivateCloud(c, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = setStateLocal()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Self Hosting Started!")
		fmt.Printf("Docker Image: %s, Cloudflare Tunnel: %s, Host name: %s\n", imgName, tunnelName, hostName)

	},
}

func init() {
	localCmd.Flags().String("image", "", "docker image name to run")
}
