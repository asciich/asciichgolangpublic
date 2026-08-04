package diskimageutils_test

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/diskimageutils"
)

const (
	testMBRSignature       = 0xAA55
	testMBRPartitionOffset = 446
	testMBRPartitionSize   = 16
	testMBRMaxPartitions   = 4
	testSectorSize         = 512
	testGPTHeaderLBA       = 1
	testGPTSignature       = "EFI PART"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestReadPartitionTable_EmptyPath(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.ReadPartitionTable(ctx, "")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path")
}

func TestReadPartitionTable_NonExistentFile(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.ReadPartitionTable(ctx, "/tmp/nonexistent_disk_image_12345.img")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestReadPartitionTable_MBR_SinglePartition(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "mbr_single.img")
	createMBRDiskImage(t, path, []mbrTestPartition{
		{
			typeCode: 0x83,
			bootable: true,
			startLBA: 2048,
			sizeLBA:  1048576,
		},
	})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "MBR", result.TableType)
	assert.Equal(t, path, result.DiskImagePath)
	assert.Len(t, result.Partitions, 1)

	p := result.Partitions[0]
	assert.Equal(t, 1, p.Index)
	assert.Equal(t, uint64(2048), p.StartLBA)
	assert.Equal(t, uint64(1048576), p.SizeLBA)
	assert.Equal(t, uint64(1048576*512), p.SizeBytes)
	assert.Equal(t, byte(0x83), p.TypeCode)
	assert.True(t, p.Bootable)
	assert.Equal(t, "Linux", p.FileSystem)
}

func TestReadPartitionTable_MBR_MultiplePartitions(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "mbr_multi.img")
	createMBRDiskImage(t, path, []mbrTestPartition{
		{
			typeCode: 0x0C,
			bootable: true,
			startLBA: 2048,
			sizeLBA:  204800,
		},
		{
			typeCode: 0x83,
			bootable: false,
			startLBA: 206848,
			sizeLBA:  1048576,
		},
		{
			typeCode: 0x82,
			bootable: false,
			startLBA: 1255424,
			sizeLBA:  524288,
		},
	})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "MBR", result.TableType)
	assert.Len(t, result.Partitions, 3)

	assert.Equal(t, byte(0x0C), result.Partitions[0].TypeCode)
	assert.Equal(t, "FAT32 (LBA)", result.Partitions[0].FileSystem)
	assert.True(t, result.Partitions[0].Bootable)

	assert.Equal(t, byte(0x83), result.Partitions[1].TypeCode)
	assert.Equal(t, "Linux", result.Partitions[1].FileSystem)
	assert.False(t, result.Partitions[1].Bootable)

	assert.Equal(t, byte(0x82), result.Partitions[2].TypeCode)
	assert.Equal(t, "Linux Swap", result.Partitions[2].FileSystem)
	assert.False(t, result.Partitions[2].Bootable)
}

func TestReadPartitionTable_MBR_NoPartitions(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "mbr_empty.img")
	createMBRDiskImage(t, path, []mbrTestPartition{})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "MBR", result.TableType)
	assert.Len(t, result.Partitions, 0)
}

func TestReadPartitionTable_MBR_InvalidSignature(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "invalid.img")
	buf := make([]byte, 512)
	binary.LittleEndian.PutUint16(buf[510:512], 0x0000)

	err := os.WriteFile(path, buf, 0644)
	require.NoError(t, err)

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature")
}

func TestReadPartitionTable_GPT_SinglePartition(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "gpt_single.img")
	createGPTDiskImage(t, path, []gptTestPartition{
		{
			startLBA: 2048,
			endLBA:   1050623,
			typeGUID: linuxFilesystemGUID(),
		},
	})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "GPT", result.TableType)
	assert.Equal(t, path, result.DiskImagePath)
	assert.Len(t, result.Partitions, 1)

	p := result.Partitions[0]
	assert.Equal(t, 1, p.Index)
	assert.Equal(t, uint64(2048), p.StartLBA)
	assert.Equal(t, uint64(1050623-2048+1), p.SizeLBA)
	assert.Equal(t, uint64((1050623-2048+1)*512), p.SizeBytes)
}

func TestReadPartitionTable_GPT_MultiplePartitions(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "gpt_multi.img")
	createGPTDiskImage(t, path, []gptTestPartition{
		{
			startLBA: 2048,
			endLBA:   206847,
			typeGUID: efiSystemPartitionGUID(),
		},
		{
			startLBA: 206848,
			endLBA:   1255423,
			typeGUID: linuxFilesystemGUID(),
		},
		{
			startLBA: 1255424,
			endLBA:   1779711,
			typeGUID: linuxSwapGUID(),
		},
	})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "GPT", result.TableType)
	assert.Len(t, result.Partitions, 3)

	assert.Equal(t, uint64(2048), result.Partitions[0].StartLBA)
	assert.Equal(t, uint64(206848), result.Partitions[1].StartLBA)
	assert.Equal(t, uint64(1255424), result.Partitions[2].StartLBA)
}

func TestReadPartitionTable_GPT_NoPartitions(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "gpt_empty.img")
	createGPTDiskImage(t, path, []gptTestPartition{})

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "GPT", result.TableType)
	assert.Len(t, result.Partitions, 0)
}

