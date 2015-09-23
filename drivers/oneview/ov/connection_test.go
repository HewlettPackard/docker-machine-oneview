package ov
import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test working with connections
// Acceptance test ->
// /rest/server-profiles
// ?filter=serialNumber matches '2M25090RMW'&sort=name:asc
func TestConnections(t *testing.T) {
	var (
		d *OVTest
		c *OVClient
	)
	if os.Getenv("ONEVIEW_TEST_ACCEPTANCE") == "true" {
		d, c = getTestDriverA()
		if c == nil {
			t.Fatalf("Failed to execute getTestDriver() ")
		}
		data, err := c.GetProfileBySN(d.Tc.GetTestData(d.Env, "SerialNumber").(string))
		assert.NoError(t, err, "GetProfileBySN threw error -> %s", err)
		// fmt.Printf("data.Connections -> %+v\n", data)
		if len(data.Connections) > 0 {
			assert.Equal(t, d.Tc.GetExpectsData(d.Env, "MACAddress").(string), data.Connections[0].MAC)
		}

	}
}
