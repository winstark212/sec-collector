package collector

var (
	// LocalIP Native active IP
	LocalIP string
	Config ClientConfig
)



type ClientConfig struct {
	Cycle  int    // 信息传输频率，单位：分钟
	UDP    bool   // 是否记录UDP请求
	LAN    bool   // 是否本地网络请求
	Mode   string // 模式，考虑中
	Filter struct {
		File    []string // 文件hash、文件名
		IP      []string // IP地址
		Process []string // 进程名、参数
	} // 直接过滤不回传的规则
	MonitorPath []string // 监控目录列表
	Lasttime    string   // 最后一条登录日志时间
}

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
