package metrics

import (
	"strconv"
	"sync"

	"github.com/orangeAppsRu/custom-exporter/pkg/config"
	"github.com/orangeAppsRu/custom-exporter/pkg/filehash"
	"github.com/orangeAppsRu/custom-exporter/pkg/network"
	"github.com/orangeAppsRu/custom-exporter/pkg/hetzner"
	"github.com/orangeAppsRu/custom-exporter/pkg/hetznercloud"
	"github.com/orangeAppsRu/custom-exporter/pkg/yandex"
	"github.com/orangeAppsRu/custom-exporter/pkg/aws"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	hashGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file_hash",
			Help: "SHA256 hash of files",
		},
		[]string{"file"},
	)

	networkTargetGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_target",
			Help: "Network port availability",
		},
		[]string{"host", "port", "protocol"},
	)

	processCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "processes_count",
			Help: "Number of running processes",
		},
		[]string{"type"},
	)

	processRunningStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_running_status",
			Help: "Status of process (running or not)",
		},
		[]string{"process"},
	)

	processCpuTimeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "process_cpu_time_total",
			Help: "CPU time consumed by processes",
		},
		[]string{"process"},
	)

	processMemoryResidentGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_resident",
			Help: "Resident memory of processes in bytes",
		},
		[]string{"process"},
	)

	hostnameChecksumGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "hostname_checksum",
			Help: "Checksum of hostname",
		},
	)
	
	hostnameGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hostname",
			Help: "Hostname of the machine",
		},
		[]string{"hostname"},
	)

	unameChecksumGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "uname_checksum",
			Help: "Checksum of uname",
		},
	)

	uptimeSecondsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "uptime_seconds",
			Help: "Uptime of the machine in seconds",
		},
	)

	countLoginUsersGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "login_users_count",
			Help: "Number of login users",
		},
	)

	puppetCatalogLastCompileTimestampGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puppet_catalog_last_compile_timestamp",
			Help: "Timestamp of the last puppet catalog compile",
		},
	)

	puppetCatalogLastCompileStatusGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puppet_catalog_last_compile_status",
			Help: "Status of the last puppet catalog compile",
		},
	)

	hetznerServersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hetzner_robot_server",
			Help: "Hetzner robot server",
		},
		[]string{"id", "name", "type", "zone", "region", "ip"},
	)

	hetznerCloudServersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hetzner_cloud_server",
			Help: "Hetzner cloud server",
		},
		[]string{"id", "name", "type", "zone", "region", "ip"},
	)

	yandexCloudServersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "yandex_cloud_server",
			Help: "Yandex cloud server",
		},
		[]string{"id", "name", "type", "zone", "region", "public_ip", "private_ip", "cpu_count", "memory"},
	)

	awsCloudServersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_cloud_server",
			Help: "AWS cloud server",
		},
		[]string{"id", "name", "type", "zone", "region", "public_ip", "private_ip", "private_dns_name"},
	)


	previousHostnameLabel string

	fileHashMutex sync.Mutex
	networkTargetMutex sync.Mutex
	processCountMutex sync.Mutex
	processCpuTimeMutex sync.Mutex
	processMemoryResidentMutex sync.Mutex
	processRunningStatusMutex sync.Mutex
	hostnameChecksumMutex sync.Mutex
	hostnameMutex sync.Mutex
	unameChecksumMutex sync.Mutex
	uptimeSecondsMutex sync.Mutex
	countLoginUsersMutex sync.Mutex
	puppetCatalogLastCompileTimestampMutex sync.Mutex
	puppetCatalogLastCompileStatusMutex sync.Mutex
	hetznerServersMutex sync.Mutex
	hetznerCloudServersMutex sync.Mutex
	yandexCloudServersMutex sync.Mutex
	awsCloudServersMutex sync.Mutex
)

