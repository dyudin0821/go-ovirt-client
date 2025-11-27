package ovirtclient_test

import (
	"testing"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

func TestStorageDomainListFiles(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	storageDomainID := helper.GetStorageDomainID()

	// List files (should be empty initially in mock)
	files, err := client.ListStorageDomainFiles(storageDomainID, false)
	if err != nil {
		t.Fatalf("Failed to list storage domain files (%v)", err)
	}

	if len(files) != 0 {
		t.Fatalf("Expected 0 files, got %d", len(files))
	}
}

func TestStorageDomainGetFileNotFound(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	storageDomainID := helper.GetStorageDomainID()

	// Try to get a non-existent file
	_, err = client.GetStorageDomainFile(storageDomainID, ovirtclient.FileID("non-existent"))
	if err == nil {
		t.Fatal("Expected error when getting non-existent file, got nil")
	}

	if !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		t.Fatalf("Expected ENotFound error, got: %v", err)
	}
}

func TestStorageDomainGetFileInvalidStorageDomain(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()

	// Try to get a file from a non-existent storage domain
	_, err = client.GetStorageDomainFile(ovirtclient.StorageDomainID("invalid"), ovirtclient.FileID("test"))
	if err == nil {
		t.Fatal("Expected error when using invalid storage domain, got nil")
	}

	if !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		t.Fatalf("Expected ENotFound error, got: %v", err)
	}
}

func TestStorageDomainListFilesInvalidStorageDomain(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()

	// Try to list files from a non-existent storage domain
	_, err = client.ListStorageDomainFiles(ovirtclient.StorageDomainID("invalid"), false)
	if err == nil {
		t.Fatal("Expected error when using invalid storage domain, got nil")
	}

	if !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		t.Fatalf("Expected ENotFound error, got: %v", err)
	}
}
