package ovirtclient_test

import (
	"fmt"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// This example demonstrates how to attach an ISO image to a virtual machine.
func ExampleClient_AttachCDROM() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// Get a VM to attach the ISO to
	vms, err := client.ListVMs()
	if err != nil {
		panic(fmt.Errorf("failed to list VMs (%w)", err))
	}
	if len(vms) == 0 {
		panic("no VMs available")
	}
	vm := vms[0]

	// Attach an ISO image to the VM's CDROM
	// Replace "iso-image-id" with the actual ID of an ISO file in your storage domain
	isoImageID := "iso-image-id"
	cdrom, err := client.AttachCDROM(vm.ID(), isoImageID)
	if err != nil {
		panic(fmt.Errorf("failed to attach ISO (%w)", err))
	}

	fmt.Printf("Attached ISO %s to CDROM %s on VM %s\n", cdrom.FileID(), cdrom.ID(), cdrom.VMID())
}

// This example demonstrates how to list all CDROM attachments on a VM.
func ExampleClient_ListCDROMs() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// Get a VM to list CDROMs from
	vms, err := client.ListVMs()
	if err != nil {
		panic(fmt.Errorf("failed to list VMs (%w)", err))
	}
	if len(vms) == 0 {
		panic("no VMs available")
	}
	vm := vms[0]

	// List all CDROM attachments for the VM
	cdroms, err := client.ListCDROMs(vm.ID())
	if err != nil {
		panic(fmt.Errorf("failed to list CDROMs (%w)", err))
	}

	for _, cdrom := range cdroms {
		if cdrom.FileID() != "" {
			fmt.Printf("CDROM %s has ISO %s (%s) attached\n", cdrom.ID(), cdrom.FileID(), cdrom.FileName())
		} else {
			fmt.Printf("CDROM %s is empty\n", cdrom.ID())
		}
	}
}

// This example demonstrates how to change the ISO in a CDROM.
func ExampleClient_ChangeCDROM() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// Get a VM and its CDROMs
	vms, err := client.ListVMs()
	if err != nil {
		panic(fmt.Errorf("failed to list VMs (%w)", err))
	}
	if len(vms) == 0 {
		panic("no VMs available")
	}
	vm := vms[0]

	cdroms, err := client.ListCDROMs(vm.ID())
	if err != nil {
		panic(fmt.Errorf("failed to list CDROMs (%w)", err))
	}
	if len(cdroms) == 0 {
		panic("no CDROMs available")
	}
	cdrom := cdroms[0]

	// Change to a different ISO
	// Replace "new-iso-image-id" with the actual ID of an ISO file in your storage domain
	newISOImageID := "new-iso-image-id"
	updatedCDROM, err := client.ChangeCDROM(vm.ID(), cdrom.ID(), newISOImageID)
	if err != nil {
		panic(fmt.Errorf("failed to change ISO (%w)", err))
	}

	fmt.Printf("Changed CDROM %s to ISO %s\n", updatedCDROM.ID(), updatedCDROM.FileID())
}

// This example demonstrates how to eject an ISO from a CDROM.
func ExampleClient_EjectCDROM() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// Get a VM and its CDROMs
	vms, err := client.ListVMs()
	if err != nil {
		panic(fmt.Errorf("failed to list VMs (%w)", err))
	}
	if len(vms) == 0 {
		panic("no VMs available")
	}
	vm := vms[0]

	cdroms, err := client.ListCDROMs(vm.ID())
	if err != nil {
		panic(fmt.Errorf("failed to list CDROMs (%w)", err))
	}
	if len(cdroms) == 0 {
		panic("no CDROMs available")
	}
	cdrom := cdroms[0]

	// Eject the ISO
	ejectedCDROM, err := client.EjectCDROM(vm.ID(), cdrom.ID())
	if err != nil {
		panic(fmt.Errorf("failed to eject ISO (%w)", err))
	}

	if ejectedCDROM.FileID() == "" {
		fmt.Printf("Successfully ejected ISO from CDROM %s\n", ejectedCDROM.ID())
	}
}
