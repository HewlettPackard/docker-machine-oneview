package oneview

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/drivers/oneview/icsp"
	"github.com/docker/machine/drivers/oneview/ov"
	"github.com/docker/machine/log"
	"github.com/docker/machine/ssh"
	"github.com/docker/machine/state"
)

// Driver OneView driver structure
type Driver struct {
	*drivers.BaseDriver
	ClientICSP           *icsp.ICSPClient
	ClientOV             *ov.OVClient
	IloUser              string
	IloPassword          string
	IloPort              int
	OSBuildPlan          string
	SSHUser              string
	SSHPort              int
	SSHPublicKey         string
	ServerTemplate       string
	PublicSlotID         int
	PublicConnectionName string
	Profile              ov.ServerProfile
	Hardware             ov.ServerHardware
	Server               icsp.Server
}

const (
	defaultTimeout = 1 * time.Second
)

var (
	ErrDriverMissingEndPointOptionOV   = errors.New("Missing option --oneview-ov-endpoint or environment ONEVIEW_OV_ENDPOINT")
	ErrDriverMissingEndPointOptionICSP = errors.New("Missing option --oneview-icsp-endpoint or environment ONEVIEW_ICSP_ENDPOINT")
	ErrDriverMissingTemplateOption     = errors.New("Missing option --oneview-server-template or environment ONEVIEW_SERVER_TEMPLATE")
	ErrDriverMissingBuildPlanOption    = errors.New("Missing option --oneview-os-plan or ONEVIEW_OS_PLAN")
)

func init() {
	drivers.Register("oneview", &drivers.RegisteredDriver{
		New:            NewDriver,
		GetCreateFlags: GetCreateFlags,
	})
}

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
//
func GetCreateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "oneview-ov-user",
			Usage:  "User Name to OneView Server",
			Value:  "",
			EnvVar: "ONEVIEW_OV_USER",
		},
		cli.StringFlag{
			Name:   "oneview-ov-password",
			Usage:  "Password to OneView Server",
			Value:  "",
			EnvVar: "ONEVIEW_OV_PASSWORD",
		},
		cli.StringFlag{
			Name:   "oneview-ov-domain",
			Usage:  "Domain to OneView Server",
			Value:  "LOCAL",
			EnvVar: "ONEVIEW_OV_DOMAIN",
		},
		cli.StringFlag{
			Name:   "oneview-ov-endpoint",
			Usage:  "OneView Server URL Endpoint",
			Value:  "",
			EnvVar: "ONEVIEW_OV_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "oneview-icsp-user",
			Usage:  "User Name to OneView Insight Controller",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_USER",
		},
		cli.StringFlag{
			Name:   "oneview-icsp-password",
			Usage:  "Password to OneView Insight Controller",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "oneview-icsp-domain",
			Usage:  "Domain to OneView Insight Controller",
			Value:  "LOCAL",
			EnvVar: "ONEVIEW_ICSP_DOMAIN",
		},
		cli.StringFlag{
			Name:   "oneview-icsp-endpoint",
			Usage:  "OneView Insight Controller URL Endpoint",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "oneview-sslverify",
			Usage:  "SSH private key path",
			EnvVar: "ONEVIEW_SSLVERIFY",
		},
		cli.StringFlag{
			Name:   "oneview-ssh-user",
			Usage:  "OneView build plan ssh user account",
			Value:  "docker",
			EnvVar: "ONEVIEW_SSH_USER",
		},
		cli.IntFlag{
			Name:   "oneview-ssh-port",
			Usage:  "OneView build plan ssh host port",
			Value:  22,
			EnvVar: "ONEVIEW_SSH_PORT",
		},
		cli.StringFlag{
			Name:   "oneview-server-template",
			Usage:  "OneView server template to use for blade provisioning, see OneView Server Template for setup.",
			Value:  "DOCKER_1.8_OVTEMP",
			EnvVar: "ONEVIEW_SERVER_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "oneview-os-plan",
			Usage:  "OneView ICSP OS Build plan to use for OS provisioning, see ICS OS Plan for setup.",
			Value:  "RHEL71_DOCKER_1.8",
			EnvVar: "ONEVIEW_OS_PLAN",
		},
		cli.StringFlag{
			Name:   "oneview-ilo-user",
			Usage:  "ILO User id that is used during ICSP server creation.",
			Value:  "docker",
			EnvVar: "ONEVIEW_ILO_USER",
		},
		cli.StringFlag{
			Name:   "oneview-ilo-password",
			Usage:  "ILO password that is used during ICSP server creation.",
			Value:  "",
			EnvVar: "ONEVIEW_ILO_PASSWORD",
		},
		cli.IntFlag{
			Name:   "oneview-ilo-port",
			Usage:  "optional ILO port to use.",
			Value:  443,
			EnvVar: "ONEVIEW_ILO_PORT",
		},
		cli.IntFlag{
			Name:   "oneview-public-slotid",
			Usage:  "Optional slot id of the public interface to use for connecting with docker.",
			Value:  1,
			EnvVar: "ONEVIEW_PUBLIC_SLOTID",
		},
		cli.StringFlag{
			Name:   "oneview-public-connection-name",
			Usage:  "Optional Connection name to use for public interface to connect with docker, the name can be defined in server template.  Overrides slotid option.",
			Value:  "",
			EnvVar: "ONEVIEW_PUBLIC_CONNECTION_NAME",
		},
	}
}

