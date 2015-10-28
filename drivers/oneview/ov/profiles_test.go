package ov

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/machine/log"
	"github.com/stretchr/testify/assert"
)

// CreateProfileFromTemplate(name string, new_template ServerProfile, blade ServerHardware)
// test create profile
func TestCreateProfileFromTemplate(t *testing.T) {
	var (
		d                *OVTest
		c                *OVClient
		testHostName     string
		testBladeSerial  string
		testTemplateName string
		testBlades       ServerHardwareList
		testTemplate     ServerProfile
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		testHostName = d.Tc.GetTestData(d.Env, "HostName").(string)
		testBladeSerial = d.Tc.GetTestData(d.Env, "FreeBladeSerialNumber").(string)
		testTemplateName = d.Tc.GetTestData(d.Env, "TemplateProfile").(string)

		testBlades, _ = c.GetServerHardwareList([]string{fmt.Sprintf("serialNumber matches '%s'", testBladeSerial)}, "name:asc")
		assert.True(t, (len(testBlades.Members) > 0), "Did not get any blades from server hardware list")

		testTemplate, _ = c.GetProfileTemplateByName(testTemplateName)
		assert.Equal(t, testTemplateName, testTemplate.Name, fmt.Sprintf("Problem getting template name, %+v", testTemplate))

		// find out if the test profile already exist
		testProfile, err := c.GetProfileByName(testHostName)
		assert.NoError(t, err, "CreateProfileFromTemplate get the server profile error -> %s", err)

		if len(testBlades.Members) > 0 && testProfile.URI.IsNil() {
			err := c.CreateProfileFromTemplate(testHostName, testTemplate, testBlades.Members[0])
			assert.NoError(t, err, "CreateProfileFromTemplate error -> %s", err)

			err = c.CreateProfileFromTemplate(testHostName, testTemplate, testBlades.Members[0])
			assert.Error(t, err, "CreateProfileFromTemplate should error because a template already exist, err-> %s", err)

		} else {
			log.Warnf("The testHostName already exist, so skipping CreateProfileTemplate test for %s", testHostName)
		}

		// reload the test profile that we just created
		testProfile, err = c.GetProfileByName(testHostName)
		assert.NoError(t, err, "CreateProfileFromTemplate get the server profile error -> %s", err)

		// get the server hardware associated with that test profile
		log.Debugf("testProfile -> %+v", testProfile)
		testBlade, err := c.GetServerHardware(testProfile.ServerHardwareURI)
		assert.NoError(t, err, "CreateProfileFromTemplate call to GetServerHardware got error -> %s", err)

		// power on the server, and leave it in that state
		var pt *PowerTask
		log.Debugf("testBlade -> %+v", testBlade)
		pt = pt.NewPowerTask(testBlade)
		pt.Timeout = 46 // timeout is 20 sec
		log.Info("------- Setting Power to On")
		err = pt.PowerExecutor(P_ON)
		assert.NoError(t, err, "PowerExecutor threw no errors -> %s", err)
	}

}

// find Server_Profile_scs79
func TestGetProfileByName(t *testing.T) {
	var (
		d        *OVTest
		c        *OVClient
		testname string
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		testname = d.Tc.GetTestData(d.Env, "ServerProfileName").(string)
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetProfileByName(testname)
		assert.NoError(t, err, "GetProfileByName threw error -> %s", err)
		assert.Equal(t, testname, data.Name)

		data, err = c.GetProfileByName("foo")
		assert.NoError(t, err, "GetProfileByName with fake name -> %+v", err)
		assert.Equal(t, "", data.Name)

	} else {
		d, c = getTestDriverU()
		testname = d.Tc.GetTestData(d.Env, "ServerProfileName").(string)
		data, err := c.GetProfileByName(testname)
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, data))
	}
}

