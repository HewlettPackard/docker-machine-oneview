package icsp

import (
	"encoding/json"
	"strings"

	"github.com/docker/machine/drivers/oneview/rest"
	"github.com/docker/machine/drivers/oneview/utils"
	"github.com/docker/machine/log"
)

// StorageDevice storage device type
type StorageDevice struct {
	Capacity   int    `json:"capacity,omitempty"`   // capacity Capacity of the storage in megabytes integer
	DeviceName string `json:"deviceName,omitempty"` // deviceName Device name, such as "C:" for Windows or "sda" for Linux string
	MediaType  string `json:"mediaType,omitempty"`  // mediaType Media type, such as "CDROM", "SCSI DISK" and etc. string
	Model      string `json:"model,omitempty"`      // model Model of the device string
	Vendor     string `json:"vendor,omitempty"`     // vendor Manufacturer of the device string
}

// Stage stage const
type Stage int

const (
	S_IN_DEPLOYMENT Stage = 1 + iota
	S_LIVE
	S_OFFLINE
	S_OPS_READY
	S_UNKNOWN
)

var stagelist = [...]string{
	"IN DEPLOYMENT", // - The managed Server is in process of deployment;
	"LIVE",          // - Deployment complete, the Server is live in production;
	"OFFLINE",       // - The managed Server is off-line;
	"OPS_READY",     // - The managed Server is available to operations;
	"UNKNOWN",       // - The managed Server is in an unknown stage - this is the default value for the field;
}

// String helper for stage
func (o Stage) String() string { return stagelist[o-1] }

// Equal helper for stage
func (o Stage) Equal(s string) bool { return (strings.ToUpper(s) == strings.ToUpper(o.String())) }

// ServerLocationItem server location type
type ServerLocationItem struct {
	Bay       string `json:"bay,omitempty"`       // bay Slot number in a rack where the Server is located string
	Enclosure string `json:"enclosure,omitempty"` // enclosure Name of an enclosure where the Server is physically located string
	Rack      string `json:"rack,omitempty"`      // rack Name of a rack where the Server is physically located string
}

// OpswLifecycle opsw lifecycle
type OpswLifecycle int

const (
	DEACTIVATED OpswLifecycle = 1 + iota
	MANAGED
	PROVISION_FAILED
	PROVISIONING
	UNPROVISIONED
	PRE_UNPROVISIONED
)

var opswlifecycle = [...]string{
	"DEACTIVATED",       // - No management activities can occur once a Server is deactivated;
	"MANAGED",           // - A production OS is installed and running on the target server. Normal management activities can occur when a Server is under management;
	"PROVISION_FAILED",  // - A managed Server enters this state when operating system installation or other provisioning activities failed;
	"PROVISIONING",      // - A managed Server is set to this state any time a job is running on the server;
	"UNPROVISIONED",     // - A managed Server in this state has booted into a service OS and is waiting to have an operating system installed;
	"PRE_UNPROVISIONED", // - A managed Server in this state is defined, but has not yet booted and registered with the appliance. An example of this is an iLO that was added without booting to maintenance;
}

// String helper for OpswLifecycle
func (o OpswLifecycle) String() string { return opswlifecycle[o-1] }

// Equal helper for OpswLifecycle
func (o OpswLifecycle) Equal(s string) bool {
	return (strings.ToUpper(s) == strings.ToUpper(o.String()))
}

// JobHistory job history type
type JobHistory struct {
	Description   string        `json:"description,omitempty"`   // description Description of the job, string
	EndDate       string        `json:"endDate,omitempty"`       // endDate Date and time when job was finished, string
	Initiator     string        `json:"initiator,omitempty"`     // initiator Name of the user who invoked the job on the Server, string
	Name          string        `json:"name,omitempty"`          // name Name of the job, string
	NameOfJobType string        `json:"nameOfJobType,omitempty"` // nameOfJobType Name of the OS Build Plan that was invoked on the Server, string
	StartDate     string        `json:"startDate,omitempty"`     // startDate Date and time when job was invoked, string
	URI           utils.Nstring `json:"uri,omitempty"`           // uri The canonical URI of the job, string
	URIOfJobType  utils.Nstring `json:"uriOfJobType,omitempty"`  // uriOfJobType The canonical URI of the OS Build Plan, string
}

