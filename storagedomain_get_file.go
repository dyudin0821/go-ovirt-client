package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetStorageDomainFile(
	storageDomainID StorageDomainID,
	fileID FileID,
	retries ...RetryStrategy,
) (result File, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting file %s from storage domain %s", fileID, storageDomainID),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.
				SystemService().
				StorageDomainsService().
				StorageDomainService(string(storageDomainID)).
				FilesService().
				FileService(string(fileID)).
				Get().
				Send()
			if e != nil {
				return e
			}
			sdkFile, ok := response.File()
			if !ok {
				return newError(
					ENotFound,
					"no file returned when getting file ID %s from storage domain %s",
					fileID,
					storageDomainID,
				)
			}
			result, err = convertSDKFile(sdkFile, storageDomainID, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert file %s from storage domain %s",
					fileID,
					storageDomainID,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetStorageDomainFile(
	storageDomainID StorageDomainID,
	fileID FileID,
	_ ...RetryStrategy,
) (File, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.storageDomains[storageDomainID]; !ok {
		return nil, newError(ENotFound, "storage domain with ID %s not found", storageDomainID)
	}

	files, ok := m.storageDomainFiles[storageDomainID]
	if !ok {
		return nil, newError(ENotFound, "file %s not found in storage domain %s", fileID, storageDomainID)
	}

	file, ok := files[fileID]
	if !ok {
		return nil, newError(ENotFound, "file %s not found in storage domain %s", fileID, storageDomainID)
	}

	return file, nil
}