// TestGetProfileConnectionByName verify functionality
func TestGetConnectionByName(t *testing.T) {
	var (
		d        *OVTest
		c        *OVClient
		testname string
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		testname = d.Tc.GetTestData(d.Env, "ServerProfileName").(string)
		pubcname := d.Tc.GetTestData(d.Env, "FreePublicConnection").(string)
		expectsmac := d.Tc.GetExpectsData(d.Env, "MACAddress").(string)
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		profile, err := c.GetProfileByName(testname)
		assert.NoError(t, err, "GetProfileByName threw error -> %s", err)

		for _, c := range profile.Connections {
			log.Debugf("connection -> %d %s %s %s", c.ID, c.Name, c.MAC, c.MacType)
		}

		connection, err := profile.GetConnectionByName(pubcname)
		log.Debugf("Got connection -> %+v", connection)
		assert.NoError(t, err, "GetConnectionByName threw error -> %s", err)
		assert.Equal(t, expectsmac, connection.MAC.String(), "GetConnectionByName failed to get connection %+v", err, connection)
	}
}

// TestGetProfileNameBySN
// Acceptance test ->
// /rest/server-profiles
// ?filter=serialNumber matches '2M25090RMW'&sort=name:asc
func TestGetProfileBySN(t *testing.T) {
	var (
		d          *OVTest
		c          *OVClient
		testSerial string
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		testSerial = d.Tc.GetTestData(d.Env, "SerialNumber").(string)
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetProfileBySN(testSerial)
		assert.NoError(t, err, "GetProfileBySN threw error -> %s", err)
		assert.Equal(t, testSerial, data.SerialNumber)

		data, err = c.GetProfileBySN("foo")
		assert.NoError(t, err, "GetProfileBySN with fake serial number -> %+v", err)
		assert.Equal(t, "", data.SerialNumber)

	} else {
		d, c = getTestDriverU()
		testSerial = d.Tc.GetTestData(d.Env, "SerialNumber").(string)
		data, err := c.GetProfileBySN(testSerial)
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, data))
	}
}

// TestGetProfiles
func TestGetProfiles(t *testing.T) {
	var (
		// d *OVTest
		c *OVClient
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		_, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetProfiles("", "")
		assert.NoError(t, err, "GetProfiles threw error -> %s, %+v\n", err, data)

		data, err = c.GetProfiles("", "name:asc")
		assert.NoError(t, err, "GetProfiles name:asc error -> %s, %+v", err, data)

	} else {
		_, c = getTestDriverU()
		data, err := c.GetProfiles("", "")
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, data))
	}
}

// test for not found profile
// should not delete a profile that doesn't exist
func TestDeleteProfileNotFound(t *testing.T) {
	var (
		c               *OVClient
		testProfileName = "fake_profile_doesnt_exist"
		testProfile     ServerProfile
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		_, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}

		err := c.DeleteProfile(testProfileName)
		assert.NoError(t, err, "DeleteProfile err-> %s", err)

		testProfile, err = c.GetProfileByName(testProfileName)
		assert.NoError(t, err, "GetProfileByName with deleted profile -> %+v", err)
		assert.Equal(t, "", testProfile.Name, fmt.Sprintf("Problem getting template name, %+v", testProfile))
	} else {
		_, c = getTestDriverU()
		err := c.DeleteProfile(testProfileName)
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, testProfile))
	}
}

// test DeleteProfile
func TestDeleteProfile(t *testing.T) {
	var (
		d               *OVTest
		c               *OVClient
		testProfileName string
		testProfile     ServerProfile
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		testProfileName = d.Tc.GetTestData(d.Env, "HostName").(string)

		err := c.DeleteProfile(testProfileName)
		assert.NoError(t, err, "DeleteProfile err-> %s", err)

		testProfile, err = c.GetProfileByName(testProfileName)
		assert.NoError(t, err, "GetProfileByName with deleted profile -> %+v", err)
		assert.Equal(t, "", testProfile.Name, fmt.Sprintf("Problem getting template name, %+v", testProfile))
	} else {
		_, c = getTestDriverU()
		err := c.DeleteProfile("footest")
		assert.Error(t, err, fmt.Sprintf("ALL ok, no error, caught as expected: %s,%+v\n", err, testProfile))
	}

}