// Interface struct
type Interface struct {
	DHCPEnabled bool   `json:"dhcpEnabled,omitempty"` // dhcpEnabled Flag that indicates whether the interface IP address is configured using DHCP, Boolean
	Duplex      string `json:"duplex,omitempty"`      // duplex Reported duplex of the interface, string
	IPV4Addr    string `json:"ipv4Addr,omitempty"`    // ipv4Addr IPv4 address of the interface, string
	IPV6Addr    string `json:"ipv6Addr,omitempty"`    // ipv6Addr IPv6 address of the interface, string
	MACAddr     string `json:"macAddr,omitempty"`     // macAddr Interface hardware network address, string
	Netmask     string `json:"netmask,omitempty"`     // netmask Netmask in dotted decimal notation, string
	Slot        string `json:"slot,omitempty"`        // slot Interface identity reported by the Server's operating system, string
	Speed       string `json:"speed,omitempty"`       // speed Interface speed in megabits, string
	Type        string `json:"type,omitempty"`        // type Interface type. For example, ETHERNET, string
}

// Ilo struct
type Ilo struct {
	Category       string        `json:"category,omitempty"`       // category The category is used to help identify the kind of resource, string
	Created        string        `json:"created,omitempty"`        // created Date and time when iLO was first discovered by Insight Control Server Provisioning, timestamp
	Description    string        `json:"description,omitempty"`    // description General description of the iLO, string
	ETAG           string        `json:"eTag,omitempty"`           // eTag Entity tag/version ID of the resource, the same value that is returned in the ETag header on a GET of the resource, string
	HealthStatus   string        `json:"healthStatus,omitempty"`   // healthStatus Overall health status of the resource, string
	IPAddress      string        `json:"ipAddress,omitempty"`      // ipAddress The IP address of the server’s iLO, string
	Modified       string        `json:"modified,omitempty"`       // modified Date and time when the resource was last modified, timestamp
	Name           string        `json:"name,omitempty"`           // name For servers added via iLO and booted to Intelligent Provisioning service OS, host name is determined by Intelligent Provisioning. For servers added via iLO and PXE booted to LinuxPE host name is "localhost". For servers added via iLO and PXE booted to WinPE host name is a random hostname "minint-xxx", string
	Passowrd       string        `json:"password,omitempty"`       // password ILO's password, string
	Port           int           `json:"port,omitempty"`           // port The socket on which the management service listens, integer
	ResourceStatus string        `json:"resourceStatus,omitempty"` // resourceStatus Current state of the resource, string
	Server         string        `json:"server,omitempty"`         // server The canonical URI of the hosting/managed server, string
	State          string        `json:"state,omitempty"`          // state Current state of the resource, string
	Status         string        `json:"status,omitempty"`         // status Overall health status of the resource, string
	Type           string        `json:"type,omitempty"`           // type Uniquely identifies the type of the JSON object(readonly), string
	URI            utils.Nstring `json:"uri,omitempty"`            // uri Unique numerical iLO identifier, string
	Username       string        `json:"username,omitempty"`       // username Username used to log in to iLO, string
}

// DeviceGroup struct
type DeviceGroup struct {
	Name  string        `json:"name,omitempty"`  // name Display name for the resource, string
	REFID int           `json:"refID,omitempty"` // refID The unique numerical identifier, integer
	URI   utils.Nstring `json:"uri,omitempty"`   // uri The canonical URI of the device group, string
}

// CPU struct
type CPU struct {
	CacheSize string `json:"cacheSize,omitempty"` // cacheSize CPU's cache size  , string
	Family    string `json:"family,omitempty"`    // family CPU's family. For example, "x86_64"  , string
	Model     string `json:"model,omitempty"`     // model CPU's model. For example, "Xeon"  , string
	Slot      string `json:"slot,omitempty"`      // slot CPU's slot  , string
	Speed     string `json:"speed,omitempty"`     // speed CPU's speed  , string
	Status    string `json:"status,omitempty"`    // status The last reported status of the CPU. For example, on-line, off-line  , string
	Stepping  string `json:"stepping,omitempty"`  // stepping CPU's stepping information  , string
}

