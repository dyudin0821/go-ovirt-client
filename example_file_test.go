package ovirtclient_test

import (
	"fmt"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// ExampleClient_listStorageDomainFiles demonstrates how to list ISO images and VFD files in a storage domain.
func ExampleClient_listStorageDomainFiles() { //nolint:testableexamples
	// Create the helper for testing. Alternatively, you could create a production client with ovirtclient.New()
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	// Get the oVirt client
	client := helper.GetClient()

	// List all storage domains
	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		panic(fmt.Errorf("failed to list storage domains (%w)", err))
	}

	// Find an ISO storage domain (you may need to adjust this based on your environment)
	var isoStorageDomain ovirtclient.StorageDomain
	for _, sd := range storageDomains {
		// In a real environment, you would check if the storage domain is of type ISO
		// For this example, we'll use the first storage domain
		isoStorageDomain = sd
		break
	}

	if isoStorageDomain == nil {
		panic("no storage domain found")
	}

	// List files in the storage domain without forcing refresh (better performance)
	files, err := client.ListStorageDomainFiles(isoStorageDomain.ID(), false)
	if err != nil {
		panic(fmt.Errorf("failed to list files in storage domain (%w)", err))
	}

	fmt.Printf("Found %d files in storage domain %s\n", len(files), isoStorageDomain.Name())
	
	// Print details of each file
	for _, file := range files {
		fmt.Printf("- File: %s (ID: %s, Type: %s)\n", file.Name(), file.ID(), file.Type())
	}
}

// ExampleClient_listStorageDomainFilesWithRefresh demonstrates listing files with forced refresh from storage.
func ExampleClient_listStorageDomainFilesWithRefresh() { //nolint:testableexamples
	// Create the helper for testing
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	
	client := helper.GetClient()
	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		panic(fmt.Errorf("failed to list storage domains (%w)", err))
	}

	if len(storageDomains) == 0 {
		panic("no storage domains found")
	}

	storageDomainID := storageDomains[0].ID()

	// List files with refresh=true to force the server to update from storage
	// Note: This may have performance impact, use only when necessary
	files, err := client.ListStorageDomainFiles(storageDomainID, true)
	if err != nil {
		panic(fmt.Errorf("failed to list files with refresh (%w)", err))
	}

	fmt.Printf("Found %d ISO/VFD files after refresh\n", len(files))
}

// ExampleClient_getStorageDomainFile demonstrates how to get a specific file from a storage domain.
func ExampleClient_getStorageDomainFile() { //nolint:testableexamples
	// Create the helper for testing
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	
	client := helper.GetClient()
	
	// First, list files to get a file ID
	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		panic(fmt.Errorf("failed to list storage domains (%w)", err))
	}

	if len(storageDomains) == 0 {
		panic("no storage domains found")
	}

	storageDomainID := storageDomains[0].ID()
	
	files, err := client.ListStorageDomainFiles(storageDomainID, false)
	if err != nil {
		panic(fmt.Errorf("failed to list files (%w)", err))
	}

	if len(files) == 0 {
		fmt.Println("No files found in storage domain")
		return
	}

	// Get details of the first file
	fileID := files[0].ID()
	file, err := client.GetStorageDomainFile(storageDomainID, fileID)
	if err != nil {
		panic(fmt.Errorf("failed to get file (%w)", err))
	}

	fmt.Printf("File details:\n")
	fmt.Printf("  Name: %s\n", file.Name())
	fmt.Printf("  ID: %s\n", file.ID())
	fmt.Printf("  Type: %s\n", file.Type())
	fmt.Printf("  Storage Domain ID: %s\n", file.StorageDomainID())
	
	// Get the storage domain object
	storageDomain, err := file.StorageDomain()
	if err != nil {
		panic(fmt.Errorf("failed to get storage domain (%w)", err))
	}
	fmt.Printf("  Storage Domain Name: %s\n", storageDomain.Name())
}

// ExampleClient_attachISOWithFileList demonstrates how to list ISO files and attach one to a VM.
func ExampleClient_attachISOWithFileList() { //nolint:testableexamples
	// Create the helper for testing
	helper, err := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
	if err != nil {
		panic(fmt.Errorf("failed to create live test helper (%w)", err))
	}
	
	client := helper.GetClient()
	
	// Get storage domain and list ISO files
	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		panic(fmt.Errorf("failed to list storage domains (%w)", err))
	}

	if len(storageDomains) == 0 {
		panic("no storage domains found")
	}

	files, err := client.ListStorageDomainFiles(storageDomains[0].ID(), false)
	if err != nil {
		panic(fmt.Errorf("failed to list files (%w)", err))
	}

	if len(files) == 0 {
		fmt.Println("No ISO files available")
		return
	}

	// Get a VM to attach the ISO to
	vms, err := client.ListVMs()
	if err != nil {
		panic(fmt.Errorf("failed to list VMs (%w)", err))
	}

	if len(vms) == 0 {
		panic("no VMs found")
	}

	vm := vms[0]
	isoFile := files[0]

	fmt.Printf("Attaching ISO '%s' to VM '%s'\n", isoFile.Name(), vm.Name())

	// Attach the ISO file to the VM's CDROM
	cdrom, err := client.AttachCDROM(vm.ID(), string(isoFile.ID()))
	if err != nil {
		panic(fmt.Errorf("failed to attach ISO (%w)", err))
	}

	fmt.Printf("Successfully attached ISO as CDROM %s\n", cdrom.ID())
}
