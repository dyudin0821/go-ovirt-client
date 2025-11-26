package ovirtclient_test

import (
	"testing"

	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
)

func TestBootDevice(t *testing.T) {
	t.Parallel()

	// Test valid boot devices
	validDevices := []ovirtclient.BootDevice{
		ovirtclient.BootDeviceHD,
		ovirtclient.BootDeviceNetwork,
		ovirtclient.BootDeviceCDROM,
	}

	for _, device := range validDevices {
		if err := device.Validate(); err != nil {
			t.Fatalf("Valid boot device %s failed validation: %v", device, err)
		}
	}

	// Test invalid boot device
	invalidDevice := ovirtclient.BootDevice("invalid")
	if err := invalidDevice.Validate(); err == nil {
		t.Fatalf("Invalid boot device should have failed validation")
	}
}

func TestBootDeviceValues(t *testing.T) {
	t.Parallel()

	devices := ovirtclient.BootDeviceValues()
	if len(devices) != 3 {
		t.Fatalf("Expected 3 boot device values, got %d", len(devices))
	}

	expected := map[ovirtclient.BootDevice]bool{
		ovirtclient.BootDeviceHD:      true,
		ovirtclient.BootDeviceNetwork: true,
		ovirtclient.BootDeviceCDROM:   true,
	}

	for _, device := range devices {
		if !expected[device] {
			t.Fatalf("Unexpected boot device: %s", device)
		}
	}
}

func TestVMOSParametersWithBootDevices(t *testing.T) {
	t.Parallel()

	// Test WithBootDevices
	osParams := ovirtclient.NewVMOSParameters()
	bootDevices := []ovirtclient.BootDevice{
		ovirtclient.BootDeviceNetwork,
		ovirtclient.BootDeviceCDROM,
		ovirtclient.BootDeviceHD,
	}

	osParams, err := osParams.WithBootDevices(bootDevices)
	if err != nil {
		t.Fatalf("Failed to set boot devices: %v", err)
	}

	if devices := osParams.BootDevices(); len(devices) != 3 {
		t.Fatalf("Expected 3 boot devices, got %d", len(devices))
	}

	// Test WithBootDevice (adding one at a time)
	osParams2 := ovirtclient.NewVMOSParameters()
	osParams2, err = osParams2.WithBootDevice(ovirtclient.BootDeviceHD)
	if err != nil {
		t.Fatalf("Failed to add boot device: %v", err)
	}
	osParams2, err = osParams2.WithBootDevice(ovirtclient.BootDeviceNetwork)
	if err != nil {
		t.Fatalf("Failed to add second boot device: %v", err)
	}

	if devices := osParams2.BootDevices(); len(devices) != 2 {
		t.Fatalf("Expected 2 boot devices, got %d", len(devices))
	}

	// Test MustWithBootDevices
	osParams3 := ovirtclient.NewVMOSParameters().MustWithBootDevices(bootDevices)
	if devices := osParams3.BootDevices(); len(devices) != 3 {
		t.Fatalf("Expected 3 boot devices with Must method, got %d", len(devices))
	}

	// Test invalid boot device should fail
	_, err = ovirtclient.NewVMOSParameters().WithBootDevice(ovirtclient.BootDevice("invalid"))
	if err == nil {
		t.Fatalf("Invalid boot device should have failed")
	}
}

func TestVMCreateWithBootSequence(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("Failed to create test helper: %v", err)
	}

	cluster := helper.GetClusterID()
	template := helper.GetBlankTemplateID()

	// Create VM with boot sequence
	osParams := ovirtclient.NewVMOSParameters().
		MustWithType("rhel_8x64").
		MustWithBootDevices([]ovirtclient.BootDevice{
			ovirtclient.BootDeviceNetwork,
			ovirtclient.BootDeviceCDROM,
			ovirtclient.BootDeviceHD,
		})

	vmParams := ovirtclient.NewCreateVMParams().
		WithOS(osParams)

	vm, err := helper.GetClient().CreateVM(cluster, template, "test-vm-boot", vmParams)
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}
	defer func() {
		if err := helper.GetClient().RemoveVM(vm.ID()); err != nil {
			t.Logf("Failed to remove VM: %v", err)
		}
	}()

	// Verify boot sequence was set
	os := vm.OS()
	if os == nil {
		t.Fatalf("VM OS should not be nil")
	}

	bootDevices := os.BootDevices()
	if len(bootDevices) != 3 {
		t.Fatalf("Expected 3 boot devices, got %d", len(bootDevices))
	}

	expectedSequence := []ovirtclient.BootDevice{
		ovirtclient.BootDeviceNetwork,
		ovirtclient.BootDeviceCDROM,
		ovirtclient.BootDeviceHD,
	}

	for i, expected := range expectedSequence {
		if bootDevices[i] != expected {
			t.Fatalf("Boot device at position %d: expected %s, got %s", i, expected, bootDevices[i])
		}
	}
}

func TestVMWithEmptyBootSequence(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("Failed to create test helper: %v", err)
	}

	cluster := helper.GetClusterID()
	template := helper.GetBlankTemplateID()

	// Create VM without boot sequence
	osParams := ovirtclient.NewVMOSParameters().MustWithType("rhel_8x64")

	vmParams := ovirtclient.NewCreateVMParams().WithOS(osParams)

	vm, err := helper.GetClient().CreateVM(cluster, template, "test-vm-no-boot", vmParams)
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}
	defer func() {
		if err := helper.GetClient().RemoveVM(vm.ID()); err != nil {
			t.Logf("Failed to remove VM: %v", err)
		}
	}()

	// Verify boot sequence is empty
	os := vm.OS()
	if os == nil {
		t.Fatalf("VM OS should not be nil")
	}

	bootDevices := os.BootDevices()
	if len(bootDevices) != 0 {
		t.Fatalf("Expected 0 boot devices, got %d", len(bootDevices))
	}
}