// Server type
type Server struct {
	Architecture           string             `json:"architecture,omitempty"`           // architecture Server's architecture, string
	Category               string             `json:"category,omitempty"`               // category The category is used to help identify the kind of resource, string
	Cpus                   []CPU              `json:"cpus,omitempty"`                   // array of CPU's
	Created                string             `json:"created,omitempty"`                // created Date and time when the Server was discovered, timestamp
	CustomAttributes       []CustomAttribute  `json:"customAttributes,omitempty"`       // array of custom attributes
	DefaultGateway         string             `json:"defaultGateway,omitempty"`         // defaultGateway Gateway for this Server, string
	Description            string             `json:"description,omitempty"`            // description Brief description of the Server, string
	DeviceGroups           []DeviceGroup      `json:"deviceGroups,omitempty"`           // deviceGroups An array of device groups associated with the Server
	DiscoveredDate         string             `json:"discoveredDate,omitempty"`         // discoveredDate Date and time when the Server was discovered. Same as created date
	ETAG                   string             `json:"eTag,omitempty"`                   // eTag Entity tag/version ID of the resource
	Facility               string             `json:"facility,omitempty"`               // facility A facility represents the collection of servers. A facility can be all or part of a data center, Server room, or computer lab. Facilities are used as security boundaries with user groups
	HardwareModel          string             `json:"hardwareModel,omitempty"`          // hardwareModel The model name of the target server
	HostName               string             `json:"hostName,omitempty"`               // hostName The name of the server as reported by the server
	ILO                    Ilo                `json:"ilo,omitempty"`                    // information on ilo
	Interfaces             []Interface        `json:"interfaces,omitempty"`             // list of interfaces
	JobsHistory            []JobHistory       `json:"jobsHistory,omitempty"`            // array of previous run jobs
	LastScannedDate        string             `json:"lastScannedDate,omitempty"`        // lastScannedDate Date and time when the Server was detected last , string
	Locale                 string             `json:"locale,omitempty"`                 // locale Server's configured locale , string
	LoopbackIP             string             `json:"loopbackIP,omitempty"`             // loopbackIP Server's loopback IP address in dotted decimal format, string
	ManagementIP           string             `json:"managementIP,omitempty"`           // managementIP Server's management IP address in dotted decimal format, string
	Manufacturer           string             `json:"manufacturer,omitempty"`           // manufacturer Manufacturer as reported by the Server  , string
	MID                    string             `json:"mid,omitempty"`                    // mid A unique ID assigned to the Server by Server Automation, string
	Modified               string             `json:"modified,omitempty"`               // modified Date and time when the Server was last modified , timestamp
	Name                   string             `json:"name,omitempty"`                   // name The display name of the server. This is what shows on the left hand side of the UI. It is not the same as the host name. , string
	NetBios                string             `json:"netBios,omitempty"`                // netBios Server's Net BIOS name, string
	OperatingSystem        string             `json:"operatingSystem,omitempty"`        // operatingSystem Operating system installed on the Server, string
	OperatingSystemVersion string             `json:"operatingSystemVersion,omitempty"` // operatingSystemVersion Version of the operating system installed on the Server, string
	OpswLifecycle          string             `json:"opswLifcycle,omitempty"`           // Use type OpswLifcycle
	OSFlavor               string             `json:"osFlavor,omitempty"`               // osFlavor Additional information about an operating system flavor, string
	OSSPVersion            string             `json:"osSPVersion,omitempty"`            // osSPVersion Windows Service Pack version info, string
	PeerIP                 string             `json:"peerIP,omitempty"`                 // peerIP Server's peer IP address, string
	RAM                    string             `json:"ram,omitempty"`                    // ram Amount of free memory on the Server, string
	Reporting              bool               `json:"reporting,omitempty"`              // reporting Flag that indicates if the client on the Server is reporting to the core, Boolean
	Running                string             `json:"running,omitempty"`                // running Flag that indicates whether provisioning is performed on the Server, string
	SerialNumber           string             `json:"serialNumber,omitempty"`           // serialNumber The serial number assigned to the Server, string
	ServerLocation         ServerLocationItem `json:"serverLocation,omitempty"`         // serverLocation The Server location information such as rack and enclosure etc
	Stage                  string             `json:"stage,omitempty"`                  // stage type //stage When a managed Server is rolled out into production it typically passes to various stages of deployment.The following are the valid values for the stages of the Server:
	State                  string             `json:"state,omitempty"`                  // state Indicates the state of the agent on the target server. The following are the valid values for the state:
	Status                 string             `json:"status,omitempty"`                 // status Unified status of the target Server. Supported values:
	StorageDevices         []StorageDevice    `json:"storageDevices,omitempty"`         // storage devices on the server
	Swap                   string             `json:"swap,omitempty"`                   // swap Amount of swap space on the Server  , string
	Type                   string             `json:"type,omitempty"`                   // type Uniquely identifies the type of the JSON object  , string (readonly)
	URI                    utils.Nstring      `json:"uri,omitempty"`                    // uri The canonical URI of the Server  , string
	UUID                   string             `json:"uuid,omitempty"`                   // uuid Server's UUID  , string
}

