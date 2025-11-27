package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) ChangeCDROM(
	vmID VMID,
	cdromID CDROMID,
	isoImageID string,
	retries ...RetryStrategy,
) (result CDROM, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("changing CDROM %s to ISO %s for VM %s", cdromID, isoImageID, vmID),
		o.logger,
		retries,
		func() error {
			cdromBuilder := ovirtsdk.NewCdromBuilder()
			fileBuilder := ovirtsdk.NewFileBuilder()
			fileBuilder.Id(isoImageID)
			cdromBuilder.File(fileBuilder.MustBuild())

			updateRequest := o.conn.SystemService().VmsService().VmService(string(vmID)).
				CdromsService().CdromService(string(cdromID)).Update()
			updateRequest.Cdrom(cdromBuilder.MustBuild())
			response, err := updateRequest.Send()
			if err != nil {
				return wrap(
					err,
					EUnidentified,
					"failed to change CDROM %s to ISO %s for VM %s",
					cdromID,
					isoImageID,
					vmID,
				)
			}

			cdrom, ok := response.Cdrom()
			if !ok {
				return newFieldNotFound("CDROM update response", "cdrom")
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

func (m *mockClient) ChangeCDROM(
	vmID VMID,
	cdromID CDROMID,
	isoImageID string,
	_ ...RetryStrategy,
) (CDROM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[vmID]; !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", vmID)
	}

	if m.vmCDROMsByVM == nil || m.vmCDROMsByVM[vmID] == nil {
		return nil, newError(ENotFound, "CDROM with ID %s not found for VM %s", cdromID, vmID)
	}

	cdrom, ok := m.vmCDROMsByVM[vmID][cdromID]
	if !ok {
		return nil, newError(ENotFound, "CDROM with ID %s not found for VM %s", cdromID, vmID)
	}

	if isoImageID == "" {
		return nil, newError(EBadArgument, "ISO image ID cannot be empty")
	}

	cdrom.fileID = isoImageID
	cdrom.fileName = fmt.Sprintf("iso-%s.iso", isoImageID)

	return cdrom, nil
}
