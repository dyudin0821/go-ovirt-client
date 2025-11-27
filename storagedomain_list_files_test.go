package ovirtclient_test

import (
	"testing"

	ovirtclient "github.com/dyudin0821/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

func TestListStorageDomainFiles(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	storageDomainID := helper.GetStorageDomainID()

	// List files without refresh - should be empty initially
	files, err := client.ListStorageDomainFiles(storageDomainID, false)
	if err != nil {
		t.Fatalf("failed to list storage domain files (%v)", err)
	}

	// Initially should be empty
	if len(files) != 0 {
		t.Logf("storage domain has %d files initially", len(files))
	}

	// List with refresh should also work
	files, err = client.ListStorageDomainFiles(storageDomainID, true)
	if err != nil {
		t.Fatalf("failed to list storage domain files with refresh (%v)", err)
	}

	// Verify we got a list (even if empty)
	if files == nil {
		t.Fatal("expected non-nil file list")
	}
}

func TestListStorageDomainFilesEmpty(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	storageDomainID := helper.GetStorageDomainID()

	// List files from storage domain with no files
	files, err := client.ListStorageDomainFiles(storageDomainID, false)
	if err != nil {
		t.Fatalf("failed to list storage domain files (%v)", err)
	}

	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestListStorageDomainFilesNotFound(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	invalidStorageDomainID := ovirtclient.StorageDomainID("invalid-storage-domain-id")

	// Try to list files from non-existent storage domain
	_, err = client.ListStorageDomainFiles(invalidStorageDomainID, false)
	if err == nil {
		t.Fatal("expected error when listing files from non-existent storage domain, got nil")
	}

	if !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		t.Fatalf("expected ENotFound error, got %v", err)
	}
}

func TestGetStorageDomainFileNotFound(t *testing.T) {
	t.Parallel()
	helper, err := ovirtclient.NewMockTestHelper(ovirtclientlog.NewTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create mock test helper (%v)", err)
	}

	client := helper.GetClient()
	storageDomainID := helper.GetStorageDomainID()
	invalidFileID := ovirtclient.FileID("invalid-file-id")

	// Try to get non-existent file
	_, err = client.GetStorageDomainFile(storageDomainID, invalidFileID)
	if err == nil {
		t.Fatal("expected error when getting non-existent file, got nil")
	}

	if !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		t.Fatalf("expected ENotFound error, got %v", err)
	}
}
