package config

// 软件信息。
const (
	VERSION   string = "v1.0.0"                                      // 软件版本号
	AUTHOR    string = "AzuresYang"                                  // 软件作者
	NAME      string = "ARX 网络爬虫"                                    // 软件名
	FULL_NAME string = NAME + "_" + VERSION + " （by " + AUTHOR + "）" // 软件全称
	TAG       string = "arx"                                         // 软件标识符
)

// 默认配置。
const (
	WORK_ROOT string = TAG + "_pkg"                        // 运行时的目录名称
	CONFIG    string = WORK_ROOT + "/config.json"          // 配置文件路径
	CACHE_DIR string = WORK_ROOT + "/cache"                // 缓存文件目录
	LOG       string = WORK_ROOT + "/logs/arx_crawler.log" // 日志文件路径
	// LOG_ASYNC      bool   = true                            // 是否异步输出日志
	PHANTOMJS_TEMP string = CACHE_DIR // Surfer-Phantom下载器：js文件临时目录
)
