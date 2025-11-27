package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) EjectCDROM(
	vmID VMID,
	cdromID CDROMID,
	retries ...RetryStrategy,
) (result CDROM, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("ejecting CDROM %s from VM %s", cdromID, vmID),
		o.logger,
		retries,
		func() error {
			cdromBuilder := ovirtsdk.NewCdromBuilder()
			// Setting file to empty ejects the ISO
			cdromBuilder.File(ovirtsdk.NewFileBuilder().MustBuild())

			updateRequest := o.conn.SystemService().VmsService().VmService(string(vmID)).
				CdromsService().CdromService(string(cdromID)).Update()
			updateRequest.Cdrom(cdromBuilder.MustBuild())
			response, err := updateRequest.Send()
			if err != nil {
				return wrap(
					err,
					EUnidentified,
					"failed to eject CDROM %s from VM %s",
					cdromID,
					vmID,
				)
			}

			cdrom, ok := response.Cdrom()
			if !ok {
				return newFieldNotFound("CDROM eject response", "cdrom")
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

func (m *mockClient) EjectCDROM(
	vmID VMID,
	cdromID CDROMID,
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

	// Eject by clearing the file reference
	cdrom.fileID = ""
	cdrom.fileName = ""

	return cdrom, nil
}
