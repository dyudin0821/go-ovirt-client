package ovirtclient_test

import (
	"fmt"
	"testing"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
)

func TestCDROMAttachment(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("cdrom_test_%s", helper.GenerateRandomID(5)),
		ovirtclient.CreateVMParams(),
	)

	// Initially there should be no CDROMs with ISOs attached
	cdroms, err := vm.ListCDROMs()
	if err != nil {
		t.Fatalf("Failed to list CDROMs for VM %s (%v)", vm.ID(), err)
	}
	t.Logf("VM %s has %d CDROM(s)", vm.ID(), len(cdroms))

	// Attach an ISO
	isoID := "test-iso-123"
	cdrom := assertCanAttachISO(t, vm, isoID)
	assertCDROMMatches(t, cdrom, vm, isoID)

	// List CDROMs again to verify
	cdroms, err = vm.ListCDROMs()
	if err != nil {
		t.Fatalf("Failed to list CDROMs after attaching ISO (%v)", err)
	}
	if len(cdroms) == 0 {
		t.Fatalf("Expected at least one CDROM after attachment, got 0")
	}

	// Get the CDROM by ID
	retrievedCDROM := assertCanGetCDROM(t, vm, cdrom.ID())
	if retrievedCDROM.FileID() != isoID {
		t.Fatalf("Retrieved CDROM file ID mismatch (expected %s, got %s)", isoID, retrievedCDROM.FileID())
	}
}

func TestCDROMChange(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("cdrom_change_test_%s", helper.GenerateRandomID(5)),
		ovirtclient.CreateVMParams(),
	)

	// Attach initial ISO
	isoID1 := "test-iso-001"
	cdrom := assertCanAttachISO(t, vm, isoID1)

	// Change to a different ISO
	isoID2 := "test-iso-002"
	changedCDROM := assertCanChangeISO(t, vm, cdrom.ID(), isoID2)
	if changedCDROM.FileID() != isoID2 {
		t.Fatalf("Changed CDROM file ID mismatch (expected %s, got %s)", isoID2, changedCDROM.FileID())
	}
}

func TestCDROMEject(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("cdrom_eject_test_%s", helper.GenerateRandomID(5)),
		ovirtclient.CreateVMParams(),
	)

	// Attach an ISO
	isoID := "test-iso-eject"
	cdrom := assertCanAttachISO(t, vm, isoID)

	// Eject the ISO
	ejectedCDROM := assertCanEjectISO(t, vm, cdrom.ID())
	if ejectedCDROM.FileID() != "" {
		t.Fatalf("Expected empty file ID after ejecting, got %s", ejectedCDROM.FileID())
	}
	if ejectedCDROM.FileName() != "" {
		t.Fatalf("Expected empty file name after ejecting, got %s", ejectedCDROM.FileName())
	}
}

func TestCDROMConvenienceMethods(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)

	vm := assertCanCreateVM(
		t,
		helper,
		fmt.Sprintf("cdrom_convenience_test_%s", helper.GenerateRandomID(5)),
		ovirtclient.CreateVMParams(),
	)

	// Test CDROM.Change()
	isoID1 := "test-iso-conv-001"
	cdrom := assertCanAttachISO(t, vm, isoID1)

	isoID2 := "test-iso-conv-002"
	changedCDROM, err := cdrom.Change(isoID2)
	if err != nil {
		t.Fatalf("Failed to change ISO using convenience method (%v)", err)
	}
	if changedCDROM.FileID() != isoID2 {
		t.Fatalf("CDROM file ID mismatch after change (expected %s, got %s)", isoID2, changedCDROM.FileID())
	}

	// Test CDROM.Eject()
	ejectedCDROM, err := changedCDROM.Eject()
	if err != nil {
		t.Fatalf("Failed to eject ISO using convenience method (%v)", err)
	}
	if ejectedCDROM.FileID() != "" {
		t.Fatalf("Expected empty file ID after ejecting, got %s", ejectedCDROM.FileID())
	}

	// Test CDROM.VM()
	retrievedVM, err := cdrom.VM()
	if err != nil {
		t.Fatalf("Failed to retrieve VM from CDROM (%v)", err)
	}
	if retrievedVM.ID() != vm.ID() {
		t.Fatalf("Retrieved VM ID mismatch (expected %s, got %s)", vm.ID(), retrievedVM.ID())
	}
}

// Helper functions

func assertCanAttachISO(t *testing.T, vm ovirtclient.VM, isoID string) ovirtclient.CDROM {
	cdrom, err := vm.AttachISO(isoID)
	if err != nil {
		t.Fatalf("Failed to attach ISO %s to VM %s (%v)", isoID, vm.ID(), err)
	}
	if cdrom.VMID() != vm.ID() {
		t.Fatalf("CDROM VM ID mismatch (expected %s, got %s)", vm.ID(), cdrom.VMID())
	}
	if cdrom.FileID() != isoID {
		t.Fatalf("CDROM file ID mismatch (expected %s, got %s)", isoID, cdrom.FileID())
	}
	return cdrom
}

func assertCanGetCDROM(t *testing.T, vm ovirtclient.VM, cdromID ovirtclient.CDROMID) ovirtclient.CDROM {
	cdrom, err := vm.GetCDROM(cdromID)
	if err != nil {
		t.Fatalf("Failed to get CDROM %s for VM %s (%v)", cdromID, vm.ID(), err)
	}
	if cdrom.ID() != cdromID {
		t.Fatalf("Retrieved CDROM ID mismatch (expected %s, got %s)", cdromID, cdrom.ID())
	}
	return cdrom
}

func assertCanChangeISO(t *testing.T, vm ovirtclient.VM, cdromID ovirtclient.CDROMID, isoID string) ovirtclient.CDROM {
	cdrom, err := vm.ChangeISO(cdromID, isoID)
	if err != nil {
		t.Fatalf("Failed to change ISO for CDROM %s to %s (%v)", cdromID, isoID, err)
	}
	if cdrom.ID() != cdromID {
		t.Fatalf("CDROM ID mismatch after change (expected %s, got %s)", cdromID, cdrom.ID())
	}
	return cdrom
}

func assertCanEjectISO(t *testing.T, vm ovirtclient.VM, cdromID ovirtclient.CDROMID) ovirtclient.CDROM {
	cdrom, err := vm.EjectISO(cdromID)
	if err != nil {
		t.Fatalf("Failed to eject ISO from CDROM %s (%v)", cdromID, err)
	}
	if cdrom.ID() != cdromID {
		t.Fatalf("CDROM ID mismatch after eject (expected %s, got %s)", cdromID, cdrom.ID())
	}
	return cdrom
}

func assertCDROMMatches(t *testing.T, cdrom ovirtclient.CDROM, vm ovirtclient.VM, isoID string) {
	if cdrom.VMID() != vm.ID() {
		t.Fatalf("CDROM VM ID does not match (expected %s, got %s)", vm.ID(), cdrom.VMID())
	}
	if cdrom.FileID() != isoID {
		t.Fatalf("CDROM file ID does not match (expected %s, got %s)", isoID, cdrom.FileID())
	}
}
