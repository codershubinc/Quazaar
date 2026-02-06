package usage

import (
	"Quazaar/internal/system/info/disk"
	"testing"
)

func TestGetSystemUsage(t *testing.T) {
	usage, err := GetSystemUsage()
	if err != nil {
		t.Fatalf("GetSystemUsage failed: %v", err)
	}

	if usage.CPUUsage < 0 || usage.CPUUsage > 100 {
		t.Errorf("Invalid CPU Usage: %f", usage.CPUUsage)
	}
	if usage.MemoryUsage < 0 || usage.MemoryUsage > 100 {
		t.Errorf("Invalid Memory Usage: %f", usage.MemoryUsage)
	}
	// GPU might be 0 if no Nvidia card, so we don't error on it.
	t.Logf("CPU: %f%%, Mem: %f%%, GPU: %f%%", usage.CPUUsage, usage.MemoryUsage, usage.GPUUsage)
}

func TestGetStorageUsage(t *testing.T) {
	storage, err := disk.GetDiskUsage()
	if err != nil {
		t.Fatalf("GetStorageUsage failed: %v", err)
	}

	if storage.Path != "/" {
		t.Errorf("Expected path '/', got %s", storage.Path)
	}
	if storage.Total == 0 {
		t.Error("Total storage is 0")
	}
	t.Logf("Storage: %+v", storage)
}