func TestReadPartitionTable_FileTooSmall(t *testing.T) {
	ctx := getCtx()

	path := filepath.Join(t.TempDir(), "tiny.img")
	err := os.WriteFile(path, []byte{0x00, 0x01, 0x02}, 0644)
	require.NoError(t, err)

	result, err := diskimageutils.ReadPartitionTable(ctx, path)
	assert.Nil(t, result)
	assert.Error(t, err)
}

// --- Test Helpers ---

type mbrTestPartition struct {
	typeCode byte
	bootable bool
	startLBA uint32
	sizeLBA  uint32
}

func createMBRDiskImage(t *testing.T, path string, partitions []mbrTestPartition) {
	t.Helper()

	buf := make([]byte, 512)

	for i, p := range partitions {
		if i >= testMBRMaxPartitions {
			break
		}

		offset := testMBRPartitionOffset + i*testMBRPartitionSize

		if p.bootable {
			buf[offset] = 0x80
		} else {
			buf[offset] = 0x00
		}

		buf[offset+4] = p.typeCode
		binary.LittleEndian.PutUint32(buf[offset+8:offset+12], p.startLBA)
		binary.LittleEndian.PutUint32(buf[offset+12:offset+16], p.sizeLBA)
	}

	binary.LittleEndian.PutUint16(buf[510:512], testMBRSignature)

	err := os.WriteFile(path, buf, 0644)
	require.NoError(t, err)
}

type gptTestPartition struct {
	startLBA uint64
	endLBA   uint64
	typeGUID [16]byte
}

func createGPTDiskImage(t *testing.T, path string, partitions []gptTestPartition) {
	t.Helper()

	const (
		gptEntrySize      = 128
		gptMaxEntries     = 128
		partitionTableLBA = 2
	)

	totalSize := (partitionTableLBA + (gptMaxEntries*gptEntrySize+testSectorSize-1)/testSectorSize) * testSectorSize
	buf := make([]byte, totalSize)

	// Protective MBR
	buf[testMBRPartitionOffset+4] = 0xEE
	binary.LittleEndian.PutUint32(buf[testMBRPartitionOffset+8:testMBRPartitionOffset+12], 1)
	binary.LittleEndian.PutUint32(buf[testMBRPartitionOffset+12:testMBRPartitionOffset+16], 0xFFFFFFFF)
	binary.LittleEndian.PutUint16(buf[510:512], testMBRSignature)

	// GPT Header at LBA 1
	headerOffset := testGPTHeaderLBA * testSectorSize
	copy(buf[headerOffset:headerOffset+8], []byte(testGPTSignature))

	// Revision
	binary.LittleEndian.PutUint32(buf[headerOffset+8:headerOffset+12], 0x00010000)
	// Header size
	binary.LittleEndian.PutUint32(buf[headerOffset+12:headerOffset+16], 92)
	// Partition entry start LBA
	binary.LittleEndian.PutUint64(buf[headerOffset+72:headerOffset+80], partitionTableLBA)
	// Number of partition entries
	binary.LittleEndian.PutUint32(buf[headerOffset+80:headerOffset+84], gptMaxEntries)
	// Size of partition entry
	binary.LittleEndian.PutUint32(buf[headerOffset+84:headerOffset+88], gptEntrySize)

	// Write partition entries
	for i, p := range partitions {
		entryOffset := partitionTableLBA*testSectorSize + i*gptEntrySize

		copy(buf[entryOffset:entryOffset+16], p.typeGUID[:])

		// Unique partition GUID (just needs to be non-zero)
		buf[entryOffset+16] = byte(i + 1)

		binary.LittleEndian.PutUint64(buf[entryOffset+32:entryOffset+40], p.startLBA)
		binary.LittleEndian.PutUint64(buf[entryOffset+40:entryOffset+48], p.endLBA)
	}

	err := os.WriteFile(path, buf, 0644)
	require.NoError(t, err)
}

func linuxFilesystemGUID() [16]byte {
	// 0FC63DAF-8483-4772-8E79-3D69D8477DE4 in mixed-endian
	return [16]byte{
		0xAF, 0x3D, 0xC6, 0x0F, 0x83, 0x84, 0x72, 0x47,
		0x8E, 0x79, 0x3D, 0x69, 0xD8, 0x47, 0x7D, 0xE4,
	}
}

func efiSystemPartitionGUID() [16]byte {
	// C12A7328-F81F-11D2-BA4B-00A0C93EC93B in mixed-endian
	return [16]byte{
		0x28, 0x73, 0x2A, 0xC1, 0x1F, 0xF8, 0xD2, 0x11,
		0xBA, 0x4B, 0x00, 0xA0, 0xC9, 0x3E, 0xC9, 0x3B,
	}
}

func linuxSwapGUID() [16]byte {
	// 0657FD6D-A4AB-43C4-84E5-0933C84B4F4F in mixed-endian
	return [16]byte{
		0x6D, 0xFD, 0x57, 0x06, 0xAB, 0xA4, 0xC4, 0x43,
		0x84, 0xE5, 0x09, 0x33, 0xC8, 0x4B, 0x4F, 0x4F,
	}
}
