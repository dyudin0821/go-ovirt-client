package ovirtclient_test

import (
	"testing"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// Helper functions to reduce function complexity

func createVMWithBootDevices(
	t *testing.T,
	helper ovirtclient.TestHelper,
	bootDevices []ovirtclient.BootDevice,
) ovirtclient.VM {
	cluster := helper.GetClusterID()
	template := helper.GetBlankTemplateID()

	osParams := ovirtclient.NewVMOSParameters().
		MustWithType("rhel_8x64").
		MustWithBootDevices(bootDevices)

	vmParams := ovirtclient.NewCreateVMParams().
		WithOS(osParams)

	vm, err := helper.GetClient().CreateVM(cluster, template, "test-vm-boot-update", vmParams)
	if err != nil {
		t.Fatalf("Failed to create VM (%v)", err)
	}
	return vm
}

func cleanupVM(t *testing.T, vm ovirtclient.VM) {
	if err := vm.Remove(); err != nil {
		t.Logf("Failed to remove test VM: %v", err)
	}
}

func verifyBootDeviceCount(t *testing.T, vm ovirtclient.VM, expectedCount int) {
	os := vm.OS()
	if os == nil {
		t.Fatal("VM OS is nil")
	}

	bootDevices := os.BootDevices()
	if len(bootDevices) != expectedCount {
		t.Fatalf("Expected %d boot devices, got %d", expectedCount, len(bootDevices))
	}
}

func updateVMBootDevices(
	t *testing.T,
	vm ovirtclient.VM,
	bootDevices []ovirtclient.BootDevice,
) ovirtclient.VM {
	updateParams := ovirtclient.UpdateVMParams().MustWithBootDevices(bootDevices)

	updatedVM, err := vm.Update(updateParams)
	if err != nil {
		t.Fatalf("Failed to update VM boot devices (%v)", err)
	}
	return updatedVM
}

func verifyBootDevices(t *testing.T, vm ovirtclient.VM, expectedDevices []ovirtclient.BootDevice) {
	os := vm.OS()
	if os == nil {
		t.Fatal("VM OS is nil")
	}

	bootDevices := os.BootDevices()
	if len(bootDevices) != len(expectedDevices) {
		t.Fatalf("Expected %d boot devices, got %d", len(expectedDevices), len(bootDevices))
	}

	for i, expected := range expectedDevices {
		if bootDevices[i] != expected {
			t.Errorf("Expected boot device at position %d to be %s, got %s", i, expected, bootDevices[i])
		}
	}
}

func TestVMUpdateBootDevices(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	// Create VM with initial boot sequence
	vm := createVMWithBootDevices(t, helper, []ovirtclient.BootDevice{
		ovirtclient.BootDeviceNetwork,
		ovirtclient.BootDeviceCDROM,
	})
	defer cleanupVM(t, vm)

	// Verify initial boot sequence
	verifyBootDeviceCount(t, vm, 2)

	// Update boot sequence
	newBootDevices := []ovirtclient.BootDevice{
		ovirtclient.BootDeviceHD,
		ovirtclient.BootDeviceCDROM,
		ovirtclient.BootDeviceNetwork,
	}

	updatedVM := updateVMBootDevices(t, vm, newBootDevices)

	// Verify updated boot sequence
	verifyBootDevices(t, updatedVM, newBootDevices)
}

func TestVMUpdateBootDevicesInvalidDevice(t *testing.T) {
	t.Parallel()

	invalidBootDevices := []ovirtclient.BootDevice{
		ovirtclient.BootDevice("invalid"),
	}

	_, err := ovirtclient.UpdateVMParams().WithBootDevices(invalidBootDevices)
	if err == nil {
		t.Fatal("Expected error when setting invalid boot device, got nil")
	}
}

func TestVMUpdateBootDevicesEmpty(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	// Create VM with boot sequence
	vm := createVMWithBootDevices(t, helper, []ovirtclient.BootDevice{
		ovirtclient.BootDeviceNetwork,
	})
	defer cleanupVM(t, vm)

	// Update with empty boot devices (should not change anything)
	updateParams := ovirtclient.UpdateVMParams().MustWithName("new-name")

	updatedVM, err := vm.Update(updateParams)
	if err != nil {
		t.Fatalf("Failed to update VM (%v)", err)
	}

	// Verify boot sequence unchanged
	verifyBootDeviceCount(t, updatedVM, 1)
}
