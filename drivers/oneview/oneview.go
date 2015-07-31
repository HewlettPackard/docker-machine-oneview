package oneview

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/log"
	"github.com/docker/machine/state"
	"github.com/docker/machine/utils"
)

type Driver struct {
	*drivers.BaseDriver
	SSHKey string
}

const (
	defaultTimeout = 1 * time.Second
)

func init() {
	drivers.Register("oneview", &drivers.RegisteredDriver{
		New:            NewDriver,
		GetCreateFlags: GetCreateFlags,
	})
}

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
func GetCreateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "generic-ip-address",
			Usage: "IP Address of machine",
			Value: "",
		},
		cli.StringFlag{
			Name:  "generic-ssh-user",
			Usage: "SSH user",
			Value: "root",
		},
		cli.StringFlag{
			Name:  "generic-ssh-key",
			Usage: "SSH private key path",
			Value: filepath.Join(utils.GetHomeDir(), ".ssh", "id_rsa"),
		},
		cli.IntFlag{
			Name:  "generic-ssh-port",
			Usage: "SSH port",
			Value: 22,
		},
	}
}

func NewDriver(machineName string, storePath string, caCert string, privateKey string) (drivers.Driver, error) {
	inner := drivers.NewBaseDriver(machineName, storePath, caCert, privateKey)
	return &Driver{BaseDriver: inner}, nil
}

func (d *Driver) DriverName() string {
	return "oneview"
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHUsername() string {
	return d.SSHUser
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIUrl = flags.String("oneview-api-url")
	d.SSHUser = flags.String("oneview-ssh-user")
	d.SSHKey = flags.String("oneview-ssh-key")
	d.SSHPort = flags.Int("oneview-ssh-port")

	if d.APIUrl == "" {
		return fmt.Errorf("oneview driver requires the --oneview-api-url option")
	}

	if d.SSHKey == "" {
		return fmt.Errorf("generic driver requires the --oneview-ssh-key option")
	}

	return nil
}

func (d *Driver) PreCreateCheck() error {
	return nil
}

func (d *Driver) Create() error {
	log.Infof("Importing SSH key...")

	if err := utils.CopyFile(d.SSHKey, d.GetSSHKeyPath()); err != nil {
		return fmt.Errorf("unable to copy ssh key: %s", err)
	}

	if err := os.Chmod(d.GetSSHKeyPath(), 0600); err != nil {
		return err
	}

	log.Debugf("IP: %s", d.IPAddress)

	return nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

func (d *Driver) GetState() (state.State, error) {
	addr := fmt.Sprintf("%s:%d", d.IPAddress, d.SSHPort)
	_, err := net.DialTimeout("tcp", addr, defaultTimeout)
	var st state.State
	if err != nil {
		st = state.Stopped
	} else {
		st = state.Running
	}
	return st, nil
}

func (d *Driver) Start() error {
	return fmt.Errorf("oneview driver does not support start")
}

func (d *Driver) Stop() error {
	return fmt.Errorf("oneview driver does not support stop")
}

func (d *Driver) Remove() error {
	return nil
}

func (d *Driver) Restart() error {
	log.Debug("Restarting...")

	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo shutdown -r now"); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Kill() error {
	log.Debug("Killing...")

	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo shutdown -P now"); err != nil {
		return err
	}

	return nil
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}
