package ovirtclient_test

import (
	"fmt"
	"testing"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// TestVMKernelParametersCreate tests creating a VM with kernel parameters.
func TestVMKernelParametersCreate(t *testing.T) {
	helper := getHelper(t)

	// Create VM with kernel parameters
	cmdline := testCmdlineSerial
	customKernelCmdline := "quiet splash"
	initrd := "/boot/initrd.img"
	kernel := "/boot/vmlinuz"

	params := ovirtclient.NewCreateVMParams().
		WithOS(
			ovirtclient.NewVMOSParameters().
				MustWithType("rhel_8x64").
				MustWithCmdline(cmdline).
				MustWithCustomKernelCmdline(customKernelCmdline).
				MustWithInitrd(initrd).
				MustWithKernel(kernel),
		)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("test_%s", helper.GenerateRandomID(5)),
		params,
	)
	defer func() {
		if err := vm.Remove(); err != nil {
			t.Logf("Failed to remove VM: %v", err)
		}
	}()

	// Verify kernel parameters
	os := vm.OS()
	if os.Cmdline() == nil || *os.Cmdline() != cmdline {
		t.Fatalf("Cmdline mismatch: expected %s, got %v", cmdline, os.Cmdline())
	}
	if os.CustomKernelCmdline() == nil || *os.CustomKernelCmdline() != customKernelCmdline {
		t.Fatalf("CustomKernelCmdline mismatch: expected %s, got %v", customKernelCmdline, os.CustomKernelCmdline())
	}
	if os.Initrd() == nil || *os.Initrd() != initrd {
		t.Fatalf("Initrd mismatch: expected %s, got %v", initrd, os.Initrd())
	}
	if os.Kernel() == nil || *os.Kernel() != kernel {
		t.Fatalf("Kernel mismatch: expected %s, got %v", kernel, os.Kernel())
	}
}

// TestVMKernelParametersUpdate tests updating kernel parameters on an existing VM.
func TestVMKernelParametersUpdate(t *testing.T) {
	helper := getHelper(t)
	client := helper.GetClient()

	// Create VM without kernel parameters
	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("test_%s", helper.GenerateRandomID(5)),
		nil,
	)
	defer func() {
		if err := vm.Remove(); err != nil {
			t.Logf("Failed to remove VM: %v", err)
		}
	}()

	// Update VM with kernel parameters
	newCmdline := "console=tty1"
	newCustomKernelCmdline := "noquiet nosplash"
	newInitrd := "/boot/initrd-new.img"
	newKernel := "/boot/vmlinuz-new"

	updateParams := ovirtclient.UpdateVMParams().
		MustWithCmdline(newCmdline).
		MustWithCustomKernelCmdline(newCustomKernelCmdline).
		MustWithInitrd(newInitrd).
		MustWithKernel(newKernel)

	updatedVM, err := client.UpdateVM(vm.ID(), updateParams)
	if err != nil {
		t.Fatalf("Failed to update VM with kernel parameters: %v", err)
	}

	// Verify updated kernel parameters
	os := updatedVM.OS()
	if os.Cmdline() == nil || *os.Cmdline() != newCmdline {
		t.Fatalf("Updated cmdline mismatch: expected %s, got %v", newCmdline, os.Cmdline())
	}
	if os.CustomKernelCmdline() == nil || *os.CustomKernelCmdline() != newCustomKernelCmdline {
		t.Fatalf("Updated customKernelCmdline mismatch: expected %s, got %v", newCustomKernelCmdline, os.CustomKernelCmdline())
	}
	if os.Initrd() == nil || *os.Initrd() != newInitrd {
		t.Fatalf("Updated initrd mismatch: expected %s, got %v", newInitrd, os.Initrd())
	}
	if os.Kernel() == nil || *os.Kernel() != newKernel {
		t.Fatalf("Updated kernel mismatch: expected %s, got %v", newKernel, os.Kernel())
	}
}

const testCmdlineSerial = "console=ttyS0"

// TestVMKernelParametersPartial tests setting only some kernel parameters.
func TestVMKernelParametersPartial(t *testing.T) {
	helper := getHelper(t)

	// Create VM with only cmdline parameter
	cmdline := testCmdlineSerial

	params := ovirtclient.NewCreateVMParams().
		WithOS(
			ovirtclient.NewVMOSParameters().
				MustWithCmdline(cmdline),
		)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("test_%s", helper.GenerateRandomID(5)),
		params,
	)
	defer func() {
		if err := vm.Remove(); err != nil {
			t.Logf("Failed to remove VM: %v", err)
		}
	}()

	// Verify only cmdline is set
	os := vm.OS()
	if os.Cmdline() == nil || *os.Cmdline() != cmdline {
		t.Fatalf("Cmdline mismatch: expected %s, got %v", cmdline, os.Cmdline())
	}
	if os.CustomKernelCmdline() != nil {
		t.Fatalf("CustomKernelCmdline should be nil, got %v", *os.CustomKernelCmdline())
	}
	if os.Initrd() != nil {
		t.Fatalf("Initrd should be nil, got %v", *os.Initrd())
	}
	if os.Kernel() != nil {
		t.Fatalf("Kernel should be nil, got %v", *os.Kernel())
	}
}

// TestVMKernelParametersMock tests kernel parameters with mock client.
func TestVMKernelParametersMock(t *testing.T) {
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("Failed to create test helper: %v", err)
	}
	client := helper.GetClient()

	// Create VM with kernel parameters
	cmdline := testCmdlineSerial
	kernel := "/boot/vmlinuz"

	params := ovirtclient.NewCreateVMParams().
		WithOS(
			ovirtclient.NewVMOSParameters().
				MustWithCmdline(cmdline).
				MustWithKernel(kernel),
		)

	vm, err := client.CreateVM(
		helper.GetClusterID(),
		helper.GetBlankTemplateID(),
		fmt.Sprintf("test_%s", helper.GenerateRandomID(5)),
		params,
	)
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// Verify kernel parameters
	os := vm.OS()
	if os.Cmdline() == nil || *os.Cmdline() != cmdline {
		t.Fatalf("Cmdline mismatch: expected %s, got %v", cmdline, os.Cmdline())
	}
	if os.Kernel() == nil || *os.Kernel() != kernel {
		t.Fatalf("Kernel mismatch: expected %s, got %v", kernel, os.Kernel())
	}

	// Test update
	newCmdline := "console=tty1"
	updateParams := ovirtclient.UpdateVMParams().MustWithCmdline(newCmdline)

	updatedVM, err := client.UpdateVM(vm.ID(), updateParams)
	if err != nil {
		t.Fatalf("Failed to update VM: %v", err)
	}

	os = updatedVM.OS()
	if os.Cmdline() == nil || *os.Cmdline() != newCmdline {
		t.Fatalf("Updated cmdline mismatch: expected %s, got %v", newCmdline, os.Cmdline())
	}
	// Kernel should remain unchanged
	if os.Kernel() == nil || *os.Kernel() != kernel {
		t.Fatalf("Kernel should remain unchanged: expected %s, got %v", kernel, os.Kernel())
	}
}
