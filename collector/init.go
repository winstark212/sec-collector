package collector

var (
	// LocalIP Native active IP
	LocalIP string
)

// ComputerInfo computer information struct
type ComputerInfo struct {
	IP       string   // IP address
	System   string   // Operating system
	Hostname string   // Computer name
	Type     string   // Server type
	Path     []string // Web directory
}

// Startup startup service struct
type Startup  struct {
	Caption   string
	Command   string
	Location  string
	User      string
}

// Service list
type Service struct {
	Caption    string
	Name       string
	PathName   string
	Started    bool
	StartMode  string
	StartName  string
}

type utmp struct {
	UtType uint32
	UtPid  uint32    // PID of login process
	UtLine [32]byte  // device name of tty - "/dev/"
	UtID   [4]byte   // init id or abbrev. ttyname
	UtUser [32]byte  // user name
	UtHost [256]byte // hostname for remote login
	UtExit struct {
		ETermination uint16 // process termination status
		EExit        uint16 // process exit status
	}
	UtSession uint32 // Session ID, used for windowing
	UtTv      struct {
		TvSec  uint32 /* Seconds */
		TvUsec uint32 /* Microseconds */
	}
	UtAddrV6 [4]uint32 // IP address of remote host
	Unused   [20]byte  // Reserved for future use
}

// Collector interface
type Collector interface {
	GetComInfo() (info ComputerInfo)
	GetCrontab() []map[string]string
	GetListening() []map[string]string
	GetLoginLog() []map[string]string
	GetProcessList() []map[string]string
	GetServiceInfo() []map[string]string
	GetStartup() []map[string]string
	GetUser() []map[string]string
	
}