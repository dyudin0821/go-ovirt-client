package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListCDROMs(vmID VMID, retries ...RetryStrategy) (result []CDROM, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing CDROMs for VM %s", vmID),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().VmsService().VmService(string(vmID)).CdromsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Cdroms()
			if !ok {
				return nil
			}
			result = make([]CDROM, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKCDROM(sdkObject, vmID, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert CDROM during listing item #%d", i)
				}
			}
			return nil
		},
	)
	return
}

func (m *mockClient) ListCDROMs(vmID VMID, _ ...RetryStrategy) ([]CDROM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var result []CDROM
	if m.vmCDROMsByVM != nil && m.vmCDROMsByVM[vmID] != nil {
		for _, item := range m.vmCDROMsByVM[vmID] {
			result = append(result, item)
		}
	}
	return result, nil
}
