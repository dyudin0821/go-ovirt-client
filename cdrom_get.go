package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetCDROM(vmID VMID, id CDROMID, retries ...RetryStrategy) (result CDROM, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting CDROM %s for VM %s", id, vmID),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().VmsService().VmService(string(vmID)).CdromsService().
				CdromService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Cdrom()
			if !ok {
				return newError(
					ENotFound,
					"no CDROM returned when getting CDROM %s on VM %s",
					id,
					vmID,
				)
			}
			result, err = convertSDKCDROM(sdkObject, vmID, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert CDROM %s on VM %s",
					id,
					vmID,
				)
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) GetCDROM(vmID VMID, id CDROMID, _ ...RetryStrategy) (CDROM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.vmCDROMsByVM != nil && m.vmCDROMsByVM[vmID] != nil {
		if cdrom, ok := m.vmCDROMsByVM[vmID][id]; ok {
			return cdrom, nil
		}
	}
	return nil, newError(ENotFound, "CDROM with ID %s not found for VM %s", id, vmID)
}