// ServerList List of Servers
type ServerList struct {
	Category    string        `json:"category,omitempty"`    // Resource category used for authorizations and resource type groupings
	Count       int           `json:"count,omitempty"`       // The actual number of resources returned in the specified page
	Created     string        `json:"created,omitempty"`     // timestamp for when resource was created
	ETAG        string        `json:"eTag,omitempty"`        // entity tag version id
	Members     []Server      `json:"members,omitempty"`     // array of Server types
	Modified    string        `json:"modified,omitempty"`    // timestamp resource last modified
	NextPageURI utils.Nstring `json:"nextPageUri,omitempty"` // Next page resources
	PrevPageURI utils.Nstring `json:"prevPageUri,omitempty"` // Previous page resources
	Start       int           `json:"start,omitempty"`       // starting row of resource for current page
	Total       int           `json:"total,omitempty"`       // total number of pages
	Type        string        `json:"type,omitempty"`        // type of paging
	URI         utils.Nstring `json:"uri,omitempty"`         // uri to page
}

// ServerCreate structure for create server
type ServerCreate struct {
	IPAddress string `json:"ipAddress,omitempty"` // PXE managed ip address
	Port      int    `json:"port,omitempty"`      // port number to use
	UserName  string `json:"username,omitempty"`  // iLo username
	Password  string `json:"password,omitempty"`  // iLO password
}

// NewServerCreate make a new servercreate object
func (sc ServerCreate) NewServerCreate(user string, pass string, ip string, port int) ServerCreate {
	return ServerCreate{
		IPAddress: ip,
		Port:      port,
		UserName:  user,
		Password:  pass,
	}
}

// SubmitNewServer submit new profile template
func (c *ICSPClient) SubmitNewServer(sc ServerCreate) (jt *JobTask, err error) {
	log.Infof("Initializing creation of server for ICSP, %s.", sc.IPAddress)
	var (
		uri  = "/rest/os-deployment-servers"
		juri ODSUri
	)
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	jt = jt.NewJobTask(c)
	jt.Reset()
	data, err := c.RestAPICall(rest.POST, uri, sc)
	if err != nil {
		jt.IsDone = true
		log.Errorf("Error submitting new server request: %s", err)
		return jt, err
	}

	log.Debugf("Response submit new server %s", data)
	if err := json.Unmarshal([]byte(data), &juri); err != nil {
		jt.IsDone = true
		jt.JobURI = juri
		log.Errorf("Error with task un-marshal: %s", err)
		return jt, err
	}
	jt.JobURI = juri

	return jt, err
}

// CreateServer create profile from template
func (c *ICSPClient) CreateServer(user string, pass string, ip string, port int) error {

	var sc ServerCreate
	sc = sc.NewServerCreate(user, pass, ip, port)

	jt, err := c.SubmitNewServer(sc)
	err = jt.Wait()
	if err != nil {
		return err
	}
	return nil
}

// GetServers get a servers from icsp
func (c *ICSPClient) GetServers() (ServerList, error) {
	var (
		uri     = "/rest/os-deployment-servers"
		servers ServerList
	)

	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())
	data, err := c.RestAPICall(rest.GET, uri, nil)
	if err != nil {
		return servers, err
	}

	log.Debugf("GetServers %s", data)
	if err := json.Unmarshal([]byte(data), &servers); err != nil {
		return servers, err
	}
	return servers, nil
}

// DeleteServer deletes a server in icsp appliance instance
func (c *ICSPClient) DeleteServer(uuid string) error {
	var (
		uri = "/rest/os-deployment-servers"
	)

	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())
	//TODO: should check to make sure server uri has a real server
	//      and it's status is managed
	data, err := c.RestAPICall(rest.DELETE, uri, nil)
	//TODO : this should be returning a jobs uri
	if err != nil {
		return err
	}
	//TODO: implement wait for delete

	log.Debugf("DeleteServer %+v", data)
	return nil
}

// SaveServer save Server
// submit new profile template
func (c *ICSPClient) SaveServer(s Server) (o Server, err error) {
	log.Infof("Saving server attributes for %s.", s.Name)
	var (
		uri = s.URI
	)
	log.Debugf("REST : %s \n %+v\n", uri, s)
	data, err := c.RestAPICall(rest.PUT, uri.String(), s)
	if err != nil {
		log.Errorf("Error submitting new server request: %s", err)
		return o, err
	}
	if err := json.Unmarshal([]byte(data), &o); err != nil {
		return o, err
	}

	return o, err
}
