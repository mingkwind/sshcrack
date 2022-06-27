package vars

import (
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/patrickmn/go-cache"
)

type Service struct {
	Ip       string
	Port     int
	Protocol string
	Username string
	Password string
}

var (
	Timeout          = 3 * time.Second
	SupportProtocols map[string]bool

	// 检测端口是否开放的进度条
	ProcessBarActive *pb.ProgressBar
	// 弱口令扫描进度条
	ProgressBar *pb.ProgressBar

	// 默认协程数
	ThreadNum = 50

	DebugMode bool

	// 标记特定服务的特定用户是否破解成功，成功的话不再尝试破解该用户
	SuccessHash sync.Map

	// 扫描结果保存到一个cache中，该cache库支持内存数据落盘
	// Cache 是一个线程安全的Map
	CacheService *cache.Cache

	// 未能正确进行爆破而进行二次扫描的服务组
	TmpService *cache.Cache

	// 开始时间
	StartTime time.Time

	// ip_list，用户名与密码字典
	IpList   = "ip_list.txt"
	UserDict = "user.txt"
	PassDict = "pass.txt"

	// 结果保存文件
	ResultFile = "password_crack.txt"

	// 对于包含以下字符串的错误的ssh连接，进行重新破解
	ReScanStrs = []string{}
)

func init() {
	CacheService = cache.New(cache.NoExpiration, cache.DefaultExpiration)
	TmpService = cache.New(cache.NoExpiration, cache.DefaultExpiration)
	SupportProtocols = make(map[string]bool)
}