// NewDriver - create a OneView object driver
func NewDriver(machineName string, storePath string, caCert string, privateKey string) (drivers.Driver, error) {
	inner := drivers.NewBaseDriver(machineName, storePath, caCert, privateKey)
	return &Driver{BaseDriver: inner}, nil
}

// DriverName - get the name of the driver
func (d *Driver) DriverName() string {
	log.Debug("DriverName...")
	return "oneview"
}

// GetSSHHostname - gets the hostname that docker-machine connects to
func (d *Driver) GetSSHHostname() (string, error) {
	log.Debug("GetSSHHostname...")
	return d.GetIP()
}

// GetSSHUsername - gets the ssh user that will be connected to
func (d *Driver) GetSSHUsername() string {
	log.Debug("GetSSHUsername...")
	return d.SSHUser
}

// SetConfigFromFlags - gets the cli configuration flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	log.Debug("SetConfigFromFlags...")

	d.ClientICSP = d.ClientICSP.NewICSPClient(flags.String("oneview-icsp-user"),
		flags.String("oneview-icsp-password"),
		flags.String("oneview-icsp-domain"),
		flags.String("oneview-icsp-endpoint"),
		flags.Bool("oneview-sslverify"),
		1)

	d.ClientOV = d.ClientOV.NewOVClient(flags.String("oneview-ov-user"),
		flags.String("oneview-ov-password"),
		flags.String("oneview-ov-domain"),
		flags.String("oneview-ov-endpoint"),
		flags.Bool("oneview-sslverify"),
		1)

	d.ClientICSP.RefreshVersion()
	d.ClientOV.RefreshVersion()

	d.IloUser = flags.String("oneview-ilo-user")
	d.IloPassword = flags.String("oneview-ilo-password")
	d.IloPort = flags.Int("oneview-ilo-port")

	d.PublicSlotID = flags.Int("oneview-public-slotid")
	d.PublicConnectionName = flags.String("oneview-public-connection-name")

	d.SSHUser = flags.String("oneview-ssh-user")
	d.SSHPort = flags.Int("oneview-ssh-port")

	d.ServerTemplate = flags.String("oneview-server-template")
	d.OSBuildPlan = flags.String("oneview-os-plan")

	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")

	// TODO : we should verify settings for each client

	// check for the ov endpoint
	if d.ClientOV.Endpoint == "" {
		return ErrDriverMissingEndPointOptionOV
	}
	// check for the icsp endpoint
	if d.ClientICSP.Endpoint == "" {
		return ErrDriverMissingEndPointOptionICSP
	}
	// check for the template name
	if d.ServerTemplate == "" {
		return ErrDriverMissingTemplateOption
	}

	if d.OSBuildPlan == "" {
		return ErrDriverMissingBuildPlanOption
	}

	return nil
}

