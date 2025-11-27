package ovirtclient_test

import (
	"fmt"
	"log"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
)

// This example demonstrates how to attach an ISO image to a virtual machine.
func ExampleClient_AttachCDROM() {
	// Create a client (in real usage, use ovirtclient.New() with proper credentials)
	client := ovirtclient.NewMock()

	// Get or create a VM
	vmID := ovirtclient.VMID("vm-id-here")

	// Attach an ISO image to the VM's CDROM
	// The isoImageID should be the ID of an ISO file available in your storage domain
	isoImageID := "my-iso-image-id"
	cdrom, err := client.AttachCDROM(vmID, isoImageID)
	if err != nil {
		log.Fatalf("Failed to attach ISO: %v", err)
	}

	fmt.Printf("Attached ISO %s to CDROM %s on VM %s\n", cdrom.FileID(), cdrom.ID(), cdrom.VMID())
	// Output:
}

// This example demonstrates how to list all CDROM attachments on a VM.
func ExampleClient_ListCDROMs() {
	// Create a client (in real usage, use ovirtclient.New() with proper credentials)
	client := ovirtclient.NewMock()

	// List all CDROM attachments for a VM
	vmID := ovirtclient.VMID("vm-id-here")
	cdroms, err := client.ListCDROMs(vmID)
	if err != nil {
		log.Fatalf("Failed to list CDROMs: %v", err)
	}

	for _, cdrom := range cdroms {
		if cdrom.FileID() != "" {
			fmt.Printf("CDROM %s has ISO %s (%s) attached\n", cdrom.ID(), cdrom.FileID(), cdrom.FileName())
		} else {
			fmt.Printf("CDROM %s is empty\n", cdrom.ID())
		}
	}
	// Output:
}

// This example demonstrates how to change the ISO in a CDROM.
func ExampleClient_ChangeCDROM() {
	// Create a client (in real usage, use ovirtclient.New() with proper credentials)
	client := ovirtclient.NewMock()

	vmID := ovirtclient.VMID("vm-id-here")
	cdromID := ovirtclient.CDROMID("cdrom-id-here")

	// Change to a different ISO
	newISOImageID := "new-iso-image-id"
	updatedCDROM, err := client.ChangeCDROM(vmID, cdromID, newISOImageID)
	if err != nil {
		log.Fatalf("Failed to change ISO: %v", err)
	}

	fmt.Printf("Changed CDROM %s to ISO %s\n", updatedCDROM.ID(), updatedCDROM.FileID())
	// Output:
}

// This example demonstrates how to eject an ISO from a CDROM.
func ExampleClient_EjectCDROM() {
	// Create a client (in real usage, use ovirtclient.New() with proper credentials)
	client := ovirtclient.NewMock()

	vmID := ovirtclient.VMID("vm-id-here")
	cdromID := ovirtclient.CDROMID("cdrom-id-here")

	// Eject the ISO
	ejectedCDROM, err := client.EjectCDROM(vmID, cdromID)
	if err != nil {
		log.Fatalf("Failed to eject ISO: %v", err)
	}

	if ejectedCDROM.FileID() == "" {
		fmt.Printf("Successfully ejected ISO from CDROM %s\n", ejectedCDROM.ID())
	}
	// Output:
}
