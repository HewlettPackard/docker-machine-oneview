package icsp

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/docker/machine/drivers/oneview/rest"
	"github.com/docker/machine/drivers/oneview/testconfig"
	"github.com/docker/machine/drivers/oneview/utils"
	"github.com/docker/machine/log"
	"github.com/stretchr/testify/assert"
)

//TODO: need to learn a better way of how integration testing works with bats
// and use that instead
// Acceptance test
// 1) craete an environment, /.oneview.env, script to export these values:
// ONEVIEW_ICSP_ENDPOINT=https://blah
// ONEVIEW_ICSP_PASSWORD=pass
// ONEVIEW_ICSP_USER=user
// ONEVIEW_ICSP_DOMAIN=LOCAL
// ONEVIEW_SSLVERIFY=true
// 2) setup gotest container
//    docker run -it --env-file ./.oneview.env -v $(pwd):/go/src/github.com/docker/machine --name testicsp docker-machine
//    exit
//    docker start goaccept
// 3) setup alias:
//   alias testicsp='docker exec goaccept godep go test -test.v=true --short'
// 4) to refresh env options, use
//    docker rm -f testicsp
//    ... repeat steps 2
type ICSPTest struct {
	Tc     *testconfig.TestConfig
	Client *ICSPClient
	Env    string
}

// get Environment
func (ot *ICSPTest) GetEnvironment() {
	if os.Getenv("ONEVIEW_TEST_ENV") != "" {
		ot.Env = os.Getenv("ONEVIEW_TEST_ENV")
		return
	}
	ot.Env = "dev"
	return
}

// get a test driver for acceptance testing
func getTestDriverA() (*ICSPTest, *ICSPClient) {
	// os.Setenv("DEBUG", "true")  // remove comment to debug logs
	var ot *ICSPTest
	var tc *testconfig.TestConfig
	ot = &ICSPTest{Tc: tc.NewTestConfig(), Env: "dev"}
	ot.GetEnvironment()
	ot.Tc.GetTestingConfiguration(os.Getenv("ONEVIEW_TEST_DATA"))
	ot.Client = &ICSPClient{
		rest.Client{
			User:     os.Getenv("ONEVIEW_ICSP_USER"),
			Password: os.Getenv("ONEVIEW_ICSP_PASSWORD"),
			Domain:   os.Getenv("ONEVIEW_ICSP_DOMAIN"),
			Endpoint: os.Getenv("ONEVIEW_ICSP_ENDPOINT"),
			// ConfigDir:
			SSLVerify: false,
			APIKey:    "none",
		},
	}
	ot.Client.RefreshVersion()
	// fmt.Println("Setting up test with getTestDriverA")
	return ot, ot.Client
}

// Unit test
func getTestDriverU() (*ICSPTest, *ICSPClient) {
	var ot *ICSPTest
	var tc *testconfig.TestConfig
	ot = &ICSPTest{Tc: tc.NewTestConfig(), Env: "dev"}
	ot.GetEnvironment()
	ot.Tc.GetTestingConfiguration(os.Getenv("ONEVIEW_TEST_DATA"))
	ot.Client = &ICSPClient{
		rest.Client{
			User:       "foo",
			Password:   "bar",
			Domain:     "LOCAL",
			Endpoint:   "https://icsptestcase",
			SSLVerify:  false,
			APIVersion: 108,
			APIKey:     "none",
		},
	}
	// fmt.Println("Setting up test with getTestDriverU")
	return ot, ot.Client
}

// implement create server unt test
//TODO: This test requires a server profile to have been created
func TestCreateServer(t *testing.T) {
	var (
		d              *ICSPTest
		c              *ICSPClient
		user, pass, ip string
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		user = os.Getenv("ONEVIEW_ILO_USER")
		pass = os.Getenv("ONEVIEW_ILO_PASSWORD")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		ip = d.Tc.GetTestData(d.Env, "IloIPAddress").(string)
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		log.Debug("implements acceptance test for TestCreateServer")
		s, err := c.GetServerBySerialNumber(serialNumber) // fake serial number
		assert.NoError(t, err, "GetServerBySerialNumber fake threw error -> %s, %+v\n", err, s)
		if s.URI.IsNil() {
			// create the server
			err := c.CreateServer(user, pass, ip, 443)
			assert.NoError(t, err, "CreateServer threw error -> %s\n", err)
		} else {
			// create the server
			err := c.CreateServer(user, pass, ip, 443)
			assert.Error(t, err, "CreateServer should throw conflict error  -> %s\n", err)
		}
		// check if the server now exist
	} else {
		_, c = getTestDriverU()
		log.Debug("implements unit test for TestCreateServer")
		err := c.CreateServer("foo", "bar", "127.0.0.1", 443)
		assert.Error(t, err, "CreateServer should throw error  -> %s\n", err)
	}
}