// PreCreateCheck - pre create check
func (d *Driver) PreCreateCheck() (err error) {
	log.Debug("PreCreateCheck...")
	// verify you can connect to ov
	ovVersion, err := d.ClientOV.GetAPIVersion()
	if err != nil {
		return err
	}
	if ovVersion.CurrentVersion <= 0 {
		return fmt.Errorf("Unable to get a valid version from OneView,  %+v\n", ovVersion)
	}
	// verify you can connect to icsp
	icspVersion, err := d.ClientICSP.GetAPIVersion()
	if err != nil {
		return err
	}
	if icspVersion.CurrentVersion <= 0 {
		return fmt.Errorf("Unable to get a valid version from ICsp,  %+v\n", icspVersion)
	}
	return nil
}

// Create - create server for docker
func (d *Driver) Create() error {
	log.Infof("Generating SSH keys...")
	if err := d.createKeyPair(); err != nil {
		return fmt.Errorf("unable to create key pair: %s", err)
	}

	log.Debugf("ICSP Endpoint is: %s", d.ClientICSP.Endpoint)
	log.Debugf("OV Endpoint is: %s", d.ClientOV.Endpoint)
	// create the server profile in oneview, we need a hostname and a template name

	log.Debugf("***> CreateMachine")
	// create d.Hardware and d.Profile
	if err := d.ClientOV.CreateMachine(d.MachineName, d.ServerTemplate); err != nil {
		return err
	}

	if err := d.getBlade(); err != nil {
		return err
	}

	// power on the server, and leave it in that state
	if err := d.Hardware.PowerOn(); err != nil {
		return err
	}
	// power off let customization bring the server online
	if err := d.Hardware.PowerOff(); err != nil {
		return err
	}

	// add the server to icsp, TestCreateServer
	// apply a build plan, TestApplyDeploymentJobs
	var sp *icsp.CustomServerAttributes
	sp = sp.New()
	sp.Set("docker_user", d.SSHUser)
	sp.Set("public_key", d.SSHPublicKey)
	// TODO: make a util for this
	if len(os.Getenv("proxy_enable")) > 0 {
		sp.Set("proxy_enable", os.Getenv("proxy_enable"))
	} else {
		sp.Set("proxy_enable", "false")
	}

	strProxy := os.Getenv("proxy_config")
	sp.Set("proxy_config", strProxy)

	sp.Set("docker_hostname", d.MachineName+"-@server_name@")

	sp.Set("interface", "@interface@") // this is populated later

	// Get the mac address for public Connection on server profile
	var publicmac string
	if d.PublicConnectionName != "" {
		conn, err := d.Profile.GetConnectionByName(d.PublicConnectionName)
		if err != nil {
			return err
		}
		publicmac = conn.MAC.String()
	} else {
		publicmac = ""
	}

	// arguments for customize server
	cs := icsp.CustomizeServer{
		HostName:         d.MachineName,                   // machine-rack-enclosure-bay
		SerialNumber:     d.Profile.SerialNumber.String(), // get it
		ILoUser:          d.IloUser,
		IloPassword:      d.IloPassword,
		IloIPAddress:     d.Hardware.GetIloIPAddress(), // MpIpAddress for v1
		IloPort:          d.IloPort,
		OSBuildPlan:      d.OSBuildPlan,  // name of the OS build plan
		PublicSlotID:     d.PublicSlotID, // this is the slot id of the public interface
		PublicMAC:        publicmac,      // Server profile mac address, overrides slotid
		ServerProperties: sp,
	}
	// create d.Server and apply a build plan and configure the custom attributes
	if err := d.ClientICSP.CustomizeServer(cs); err != nil {
		return err
	}

	ip, err := d.GetIP()
	if err != nil {
		return err
	}
	d.IPAddress = ip

	// use ssh to set keys, and test ssh
	sshClient, err := d.getLocalSSHClient()
	if err != nil {
		return err
	}

	pubKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return err
	}

	if out, err := sshClient.Output(fmt.Sprintf(
		"printf '%%s' '%s' | tee /home/%s/.ssh/authorized_keys",
		string(pubKey),
		d.GetSSHUsername(),
	)); err != nil {
		log.Error(out)
		return err
	}
	log.Infof("%s, Completed all create steps, docker provisioning will continue.", d.DriverName())

	defer closeAll(d)
	return nil
}

