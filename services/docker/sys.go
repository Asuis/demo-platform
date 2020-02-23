package docker

import (
	"runtime"
	"syscall"
)

type DiskStatus struct {
	All uint64
	Used uint64
	Free uint64
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks*uint64(fs.Bsize)
	disk.Free = fs.Bfree*uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

type MemStats struct {
	All uint64
	Used uint64
	Free uint64
	Self uint64
}

func MemStat() MemStats{
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	mem := MemStats{}
	mem.Self = memStat.Alloc

	sysInfo := new(syscall.Sysinfo_t)
	err:=syscall.Sysinfo(sysInfo)
	if err == nil {
		mem.All = sysInfo.Totalram * uint64(syscall.Getpagesize())
		mem.Free = sysInfo.Freeram * uint64(syscall.Getpagesize())
		mem.Used = mem.All - mem.Free
	}
	return mem
}