// TestInterface verify we can parse interfaces
func TestInterface(t *testing.T) {
	var testSlots []string
	testSlots = []string{"eth0", "eth1", "eth2"}
	for i, slot := range testSlots {
		x, _ := strconv.Atoi(slot[len(slot)-1:])
		assert.Equal(t, i, x, "the slot should match the index")
	}
}

// TestPreApplyDeploymentJobs - setup some information from icsp
//TODO: This test requires a server profile to have been created
func TestPreApplyDeploymentJobs(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for ApplyDeploymentJobs")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber)
		assert.NoError(t, err, "GetServerBySerialNumber threw error -> %s, %+v\n", err, s)

		log.Infof("server opslf -> %+v", s.OpswLifecycle)
		assert.True(t, PRE_UNPROVISIONED.Equal(s.OpswLifecycle), "should be unprovisioned")

		log.Infof("server state -> %+v", s.State)
		assert.True(t, S_MAINTENANCE.Equal(s.State), "server should be in maintenance mode")

		pubinet, err := s.GetInterface(1)
		err = c.PreApplyDeploymentJobs(s, pubinet) // responsible for configuring the Pulbic IP CustomAttributes
		assert.NoError(t, err, "ApplyDeploymentJobs threw error -> %+v, %+v", err, s)

		// verify that the server attribute was saved by getting the server again and checking the value
		_, testValue2 := s.GetValueItem("public_interface", "server")
		// unmarshal the custom attribute
		var inet *Interface
		err = json.Unmarshal([]byte(testValue2.Value), &inet)
		assert.NoError(t, err, "Unmarshal Interface threw error -> %s, %+v\n", err, testValue2.Value)

		log.Infof("We got public ip addr -> %s", inet.IPV4Addr)
		assert.True(t, len(inet.IPV4Addr) > 0, "Should return the saved custom attribute for ipaddress")
	}
}

// integrated acceptance test
// TestSaveServer implement save server
//TODO: a workaround to figuring out how to bubble up public ip address information from the os to icsp after os build plan provisioning
func TestApplyDeploymentJobs(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for ApplyDeploymentJobs")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		// get a Server
		osBuildPlan := d.Tc.GetTestData(d.Env, "OSBuildPlan").(string)
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber)
		assert.NoError(t, err, "GetServerBySerialNumber threw error -> %s, %+v\n", err, s)
		// set a custom attribute
		s.SetCustomAttribute("docker_user", "server", "docker")
		// use test keys like from https://github.com/mitchellh/vagrant/tree/master/keys
		// private key from https://raw.githubusercontent.com/mitchellh/vagrant/master/keys/vagrant
		s.SetCustomAttribute("public_key", "server", "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ==")

		// save a server
		news, err := c.SaveServer(s)
		assert.NoError(t, err, "SaveServer threw error -> %s, %+v\n", err, news)
		assert.Equal(t, s.UUID, news.UUID, "Should return a server with the same UUID")

		// verify that the server attribute was saved by getting the server again and checking the value
		_, testValue2 := s.GetValueItem("docker_user", "server")
		assert.Equal(t, "docker", testValue2.Value, "Should return the saved custom attribute")

		_, err = c.ApplyDeploymentJobs(osBuildPlan, s)
		assert.NoError(t, err, "ApplyDeploymentJobs threw error -> %s, %+v\n", err, news)
	} else {
		var s Server
		_, c = getTestDriverU()
		log.Debug("implements unit test for ApplyDeploymentJobs")
		_, err := c.ApplyDeploymentJobs("testbuildplan", s)
		assert.Error(t, err, "ApplyDeploymentJobs threw error -> %s, %+v\n", err, s)
	}
}

//TestPostApplyDeploymentJobs test job Task
//TODO: This test requires a server profile to have been created
func TestPostApplyDeploymentJobs(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}

		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber) // fake serial number

		// (c *ICSPClient) GetJob(u ODSUri) (Job, error) {
		// create a jt *JobTask object
		// JobURI
		var jt *JobTask
		var testURL utils.Nstring
		testURL = "/rest/os-deployment-jobs/5350001"
		jt = &JobTask{
			JobURI: ODSUri{URI: testURL},
			Client: c,
		}
		var findprops []string
		findprops = append(findprops, "public_ip")
		err = c.PostApplyDeploymentJobs(jt, s, findprops)
		assert.NoError(t, err, "PostApplyDeploymentJobs threw error -> %s, %+v\n", err, jt)
	}
}

// TestCustomServerAttributes test CustomServerAttributes
func TestCustomServerAttributes(t *testing.T) {
	var testServerAttributes *CustomServerAttributes
	testServerAttributes = testServerAttributes.New()
	testServerAttributes.Set("ssh_user", "docker")
	log.Infof("testServer -> %+v", testServerAttributes)
	assert.Equal(t, "docker", testServerAttributes.Get("ssh_user"), "should be equal when called for ssh_user")
}
