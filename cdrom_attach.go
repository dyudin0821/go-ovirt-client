package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) AttachCDROM(
	vmID VMID,
	isoImageID string,
	retries ...RetryStrategy,
) (result CDROM, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("attaching ISO %s to VM %s", isoImageID, vmID),
		o.logger,
		retries,
		func() error {
			cdromBuilder := ovirtsdk.NewCdromBuilder()
			fileBuilder := ovirtsdk.NewFileBuilder()
			fileBuilder.Id(isoImageID)
			cdromBuilder.File(fileBuilder.MustBuild())

			addRequest := o.conn.SystemService().VmsService().VmService(string(vmID)).CdromsService().Add()
			addRequest.Cdrom(cdromBuilder.MustBuild())
			response, err := addRequest.Send()
			if err != nil {
				return wrap(
					err,
					EUnidentified,
					"failed to attach ISO %s to VM %s",
					isoImageID,
					vmID,
				)
			}

			cdrom, ok := response.Cdrom()
			if !ok {
				return newFieldNotFound("CDROM response", "cdrom")
			}
			result, err = convertSDKCDROM(cdrom, vmID, o)
			if err != nil {
				return wrap(err, EUnidentified, "failed to convert SDK CDROM")
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) AttachCDROM(
	vmID VMID,
	isoImageID string,
	_ ...RetryStrategy,
) (CDROM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, ok := m.vms[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", vmID)
	}

	// Check if the ISO file exists in mock storage domains
	// For mock, we'll just validate that the ID is not empty
	if isoImageID == "" {
		return nil, newError(EBadArgument, "ISO image ID cannot be empty")
	}

	cdromID := CDROMID(m.GenerateUUID())
	cdromObj := &cdrom{
		client:   m,
		id:       cdromID,
		vmid:     vm.ID(),
		fileID:   isoImageID,
		fileName: fmt.Sprintf("iso-%s.iso", isoImageID),
	}

	if m.vmCDROMsByVM == nil {
		m.vmCDROMsByVM = make(map[VMID]map[CDROMID]*cdrom)
	}
	if m.vmCDROMsByVM[vm.ID()] == nil {
		m.vmCDROMsByVM[vm.ID()] = make(map[CDROMID]*cdrom)
	}

	m.vmCDROMsByVM[vm.ID()][cdromID] = cdromObj

	return cdromObj, nil
}