func RegistrMetrics(cfg config.Config) {
    prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
    prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())

	if cfg.FileHashCollector.Enabled {
		prometheus.MustRegister(hashGauge)
	}

	if cfg.PortCollector.Enabled {
		prometheus.MustRegister(networkTargetGauge)
	}

	if cfg.ProcessCollector.Enabled {
		prometheus.MustRegister(processCountGauge)
		prometheus.MustRegister(processCpuTimeCounter)
		prometheus.MustRegister(processMemoryResidentGauge)
		prometheus.MustRegister(processRunningStatusGauge)
	}

	if cfg.SystemCollector.Enabled {
		prometheus.MustRegister(hostnameChecksumGauge)
		prometheus.MustRegister(hostnameGauge)
		prometheus.MustRegister(unameChecksumGauge)
		prometheus.MustRegister(uptimeSecondsCounter)
		prometheus.MustRegister(countLoginUsersGauge)
	}

	if cfg.PuppetCollector.Enabled {
		prometheus.MustRegister(puppetCatalogLastCompileTimestampGauge)
		prometheus.MustRegister(puppetCatalogLastCompileStatusGauge)
	}

	if cfg.HetznerCollector.Enabled {
		prometheus.MustRegister(hetznerServersGauge)
	}

	if cfg.HetznerCloudCollector.Enabled {
		prometheus.MustRegister(hetznerCloudServersGauge)
	}

	if cfg.YandexCloudCollector.Enabled {
		prometheus.MustRegister(yandexCloudServersGauge)
	}

	if cfg.AWSCloudCollector.Enabled {
		prometheus.MustRegister(awsCloudServersGauge)
	}
}

func UpdateFileHashMetrics(filesWithHash []filehash.FileHash) {
	for _, fileInfo := range filesWithHash {
		fileHashMutex.Lock()
		hashGauge.WithLabelValues(fileInfo.File).Set(fileInfo.Hash)
		fileHashMutex.Unlock()
	}
}

func UpdateNetworkTargetsMetrics(targets []network.ResultTarget) {
	for _, t := range targets {
		value := 0
		if t.IsOpen {
			value = 1
		}
		networkTargetMutex.Lock()
		networkTargetGauge.With(prometheus.Labels{
			"host": t.Host,
			"port": strconv.Itoa(int(t.Port)),
			"protocol": t.Protocol,
		}).Set(float64(value))
		networkTargetMutex.Unlock()
	}
}

func UpdateProcessCountMetrics(typeProcess string, count int) {
	processCountMutex.Lock()
	processCountGauge.WithLabelValues(typeProcess).Set(float64(count))
	processCountMutex.Unlock()
}

func UpdateProcessCPUTimeMetrics(process string, time float64) {
	processCpuTimeMutex.Lock()
	processCpuTimeCounter.WithLabelValues(process).Add(time)
	processCpuTimeMutex.Unlock()
}

func UpdateProcessMemoryResidentMetrics(process string, memory uint64) {
	processMemoryResidentMutex.Lock()
	processMemoryResidentGauge.WithLabelValues(process).Set(float64(memory))
	processMemoryResidentMutex.Unlock()
}

func UpdateProcessRunningStatusMetrics(process string, status int) {
	processRunningStatusMutex.Lock()
	processRunningStatusGauge.WithLabelValues(process).Set(float64(status))
	processRunningStatusMutex.Unlock()
}

func UpdateHostnameChecksumMetrics(checksum float64) {
	hostnameChecksumMutex.Lock()
	hostnameChecksumGauge.Set(checksum)
	hostnameChecksumMutex.Unlock()
}

func UpdateHostnameMetrics(hostname string) {
	hostnameMutex.Lock()
	if previousHostnameLabel != "" && previousHostnameLabel != hostname {
		hostnameGauge.DeleteLabelValues(previousHostnameLabel)
	}
	hostnameGauge.WithLabelValues(hostname).Set(1)
	previousHostnameLabel = hostname
	hostnameMutex.Unlock()
}