// closeAll - cleanup sessions on the OV and ICSP appliances
func closeAll(d *Driver) {
	err := d.ClientOV.SessionLogout()
	if err != nil {
		log.Warnf("OV Session Logout : %s", err)
	}
	err = d.ClientICSP.SessionLogout()
	if err != nil {
		log.Warnf("ICSP Session Logout : %s", err)
	}
}

// GetURL - get docker url
func (d *Driver) GetURL() (string, error) {
	log.Debug("GetURL...")
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

// GetIP - get server host or ip address
// TODO: we need to get ip of server from icsp or ov??
// currently the only way i can see to get this is with sudo ifconfig|grep inet
func (d *Driver) GetIP() (string, error) {
	log.Debug("GetIP...")
	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return "", err
	}
	sPublicIPv4, err := d.Server.GetPublicIPV4()
	if err != nil {
		return "", err
	}
	if sPublicIPv4 == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return sPublicIPv4, nil
}

// GetState - get the running state of the target machine
func (d *Driver) GetState() (state.State, error) {
	log.Debug("GetState...")

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return state.Error, err
	}
	if icsp.PROVISIONING.Equal(d.Server.OpswLifecycle) {
		return state.Starting, nil
	}
	if icsp.UNPROVISIONED.Equal(d.Server.OpswLifecycle) ||
		icsp.PRE_UNPROVISIONED.Equal(d.Server.OpswLifecycle) {
		return state.Stopping, nil
	}
	if icsp.DEACTIVATED.Equal(d.Server.OpswLifecycle) {
		return state.Stopped, nil
	}
	if icsp.PROVISION_FAILED.Equal(d.Server.OpswLifecycle) {
		return state.Error, nil
	}
	// use power state to determine status
	ps, err := d.Hardware.GetPowerState()
	if err != nil {
		return state.Error, err
	}
	switch ps {
	case ov.P_ON:
		return state.Running, nil
	case ov.P_OFF:
		return state.Stopped, nil
	case ov.P_UKNOWN:
		return state.Error, nil
	default:
		return state.None, nil
	}

}

// Start - start the docker machine target
func (d *Driver) Start() error {
	log.Infof("Starting ... %s", d.MachineName)

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return err
	}

	// power on the server, and leave it in that state
	if err := d.Hardware.PowerOn(); err != nil {
		return err
	}
	// implement icsp check for is in maintenance mode or started
	isManaged, err := d.ClientICSP.IsServerManaged(d.Hardware.SerialNumber.String())
	if err != nil {
		return err
	}
	if !isManaged {
		return errors.New("Server was started but not ready, check icsp status")
	}
	return nil
}

// Stop - stop the docker machine target
func (d *Driver) Stop() error {
	log.Debug("Stop...")
	log.Infof("Stop ... %s", d.MachineName)
	// gracefully attempt to stop the os

	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo shutdown -P now"); err != nil {
		log.Warnf("Problem shutting down gracefully : %s", err)
	}

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return err
	}

	// power on the server, and leave it in that state
	if err := d.Hardware.PowerOff(); err != nil {
		return err
	}
	// cleanup
	defer closeAll(d)
	return nil
}

