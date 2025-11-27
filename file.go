package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// FileID is the identifier for a file in a storage domain (ISO image or VFD).
type FileID string

// File represents an ISO image or virtual floppy disk file in a storage domain.
type File interface {
	// ID returns the unique identifier of the file.
	ID() FileID
	// Name returns the name of the file.
	Name() string
	// Type returns the type of the file (e.g., "iso").
	Type() string
	// StorageDomainID returns the ID of the storage domain containing this file.
	StorageDomainID() StorageDomainID
	// Comment returns the comment associated with the file, if any.
	Comment() string
	// Description returns the description of the file, if any.
	Description() string

	// StorageDomain fetches the storage domain this file belongs to.
	StorageDomain(retries ...RetryStrategy) (StorageDomain, error)
}

// FileList is a list of files.
type FileList []File

type file struct {
	client Client

	id              FileID
	name            string
	fileType        string
	storageDomainID StorageDomainID
	comment         string
	description     string
}

func (f *file) ID() FileID {
	return f.id
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Type() string {
	return f.fileType
}

func (f *file) StorageDomainID() StorageDomainID {
	return f.storageDomainID
}

func (f *file) Comment() string {
	return f.comment
}

func (f *file) Description() string {
	return f.description
}

func (f *file) StorageDomain(retries ...RetryStrategy) (StorageDomain, error) {
	return f.client.GetStorageDomain(f.storageDomainID, retries...)
}

func convertSDKFile(object *ovirtsdk4.File, storageDomainID StorageDomainID, client Client) (File, error) {
	id, ok := object.Id()
	if !ok {
		return nil, newFieldNotFound("file", "id")
	}

	name, ok := object.Name()
	if !ok {
		return nil, newFieldNotFound("file", "name")
	}

	result := &file{
		client:          client,
		id:              FileID(id),
		name:            name,
		storageDomainID: storageDomainID,
	}

	// Optional fields
	if fileType, ok := object.Type(); ok {
		result.fileType = fileType
	}

	if comment, ok := object.Comment(); ok {
		result.comment = comment
	}

	if description, ok := object.Description(); ok {
		result.description = description
	}

	return result, nil
}
