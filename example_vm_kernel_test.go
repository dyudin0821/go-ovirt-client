package ovirtclient_test

import (
	"fmt"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// ExampleVMClient_createWithKernelParameters demonstrates how to create a VM with custom kernel parameters.
func ExampleVMClient_createWithKernelParameters() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// This is the cluster the VM will be created on.
	clusterID := helper.GetClusterID()
	// Use the blank template as a starting point.
	templateID := helper.GetBlankTemplateID()
	// Set the VM name
	name := "test-vm-with-kernel-params"

	// Create the optional parameters with kernel parameters
	params := ovirtclient.NewCreateVMParams().
		WithOS(
			ovirtclient.NewVMOSParameters().
				MustWithType("rhel_8x64").
				// Set custom kernel command line parameters
				MustWithCmdline("console=ttyS0 quiet splash").
				// Set custom kernel path on ISO storage domain
				MustWithKernel("/boot/vmlinuz-custom").
				// Set custom initrd path on ISO storage domain
				MustWithInitrd("/boot/initrd-custom.img").
				// Set custom host kernel command line part
				MustWithCustomKernelCmdline("intel_iommu=on"),
		)

	// Create the VM...
	vm, err := client.CreateVM(clusterID, templateID, name, params)
	if err != nil {
		panic(fmt.Sprintf("failed to create VM (%v)", err))
	}

	// Read kernel parameters from the created VM
	os := vm.OS()
	if cmdline := os.Cmdline(); cmdline != nil {
		fmt.Printf("VM created with cmdline: %s\n", *cmdline)
	}
	if kernel := os.Kernel(); kernel != nil {
		fmt.Printf("VM created with kernel: %s\n", *kernel)
	}
	if initrd := os.Initrd(); initrd != nil {
		fmt.Printf("VM created with initrd: %s\n", *initrd)
	}

	// ... and then remove it. Alternatively, you could call client.RemoveVM(vm.ID()).
	if err := vm.Remove(); err != nil {
		panic(fmt.Sprintf("failed to remove VM (%v)", err))
	}
}

// ExampleVMClient_updateKernelParameters demonstrates how to update kernel parameters on an existing VM.
func ExampleVMClient_updateKernelParameters() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// This is the cluster the VM will be created on.
	clusterID := helper.GetClusterID()
	// Use the blank template as a starting point.
	templateID := helper.GetBlankTemplateID()
	// Set the VM name
	name := "test-vm-update-kernel"

	// Create a VM without kernel parameters
	vm, err := client.CreateVM(clusterID, templateID, name, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create VM (%v)", err))
	}

	// Update the VM with kernel parameters
	updateParams := ovirtclient.UpdateVMParams().
		MustWithCmdline("console=tty1 debug").
		MustWithKernel("/boot/vmlinuz-new").
		MustWithInitrd("/boot/initrd-new.img")

	updatedVM, err := client.UpdateVM(vm.ID(), updateParams)
	if err != nil {
		panic(fmt.Sprintf("failed to update VM (%v)", err))
	}

	// Read updated kernel parameters
	os := updatedVM.OS()
	if cmdline := os.Cmdline(); cmdline != nil {
		fmt.Printf("VM updated with cmdline: %s\n", *cmdline)
	}

	// Clean up
	if err := updatedVM.Remove(); err != nil {
		panic(fmt.Sprintf("failed to remove VM (%v)", err))
	}
}
