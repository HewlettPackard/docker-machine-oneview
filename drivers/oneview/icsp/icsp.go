package icsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/machine/drivers/oneview/rest"
	"github.com/docker/machine/log"
)

// ICSPClient - wrapper class for icsp api's
type ICSPClient struct {
	rest.Client
}

// new Client
func (c *ICSPClient) NewICSPClient(user string, password string, domain string, endpoint string, sslverify bool, apiversion int) *ICSPClient {
	return &ICSPClient{
		rest.Client{
			User:       user,
			Password:   password,
			Domain:     domain,
			Endpoint:   endpoint,
			SSLVerify:  sslverify,
			APIVersion: apiversion,
			APIKey:     "none",
		},
	}
}

// CustomServerAttributes setup custom attributes to apply to a server
type CustomServerAttributes struct {
	Values map[string]string // Hash of strings for custom attributes
}

func (cs *CustomServerAttributes) New() *CustomServerAttributes {
	return &CustomServerAttributes{
		Values: make(map[string]string),
	}
}

func (cs *CustomServerAttributes) Set(key string, value string) {
	cs.Values[key] = value
}

func (cs *CustomServerAttributes) Get(key string) string {
	return cs.Values[key]
}

// CustomizeServer - use customizeserver when working with creating a new server
// server create if it's missing
// server apply deployment job
type CustomizeServer struct {
	HostName         string                  // provide a hostname to set
	SerialNumber     string                  // should be the serial number for the server
	ILoUser          string                  // should be the user name for ilo administration
	IloPassword      string                  // should be the ilo password to use
	IloIPAddress     string                  // PXE ip address for ilo
	IloPort          int                     // port number for ilo server
	OSBuildPlan      string                  // name of the OS build plan
	ServerProperties *CustomServerAttributes // name value pairs for server custom attributes
	PublicSlotID     int                     // the public interface that will be used to get public ipaddress
}

// PostApplyDeploymentJobs - post deployment task to update custom attributes with
// results of a job task that was executed on the server
func (c *ICSPClient) PostApplyDeploymentJobs(jt *JobTask, s Server, properties []string) error {
	// look at jobResultLogDetails, parse *=* strings
	job, err := c.GetJob(jt.JobURI)
	if err != nil {
		return err
	}
	for _, result := range job.JobResult {
		for _, line := range strings.Split(result.JobResultLogDetails, "\n") {
			r := regexp.MustCompile("(.*)=(.*)")
			if r.FindString(line) != "" {
				for _, property := range properties {
					a := r.FindStringSubmatch(line)
					if len(a) >= 3 && property == a[1] {
						s.SetCustomAttribute(a[1], "server", a[2])
					}
				}
			}
		}
	}
	_, err = c.SaveServer(s)
	// place those strings into custom attributes
	return err
}

// PreApplyDeploymentJobs - will attempt to identify the public interface for this job
//  for now we simply look for interfaces on ethx and save those into a custom attribute called
//  PublicInterface, this can be controlled by providing the slot id.   Example:
//  slotid = x , will map to ethx
//  slotid = 1 , will map to eth1
//  slotid = 2 , will map to eth2
//  public_interface can only be called when the server is in maintenance mode, all others
//   simply fall out
//TODO: a workaround to figuring out how to bubble up public ip address information from the os to icsp after os build plan provisioning
func (c *ICSPClient) PreApplyDeploymentJobs(s Server, slotid int) error {
	var publicinterface Interface
	if PRE_UNPROVISIONED.Equal(s.OpswLifecycle) && S_MAINTENANCE.Equal(s.State) {
		log.Debugf("Applying pre deployment job settings")
		inets := s.GetInterfaces()
		for i, inet := range inets {
			if i == slotid {
				publicinterface = inet
				break
			}
		}
		// json version of the publicinterface
		publicinterfacejson, err := json.Marshal(publicinterface)
		if err != nil {
			return err
		}
		// save the publicinterface into a custom attribute called public_interface
		s.SetCustomAttribute("public_interface", "server", fmt.Sprintf("%s", bytes.NewBuffer(publicinterfacejson)))

		// save it
		_, err = c.SaveServer(s)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Skippling the pre-apply deployment jobs settings")
	}
	return nil
}

// Customize Server
func (c *ICSPClient) CustomizeServer(cs CustomizeServer) error {
	s, err := c.GetServerBySerialNumber(cs.SerialNumber)
	if err != nil {
		return err
	}
	if s.SerialNumber != cs.SerialNumber {
		log.Infof("ICSP creating server for : %s", cs.IloIPAddress)
		if err := c.CreateServer(cs.ILoUser, cs.IloPassword, cs.IloIPAddress, cs.IloPort); err != nil {
			return err
		}
	} else {
		log.Infof("ICSP server was already created, %s, skipping", cs.HostName)
	}

	// reload that server
	s, err = c.GetServerBySerialNumber(cs.SerialNumber)
	if err != nil {
		return err
	}

	// verify that the server actually has a URI
	if s.URI.IsNil() {
		return fmt.Errorf("Server customization failure, unable to get a valid server from icsp for serial number: %s", cs.SerialNumber)
	}

	// verify that the server retrieved matches it's serial number
	if s.SerialNumber != cs.SerialNumber {
		return fmt.Errorf("Server customization failure, server serial numbers mismatch for %s.", cs.SerialNumber)
	}

	// handle getting interface name
	intname, err := s.GetInterfaceName(cs.PublicSlotID)
	if err != nil {
		return err
	}

	// save the server attributes to the server
	for k, v := range cs.ServerProperties.Values {
		// handle sepecial custom attributes
		// handle @server_name@ and replace for s.Name
		v = strings.Replace(v, "@server_name@", s.Name, -1)
		v = strings.Replace(v, "@interface@", intname, -1)
		s.SetCustomAttribute(k, "server", v)
	}

	// save it
	new_server, err := c.SaveServer(s)
	if err != nil {
		return err
	}

	// call to capture the public_interface attribute
	if err := c.PreApplyDeploymentJobs(new_server, cs.PublicSlotID); err != nil {
		return err
	}

	// apply the build Plan
	jt, err := c.ApplyDeploymentJobs(cs.OSBuildPlan, new_server)
	if err != nil {
		return err
	}

	// use jt to get additional customizations we can use on the server custom attributes
	// TODO: this needs to be evaluated on usefull ness and proper way to pass up additional deployment information back to the server in icsp
	var findprops []string
	findprops = append(findprops, "public_ip")
	return c.PostApplyDeploymentJobs(jt, new_server, findprops)
}
