package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// CDROMID is the identifier for a CDROM attachment.
type CDROMID string

// CDROMClient contains the methods required for handling CDROM (ISO) attachments to VMs.
type CDROMClient interface {
	// AttachCDROM attaches an ISO image to a VM's CDROM.
	AttachCDROM(
		vmID VMID,
		isoImageID string,
		retries ...RetryStrategy,
	) (CDROM, error)
	// GetCDROM returns a single CDROM attachment for a virtual machine.
	GetCDROM(vmID VMID, id CDROMID, retries ...RetryStrategy) (CDROM, error)
	// ListCDROMs lists all CDROM attachments for a virtual machine.
	ListCDROMs(vmID VMID, retries ...RetryStrategy) ([]CDROM, error)
	// ChangeCDROM changes the ISO image in an existing CDROM attachment.
	ChangeCDROM(vmID VMID, cdromID CDROMID, isoImageID string, retries ...RetryStrategy) (CDROM, error)
	// EjectCDROM ejects the ISO from the CDROM (removes the file reference).
	EjectCDROM(vmID VMID, cdromID CDROMID, retries ...RetryStrategy) (CDROM, error)
}

// CDROM represents a CDROM device attached to a VM that can contain an ISO image.
type CDROM interface {
	// ID returns the identifier of the CDROM.
	ID() CDROMID
	// VMID returns the ID of the virtual machine this CDROM belongs to.
	VMID() VMID
	// FileID returns the ID of the ISO file attached to this CDROM, or empty string if no ISO is attached.
	FileID() string
	// FileName returns the name of the ISO file, or empty string if no ISO is attached.
	FileName() string

	// VM fetches the virtual machine this CDROM belongs to.
	VM(retries ...RetryStrategy) (VM, error)
	// Change changes the ISO image in this CDROM.
	Change(isoImageID string, retries ...RetryStrategy) (CDROM, error)
	// Eject ejects the ISO from this CDROM.
	Eject(retries ...RetryStrategy) (CDROM, error)
}

type cdrom struct {
	client Client

	id       CDROMID
	vmid     VMID
	fileID   string
	fileName string
}

func (c *cdrom) ID() CDROMID {
	return c.id
}

func (c *cdrom) VMID() VMID {
	return c.vmid
}

func (c *cdrom) FileID() string {
	return c.fileID
}

func (c *cdrom) FileName() string {
	return c.fileName
}

func (c *cdrom) VM(retries ...RetryStrategy) (VM, error) {
	return c.client.GetVM(c.vmid, retries...)
}

func (c *cdrom) Change(isoImageID string, retries ...RetryStrategy) (CDROM, error) {
	return c.client.ChangeCDROM(c.vmid, c.id, isoImageID, retries...)
}

func (c *cdrom) Eject(retries ...RetryStrategy) (CDROM, error) {
	return c.client.EjectCDROM(c.vmid, c.id, retries...)
}

func convertSDKCDROM(object *ovirtsdk4.Cdrom, vmID VMID, o *oVirtClient) (CDROM, error) {
	id, ok := object.Id()
	if !ok {
		return nil, newFieldNotFound("cdrom", "id")
	}

	result := &cdrom{
		client: o,
		id:     CDROMID(id),
		vmid:   vmID,
	}

	// File is optional - CDROM can be empty
	if file, ok := object.File(); ok {
		if fileID, ok := file.Id(); ok {
			result.fileID = fileID
		}
		if fileName, ok := file.Name(); ok {
			result.fileName = fileName
		}
	}

	return result, nil
}
