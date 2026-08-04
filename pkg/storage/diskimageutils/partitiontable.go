package diskimageutils

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	mbrSignature       = 0xAA55
	mbrPartitionOffset = 446
	mbrPartitionSize   = 16
	mbrMaxPartitions   = 4
	sectorSize         = 512

	gptSignature = "EFI PART"
	gptHeaderLBA = 1
)

type PartitionEntry struct {
	Index      int
	StartLBA   uint64
	SizeLBA    uint64
	SizeBytes  uint64
	TypeCode   byte
	Bootable   bool
	FileSystem string
}

type PartitionTable struct {
	DiskImagePath string
	TableType     string // "MBR" or "GPT"
	Partitions    []*PartitionEntry
}

// ReadPartitionTable reads all partition tables from a disk image at 'path'.
//
// It first checks for a GPT header. If none is found, it falls back to reading the MBR partition table.
func ReadPartitionTable(ctx context.Context, path string) (*PartitionTable, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Reading partition table from disk image '%s' started.", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error opening disk image '%s': %w", path, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error getting file info for '%s': %w", path, err)
	}

	logging.LogInfoByCtxf(ctx, "Disk image size: %d bytes.", fileInfo.Size())

	partitionTable, err := tryReadGPT(ctx, file, path)
	if err == nil && partitionTable != nil {
		logging.LogInfoByCtxf(ctx, "GPT partition table detected with '%d' partitions.", len(partitionTable.Partitions))
		logging.LogInfoByCtxf(ctx, "Reading partition table from disk image '%s' finished.", path)
		return partitionTable, nil
	}

	logging.LogInfoByCtxf(ctx, "No GPT header found, falling back to MBR.")

	partitionTable, err = readMBR(ctx, file, path)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "MBR partition table detected with '%d' partitions.", len(partitionTable.Partitions))
	logging.LogInfoByCtxf(ctx, "Reading partition table from disk image '%s' finished.", path)

	return partitionTable, nil
}

func tryReadGPT(ctx context.Context, file *os.File, path string) (*PartitionTable, error) {
	headerBuf := make([]byte, sectorSize)

	_, err := file.ReadAt(headerBuf, int64(gptHeaderLBA)*sectorSize)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error reading GPT header from '%s': %w", path, err)
	}

	signature := string(headerBuf[0:8])
	if signature != gptSignature {
		return nil, tracederrors.TracedErrorf("No GPT signature found in '%s'.", path)
	}

	partitionEntryLBA := binary.LittleEndian.Uint64(headerBuf[72:80])
	numberOfEntries := binary.LittleEndian.Uint32(headerBuf[80:84])
	entrySize := binary.LittleEndian.Uint32(headerBuf[84:88])

	logging.LogInfoByCtxf(ctx, "GPT header found: %d partition entries of size %d bytes starting at LBA %d.", numberOfEntries, entrySize, partitionEntryLBA)

	partitions := []*PartitionEntry{}

	for i := uint32(0); i < numberOfEntries; i++ {
		offset := int64(partitionEntryLBA)*sectorSize + int64(i)*int64(entrySize)
		entryBuf := make([]byte, entrySize)

		_, err := file.ReadAt(entryBuf, offset)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Error reading GPT entry %d from '%s': %w", i, path, err)
		}

		// Check if entry is unused (all zeros in type GUID)
		allZero := true
		for _, b := range entryBuf[0:16] {
			if b != 0 {
				allZero = false
				break
			}
		}

		if allZero {
			continue
		}

		startLBA := binary.LittleEndian.Uint64(entryBuf[32:40])
		endLBA := binary.LittleEndian.Uint64(entryBuf[40:48])
		sizeLBA := endLBA - startLBA + 1

		partition := &PartitionEntry{
			Index:     int(i) + 1,
			StartLBA:  startLBA,
			SizeLBA:   sizeLBA,
			SizeBytes: sizeLBA * sectorSize,
		}

		partitions = append(partitions, partition)
	}

	return &PartitionTable{
		DiskImagePath: path,
		TableType:     "GPT",
		Partitions:    partitions,
	}, nil
}

func readMBR(ctx context.Context, file *os.File, path string) (*PartitionTable, error) {
	mbrBuf := make([]byte, sectorSize)

	_, err := file.ReadAt(mbrBuf, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Error reading MBR from '%s': %w", path, err)
	}

	sig := binary.LittleEndian.Uint16(mbrBuf[510:512])
	if sig != mbrSignature {
		return nil, tracederrors.TracedErrorf("Invalid MBR signature in '%s': expected 0x%X, got 0x%X.", path, mbrSignature, sig)
	}

	partitions := []*PartitionEntry{}

	for i := 0; i < mbrMaxPartitions; i++ {
		offset := mbrPartitionOffset + i*mbrPartitionSize
		entry := mbrBuf[offset : offset+mbrPartitionSize]

		typeCode := entry[4]
		if typeCode == 0x00 {
			continue
		}

		bootable := entry[0] == 0x80
		startLBA := uint64(binary.LittleEndian.Uint32(entry[8:12]))
		sizeLBA := uint64(binary.LittleEndian.Uint32(entry[12:16]))

		partition := &PartitionEntry{
			Index:      i + 1,
			StartLBA:   startLBA,
			SizeLBA:    sizeLBA,
			SizeBytes:  sizeLBA * sectorSize,
			TypeCode:   typeCode,
			Bootable:   bootable,
			FileSystem: getMBRPartitionTypeName(typeCode),
		}

		partitions = append(partitions, partition)

		logging.LogInfoByCtxf(ctx, "MBR partition %d: type=0x%02X (%s), startLBA=%d, sizeLBA=%d, bootable=%t.",
			partition.Index, typeCode, partition.FileSystem, startLBA, sizeLBA, bootable)
	}

	return &PartitionTable{
		DiskImagePath: path,
		TableType:     "MBR",
		Partitions:    partitions,
	}, nil
}

func getMBRPartitionTypeName(typeCode byte) string {
	switch typeCode {
	case 0x01:
		return "FAT12"
	case 0x04:
		return "FAT16 (<32MB)"
	case 0x05:
		return "Extended"
	case 0x06:
		return "FAT16 (>32MB)"
	case 0x07:
		return "NTFS/exFAT"
	case 0x0B:
		return "FAT32 (CHS)"
	case 0x0C:
		return "FAT32 (LBA)"
	case 0x0E:
		return "FAT16 (LBA)"
	case 0x0F:
		return "Extended (LBA)"
	case 0x82:
		return "Linux Swap"
	case 0x83:
		return "Linux"
	case 0x8E:
		return "Linux LVM"
	case 0xEE:
		return "GPT Protective"
	case 0xEF:
		return "EFI System"
	case 0xFD:
		return "Linux RAID"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", typeCode)
	}
}