// Remove - remove the docker machine target
//    Should remove the ICSP provisioned plan and the Server Profile from OV
func (d *Driver) Remove() error {
	log.Debug("Remove...")
	// remove the ssh keys
	if err := d.deleteKeyPair(); err != nil {
		return err
	}
	if err := d.Stop(); err != nil {
		return err
	}
	if err := d.getBlade(); err != nil {
		return err
	}
	// destroy the server in icsp
	isDeleted, err := d.ClientICSP.DeleteServer(d.Server.MID)
	if err != nil {
		return err
	}
	if !isDeleted {
		return fmt.Errorf("Unable to delete the server from icsp : %s, %s", d.MachineName, d.Server.MID)
	}
	// delete the server profile in ov : TestDeleteProfile
	t, err := d.ClientOV.SubmitDeleteProfile(d.Profile)
	err = t.Wait()
	if err != nil {
		return err
	}
	// cleanup
	defer closeAll(d)
	return nil
}

// Restart - restart the target machine
func (d *Driver) Restart() error {
	log.Debug("Restarting...")
	if err := d.Stop(); err != nil {
		return err
	}
	return d.Start()
}

// Kill - kill the docker machine
func (d *Driver) Kill() error {
	log.Debug("Killing...")
	//TODO: implement power off , is there a force?
	return nil
}

// publicSSHKeyPath - get the path to public ssh key
func (d *Driver) publicSSHKeyPath() string {
	log.Debug("publicSSHKeyPath...")
	return d.GetSSHKeyPath() + ".pub"
}

// /////////  HELPLERS /////////////

func (d *Driver) getBlade() (err error) {
	log.Debug("In getBlade()")

	d.Profile, err = d.ClientOV.GetProfileByName(d.MachineName)
	if err != nil {
		return err
	}

	log.Debugf("***> check if we got a profile")
	if d.Profile.URI.IsNil() {
		err = fmt.Errorf("Attempting to get machine profile information, unable to find machine in oneview: %s", d.MachineName)
		return err
	}

	// power on the server
	// get the server hardware associated with that test profile
	log.Debugf("***> GetServerHardware")
	d.Hardware, err = d.ClientOV.GetServerHardware(d.Profile.ServerHardwareURI)
	if d.Hardware.URI.IsNil() {
		err = fmt.Errorf("Attempting to get machine blade information, unable to find machine: %s", d.MachineName)
		return err
	}
	// get an icsp server
	if d.Hardware.VirtualSerialNumber.IsNil() {
		// get the server profile with SerialNumber
		d.Server, err = d.ClientICSP.GetServerBySerialNumber(d.Hardware.SerialNumber.String())
	} else {
		// get the server profile with the VirtualSerialNumber
		d.Server, err = d.ClientICSP.GetServerBySerialNumber(d.Hardware.VirtualSerialNumber.String())
	}
	if err != nil {
		return err
	}
	return err
}

// createKeyPair - generate key files needed
func (d *Driver) createKeyPair() error {

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	log.Debugf("created keys => %s", string(publicKey))
	d.SSHPublicKey = string(publicKey)
	return nil
}

// deleteKeyPair
func (d *Driver) deleteKeyPair() error {
	if err := os.Remove(d.GetSSHKeyPath()); err != nil {
		return err
	}
	if err := os.Remove(d.GetSSHKeyPath() + ".pub"); err != nil {
		return err
	}
	return nil
}

func (d *Driver) getLocalSSHClient() (ssh.Client, error) {
	sshAuth := &ssh.Auth{
		Passwords: []string{"docker"},
		Keys:      []string{d.GetSSHKeyPath()},
	}
	sshClient, err := ssh.NewNativeClient(d.GetSSHUsername(), d.IPAddress, d.SSHPort, sshAuth)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
