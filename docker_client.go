package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

type Stats struct {
	Read      time.Time `json:"read"`
	Preread   time.Time `json:"preread"`
	PidsStats struct {
		Current int `json:"current"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []interface{} `json:"io_service_bytes_recursive"`
		IoServicedRecursive     []interface{} `json:"io_serviced_recursive"`
		IoQueueRecursive        []interface{} `json:"io_queue_recursive"`
		IoServiceTimeRecursive  []interface{} `json:"io_service_time_recursive"`
		IoWaitTimeRecursive     []interface{} `json:"io_wait_time_recursive"`
		IoMergedRecursive       []interface{} `json:"io_merged_recursive"`
		IoTimeRecursive         []interface{} `json:"io_time_recursive"`
		SectorsRecursive        []interface{} `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			PercpuUsage       []int `json:"percpu_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			PercpuUsage       []int `json:"percpu_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    int `json:"usage"`
		MaxUsage int `json:"max_usage"`
		Stats    struct {
			ActiveAnon              int   `json:"active_anon"`
			ActiveFile              int   `json:"active_file"`
			Cache                   int   `json:"cache"`
			Dirty                   int   `json:"dirty"`
			HierarchicalMemoryLimit int64 `json:"hierarchical_memory_limit"`
			HierarchicalMemswLimit  int64 `json:"hierarchical_memsw_limit"`
			InactiveAnon            int   `json:"inactive_anon"`
			InactiveFile            int   `json:"inactive_file"`
			MappedFile              int   `json:"mapped_file"`
			Pgfault                 int   `json:"pgfault"`
			Pgmajfault              int   `json:"pgmajfault"`
			Pgpgin                  int   `json:"pgpgin"`
			Pgpgout                 int   `json:"pgpgout"`
			Rss                     int   `json:"rss"`
			RssHuge                 int   `json:"rss_huge"`
			Swap                    int   `json:"swap"`
			TotalActiveAnon         int   `json:"total_active_anon"`
			TotalActiveFile         int   `json:"total_active_file"`
			TotalCache              int   `json:"total_cache"`
			TotalDirty              int   `json:"total_dirty"`
			TotalInactiveAnon       int   `json:"total_inactive_anon"`
			TotalInactiveFile       int   `json:"total_inactive_file"`
			TotalMappedFile         int   `json:"total_mapped_file"`
			TotalPgfault            int   `json:"total_pgfault"`
			TotalPgmajfault         int   `json:"total_pgmajfault"`
			TotalPgpgin             int   `json:"total_pgpgin"`
			TotalPgpgout            int   `json:"total_pgpgout"`
			TotalRss                int   `json:"total_rss"`
			TotalRssHuge            int   `json:"total_rss_huge"`
			TotalSwap               int   `json:"total_swap"`
			TotalUnevictable        int   `json:"total_unevictable"`
			TotalWriteback          int   `json:"total_writeback"`
			Unevictable             int   `json:"unevictable"`
			Writeback               int   `json:"writeback"`
		} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks struct {
		Eth0 struct {
			RxBytes   int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			RxErrors  int `json:"rx_errors"`
			RxDropped int `json:"rx_dropped"`
			TxBytes   int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
			TxErrors  int `json:"tx_errors"`
			TxDropped int `json:"tx_dropped"`
		} `json:"eth0"`
	} `json:"networks"`
}

type DockerClient struct {
	cli *docker.Client
}

func (client *DockerClient) init() {
	cli, err := docker.NewEnvClient()

	if err != nil {
		panic(err)
	}

	client.cli = cli
}

//
func (client *DockerClient) GetContainers() []types.Container {

	containers, err := client.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	return containers
}

//
func (client *DockerClient) ContainerRemove(container types.Container) {
	err := client.cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		// panic(err)
	}
}

//
func (client *DockerClient) ContainerStop(container types.Container) {
	err := client.cli.ContainerStop(context.Background(), container.ID, nil)
	if err != nil {
		// panic(err)
	}
}

//
func (client *DockerClient) ContainerStart(container types.Container) {
	err := client.cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		// panic(err)
	}
}

//
func (client *DockerClient) ContainerStats(container types.Container) Stats {
	stats, err := client.cli.ContainerStats(context.Background(), container.ID, false)
	if err != nil {
		// panic(err)
	}

	body, err := ioutil.ReadAll(stats.Body)
	var stat Stats
	err = json.Unmarshal(body, &stat)
	if err != nil {

	}
	return stat
}

func CreateClient() DockerClient {
	client := DockerClient{}

	client.init()
	return client
}
