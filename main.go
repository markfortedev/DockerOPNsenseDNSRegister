package main

import (
	"DockerOPNsenseDNSRegister/docker"
	opnsenseproxyapi "DockerOPNsenseDNSRegister/opnsense-proxy-api"
	"fmt"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"time"
)

var hostFQDN string
var currentVirtualHosts []string
var opnsenseProxyAddress string

func main() {
	hostname := os.Getenv("HOST_HOSTNAME")
	if hostname == "" {
		log.Fatalf("HOST_HOSTNAME is not set")
	}
	domainName := os.Getenv("DOMAIN_NAME")
	if domainName == "" {
		log.Fatalf("DOMAIN_NAME is not set")
	}
	opnsenseProxyAddress = os.Getenv("OPNSENSE_PROXY_ADDRESS")
	if opnsenseProxyAddress == "" {
		log.Fatalf("OPNSENSE_PROXY_ADDRESS is not set")
	}
	if !strings.Contains(opnsenseProxyAddress, "http://") {
		opnsenseProxyAddress = fmt.Sprintf("http://%v", opnsenseProxyAddress)
	}
	currentVirtualHosts = []string{}
	hostFQDN = fmt.Sprintf("%v-internal-reverse-proxy.%v", hostname, domainName)
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).Minutes().Do(updateVirtualHosts)
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	log.Infof("Starting DNS Syncing")
	scheduler.StartBlocking()
}

func updateVirtualHosts() {
	virtualHosts, err := docker.GetVirtualHosts()
	if err != nil {
		log.Errorf("Error getting virtual hosts from Docker daemon: %v", err)
	}
	if reflect.DeepEqual(virtualHosts, currentVirtualHosts) {
		return
	}
	log.Infof("Virtual hosts have changed. Updating OPNsense. New Virtual Hosts: %v", strings.Join(virtualHosts, ", "))
	currentVirtualHosts = virtualHosts
	err = opnsenseproxyapi.Sync(opnsenseProxyAddress, virtualHosts, hostFQDN)
	if err != nil {
		log.Errorf("Error syncing virtual hosts with OPNsense Proxy: %v", err)
	}
}
