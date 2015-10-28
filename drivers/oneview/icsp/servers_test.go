package icsp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/docker/machine/log"
	"github.com/stretchr/testify/assert"
)

// func (s Server) GetPublicIPV4() (string, error) {
// TestGetPublicIPV4 try to test for getting interface from custom attribute
func TestGetPublicIPV4(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestGetPublicIPV4")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		serialNumber := d.Tc.GetTestData(d.Env, "FreeICSPSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber)
		testIP, err := s.GetPublicIPV4()
		assert.NoError(t, err, "Should GetPublicIPV4 without error -> %s, %+v\n", err, s)
		log.Debugf(" testIP -> %s", testIP)
		assert.True(t, (len(testIP) > 0), "Should return an ip address string")
	} else {
		// TODO: implement a test
		// need to simplate createing public_interface custom attribute object
		// need to read custom attribute object, see server_customattribute_test.go
		log.Debug("implements unit test for TestGetPublicIPV4")
	}
}

// TestGetInterfaceFromMac verify that getting an inerface with mac address works
func TestGetInterfaceFromMac(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
		s Server
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestGetInterfaceFromMac")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		macAddr := d.Tc.GetTestData(d.Env, "FreeMacAddr").(string)

		s, err := c.GetServerBySerialNumber(serialNumber)
		assert.NoError(t, err, "GetServerBySerialNumber threw error -> %s, %+v\n", err, serialNumber)

		data, err := s.GetInterfaceFromMac(macAddr)
		assert.NoError(t, err, "GetInterfaceFromMac threw error -> %s, %+v\n", err, data)
		assert.Equal(t, macAddr, data.MACAddr, "Failed to get interface -> %+v", data)
		log.Infof("Found interface -> %+v", data)
	} else {
		log.Debug("implements unit test for TestGetInterfaces")
		d, c = getTestDriverU()
		jsonServerData := d.Tc.GetTestData(d.Env, "ServerJSONString").(string)
		log.Debugf("jsonServerData => %s", jsonServerData)
		err := json.Unmarshal([]byte(jsonServerData), &s)
		assert.NoError(t, err, "Unmarshal Server threw error -> %s, %+v\n", err, jsonServerData)

		log.Debugf("server -> %v", s)

		data := s.GetInterfaces()
		log.Debugf("Interfaces -> %+v", data)
		assert.True(t, len(data) > 0, "Failed to get a valid list of interfaces -> %+v", data)
		for _, inet := range data {
			log.Debugf("inet -> %+v", inet)
			log.Debugf("inet ip -> %+v", inet.IPV4Addr)
			log.Debugf("inet ip -> %+v", inet.Slot)
			log.Infof("inet ip -> %+v", inet.MACAddr)
		}
	}
}

// TestGetInterfaces  verify that interfaces works
func TestGetInterfaces(t *testing.T) {

	var (
		d *ICSPTest
		c *ICSPClient
		s Server
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestGetInterfaces")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber)
		data := s.GetInterfaces()
		assert.NoError(t, err, "GetInterfaces threw error -> %s, %+v\n", err, data)
		assert.True(t, len(data) > 0, "Failed to get a valid list of interfaces -> %+v", data)
		for _, inet := range data {
			log.Infof("inet -> %+v", inet)
			log.Infof("inet ip -> %+v", inet.IPV4Addr)
			log.Infof("inet ip -> %+v", inet.Slot)
			log.Infof("inet ip -> %+v", inet.MACAddr)
		}
	} else {
		log.Debug("implements unit test for TestGetInterfaces")
		d, c = getTestDriverU()
		jsonServerData := d.Tc.GetTestData(d.Env, "ServerJSONString").(string)
		log.Debugf("jsonServerData => %s", jsonServerData)
		err := json.Unmarshal([]byte(jsonServerData), &s)
		assert.NoError(t, err, "Unmarshal Server threw error -> %s, %+v\n", err, jsonServerData)

		log.Debugf("server -> %v", s)

		data := s.GetInterfaces()
		log.Debugf("Interfaces -> %+v", data)
		assert.True(t, len(data) > 0, "Failed to get a valid list of interfaces -> %+v", data)
		for _, inet := range data {
			log.Debugf("inet -> %+v", inet)
			log.Debugf("inet ip -> %+v", inet.IPV4Addr)
			log.Debugf("inet ip -> %+v", inet.Slot)
			log.Debugf("inet ip -> %+v", inet.MACAddr)
		}
	}
}

