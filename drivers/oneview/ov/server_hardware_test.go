package ov
import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/docker/machine/log"
	"github.com/docker/machine/drivers/oneview/utils"
)

// get server hardware test
func TestServerHardware(t *testing.T) {
	var (
		d *OVTest
		c *OVClient
		testData utils.Nstring
		expectsData string
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		testData    = utils.Nstring(d.Tc.GetTestData(d.Env, "ServerHardwareURI").(string))
		expectsData = d.Tc.GetExpectsData(d.Env, "SerialNumber").(string)
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetServerHardware(testData)
		log.Debugf("%+v", data)
		assert.NoError(t, err, "GetServerHardware threw error -> %s", err)
		// fmt.Printf("data.Connections -> %+v\n", data)
		assert.Equal(t, expectsData, data.SerialNumber)

	}
}

// get server hardware test
func TestGetAvailableHardware(t *testing.T) {
	var (
		d *OVTest
		c *OVClient
		testHwType_URI  utils.Nstring
		testHWGroup_URI utils.Nstring
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		testHwType_URI  = utils.Nstring(d.Tc.GetTestData(d.Env, "HardwareTypeURI").(string))
		testHWGroup_URI = utils.Nstring(d.Tc.GetTestData(d.Env, "GroupURI").(string))
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetAvailableHardware(testHwType_URI, testHWGroup_URI)
		assert.NoError(t, err, "GetAvailableHardware threw error -> %s", err)
		// fmt.Printf("data.Connections -> %+v\n", data)
		log.Debugf("Abailable server -> %+v", data)
		assert.NotEqual(t, "", data.Name)

	}
}