func UpdateUnameChecksumMetrics(checksum float64) {
	unameChecksumMutex.Lock()
	unameChecksumGauge.Set(checksum)
	unameChecksumMutex.Unlock()
}

func UpdateUptimeSecondsMetrics(uptime float64) {
	uptimeSecondsMutex.Lock()
	uptimeSecondsCounter.Add(uptime)
	uptimeSecondsMutex.Unlock()
}

func UpdateLoginUsersCountMetrics(count int) {
	countLoginUsersMutex.Lock()
	countLoginUsersGauge.Set(float64(count))
	countLoginUsersMutex.Unlock()
}

func UpdatePuppetCatalogLastCompileTimestampMetrics(timestamp int64) {
	puppetCatalogLastCompileTimestampMutex.Lock()
	puppetCatalogLastCompileTimestampGauge.Set(float64(timestamp))
	puppetCatalogLastCompileTimestampMutex.Unlock()
}

func UpdatePuppetCatalogLastCompileStatusMetrics(status bool) {
	statusValue := 0
	if status {
		statusValue = 1
	}
	puppetCatalogLastCompileStatusMutex.Lock()
	puppetCatalogLastCompileStatusGauge.Set(float64(statusValue))
	puppetCatalogLastCompileStatusMutex.Unlock()
}

func UpdateHetznerServersMetrics(servers []hetzner.HrobotServer) {
	for _, s := range servers {
		hetznerServersMutex.Lock()
		hetznerServersGauge.With(prometheus.Labels{
			"id": strconv.FormatInt(s.ID, 10),
			"name": s.Name,
			"type": s.Type,
			"zone": s.Zone,
			"region": s.Region,
			"ip": s.IP.String(),
		}).Set(1)
		hetznerServersMutex.Unlock()
	}
}

func UpdateHetznerCloudServersMetrics(servers []hetznercloud.Server) {
	for _, s := range servers {
		hetznerCloudServersMutex.Lock()
		hetznerCloudServersGauge.With(prometheus.Labels{
			"id": strconv.FormatInt(s.ID, 10),
			"name": s.Name,
			"type": s.Type,
			"zone": s.Zone,
			"region": s.Region,
			"ip": s.IP.String(),
		}).Set(1)
		hetznerCloudServersMutex.Unlock()
	}
}

func UpdateYandexCloudServersMetrics(servers []yandex.Server) {
	for _, s := range servers {
		yandexCloudServersMutex.Lock()
		publicIP := s.PublicIP.String()
		if publicIP == "<nil>" {
			publicIP = ""
		}
		privateIP := s.PrivateIP.String()
		if privateIP == "<nil>" {
			privateIP = ""
		}
		yandexCloudServersGauge.With(prometheus.Labels{
			"id": s.ID,
			"name": s.Name,
			"type": s.Type,
			"zone": s.Zone,
			"region": s.Region,
			"public_ip": publicIP,
			"private_ip": privateIP,
			"cpu_count": strconv.FormatUint(uint64(s.CpuCount), 10),
			"memory": strconv.FormatUint(s.Memory, 10),
		}).Set(1)
		yandexCloudServersMutex.Unlock()
	}
}

func UpdateAWSCloudServersMetrics(servers []aws.Server) {
	for _, s := range servers {
		awsCloudServersMutex.Lock()
		publicIP := s.PublicIP.String()
		if publicIP == "<nil>" {
			publicIP = ""
		}
		privateIP := s.PrivateIP.String()
		if privateIP == "<nil>" {
			privateIP = ""
		}
		awsCloudServersGauge.With(prometheus.Labels{
			"id": s.ID,
			"name": s.Name,
			"private_dns_name": s.PrivateDnsName,
			"type": s.Type,
			"zone": s.Zone,
			"region": s.Region,
			"public_ip": publicIP,
			"private_ip": privateIP,
		}).Set(1)
		awsCloudServersMutex.Unlock()
	}
}

