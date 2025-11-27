package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListStorageDomainFiles(
	storageDomainID StorageDomainID,
	refresh bool,
	retries ...RetryStrategy,
) (result FileList, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing files in storage domain %s", storageDomainID),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.
				SystemService().
				StorageDomainsService().
				StorageDomainService(string(storageDomainID)).
				FilesService().
				List().
				Refresh(refresh).
				Send()
			if e != nil {
				return e
			}
			sdkFiles, ok := response.File()
			if !ok {
				return nil
			}
			result = make(FileList, len(sdkFiles.Slice()))
			for i, sdkFile := range sdkFiles.Slice() {
				result[i], e = convertSDKFile(sdkFile, storageDomainID, o)
				if e != nil {
					return wrap(
						e,
						EBug,
						"failed to convert file during listing item #%d in storage domain %s",
						i,
						storageDomainID,
					)
				}
			}
			return nil
		})
	return result, err
}

func (m *mockClient) ListStorageDomainFiles(
	storageDomainID StorageDomainID,
	_ bool,
	_ ...RetryStrategy,
) (FileList, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.storageDomains[storageDomainID]; !ok {
		return nil, newError(ENotFound, "storage domain with ID %s not found", storageDomainID)
	}

	files, ok := m.storageDomainFiles[storageDomainID]
	if !ok {
		return FileList{}, nil
	}

	result := make(FileList, len(files))
	i := 0
	for _, file := range files {
		result[i] = file
		i++
	}
	return result, nil
}