// TestGetInterface implements getting the name of an interface
func TestGetInterface(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
		s Server
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestGetInterfaceName")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		serialNumber := d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		s, err := c.GetServerBySerialNumber(serialNumber)

		data, err := s.GetInterface(0)
		assert.NoError(t, err, "GetInterface threw error -> %s, %+v\n", err, data)
		assert.True(t, len(data.Slot) > 0, "Failed to get a an interface name -> %+v", data)
		log.Infof("Interface name -> %+v", data.Slot)

		data, err = s.GetInterface(1)
		assert.NoError(t, err, "GetInterface threw error -> %s, %+v\n", err, data)
		assert.True(t, len(data.Slot) > 0, "Failed to get a an interface name -> %+v", data)
		log.Infof("Interface name -> %+v", data.Slot)
	} else {
		log.Debug("implements unit test for TestGetInterfaceName")
		d, c = getTestDriverU()
		jsonServerData := d.Tc.GetTestData(d.Env, "ServerJSONString").(string)
		log.Debugf("jsonServerData => %s", jsonServerData)
		err := json.Unmarshal([]byte(jsonServerData), &s)
		assert.NoError(t, err, "Unmarshal Server threw error -> %s, %+v\n", err, jsonServerData)

		log.Debugf("server -> %v", s)

		data, err := s.GetInterface(0)
		log.Infof("Interface Name -> %+v", data)
		assert.True(t, len(data.Slot) > 0, "Failed to get a valid list of interfaces -> %+v", data)
	}
}

// TestSaveServer implement save server
func TestSaveServer(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestCreateServer")
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		// get a Server
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
	} else {
		log.Debug("implements unit test for TestCreateServer")
		var s Server
		_, c = getTestDriverU()
		s, err := c.SaveServer(s)
		assert.Error(t, err, "SaveServer threw error -> %s, %+v\n", err, s)
	}
}

// TestGetProfiles
func TestGetServers(t *testing.T) {
	var (
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		_, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetServers()
		assert.NoError(t, err, "GetServers threw error -> %s, %+v\n", err, data)

	} else {
		_, c = getTestDriverU()
		data, err := c.GetServers()
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, data))
	}
}

func TestGetServerByName(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		IcspName := d.Tc.GetTestData(d.Env, "IcspName").(string)
		data, err := c.GetServerByName(IcspName)
		assert.NoError(t, err, "GetServerByName threw error -> %s, %+v\n", err, data)

		IcspName2 := d.Tc.GetTestData(d.Env, "IcspName2").(string)
		expectsIcspName2 := d.Tc.GetExpectsData(d.Env, "IcspName2").(string)
		data, err = c.GetServerByName(IcspName2)
		assert.NoError(t, err, "GetServerByName IcspName2 threw error -> %s, %+v\n", err, data)
		assert.Equal(t, expectsIcspName2, data.Name, "GetServerByName IcspName2 on fake should be nil")
	}
}

func TestGetServerByHostName(t *testing.T) {
	var (
		d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		testHostName := d.Tc.GetTestData(d.Env, "FreeHostName").(string)
		data, err := c.GetServerByHostName(testHostName)
		assert.NoError(t, err, "GetServerByName threw error -> %s, %+v\n", err, data)
	}
}

func TestGetServerBySerialNumber(t *testing.T) {
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
		data, err := c.GetServerBySerialNumber(serialNumber)
		assert.NoError(t, err, "GetServerBySerialNumber threw error -> %s, %+v\n", err, data)

		data, err = c.GetServerBySerialNumber("2M25090RMW")
		log.Infof("server -> %+v", strings.ToLower(data.Name))

		// negative test
		data, err = c.GetServerBySerialNumber("SXXXX33333") // fake serial number
		assert.NoError(t, err, "GetServerBySerialNumber fake threw error -> %s, %+v\n", err, data)
		assert.Equal(t, data.URI.String(), "null", "GetServerBySerialNumber on fake should be nil")
	}
}

//TODO: implement test for delete
func TestDeleteServer(t *testing.T) {
	var c *ICSPClient
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		log.Debug("implements acceptance test for TestDeleteServer")
		// check if the server exist
		_, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		// mid: A unique ID assigned to the Server by Server Automation
		data, err := c.DeleteServer("510001")
		assert.True(t, data)
		assert.NoError(t, err, "DeleteServer threw error -> %s, %+v\n", err, data)
	} else {
		log.Debug("implements unit test for TestDeleteServer")
	}
}

func TestIsServerManaged(t *testing.T) {
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
		data, err := c.IsServerManaged(serialNumber)
		log.Debugf("test : %v", data)
		assert.NoError(t, err, "IsServerManaged -> %s, %+v\n", err, data)
		assert.True(t, data)
	}
}

func TestGetServerById(t *testing.T) {
	var (
		//d *ICSPTest
		c *ICSPClient
	)
	if os.Getenv("ICSP_TEST_ACCEPTANCE") == "true" {
		_, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetServerById("490001")
		assert.NoError(t, err, "GetServerByName threw error -> %s, %+v\n", err, data)
	}
}
