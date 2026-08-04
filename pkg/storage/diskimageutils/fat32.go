package diskimageutils

import (
	"context"
	"encoding/binary"
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	fat32SectorSize        = 512
	fat32SectorsPerCluster = 8 // 4KB clusters
	fat32ReservedSectors   = 32
	fat32NumberOfFATs      = 2
	fat32MediaType         = 0xF8
	fat32FSInfoSectorNum   = 1
	fat32BackupBootSector  = 6
	fat32RootDirCluster    = 2
	fat32EOFMarker         = 0x0FFFFFF8
	fat32ClusterInUse      = 0x0FFFFFFF
	fat32BootSignature     = 0xAA55
	fat32FSInfoSignature1  = 0x41615252
	fat32FSInfoSignature2  = 0x61417272
	fat32FSInfoSignature3  = 0xAA550000
	fat32MinImageSize      = 33 * 1024 * 1024
)

// CreateFat32Options configures the FAT32 image creation.
type CreateFat32Options struct {
	// SizeBytes is the total size of the image in bytes.
	// Must be at least 33 MB for a valid FAT32 filesystem.
	SizeBytes int64

	// VolumeLabel is the optional volume label (max 11 characters, padded with spaces).
	VolumeLabel string
}

// CreateFat32Image creates a FAT32 formatted raw image file at 'outputPath' without a partition table.
//
// The image is a bare FAT32 filesystem that can be mounted directly (e.g. mount -o loop).
func CreateFat32Image(ctx context.Context, outputPath string, options *CreateFat32Options) error {
	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.SizeBytes < fat32MinImageSize {
		return tracederrors.TracedErrorf("FAT32 requires at least 33 MB, got '%d' bytes.", options.SizeBytes)
	}

	logging.LogInfoByCtxf(ctx, "Create FAT32 image '%s' with size '%d' bytes started.", outputPath, options.SizeBytes)

	totalSectors := uint32(options.SizeBytes / fat32SectorSize)
	fatSizeSectors := fat32CalculateFATSize(totalSectors)

	dataStartSector := uint32(fat32ReservedSectors) + fat32NumberOfFATs*fatSizeSectors
	dataSectors := totalSectors - dataStartSector
	totalClusters := dataSectors / fat32SectorsPerCluster

	volumeLabel := fat32PadVolumeLabel(options.VolumeLabel)

	file, err := os.Create(outputPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create image file '%s': %w", outputPath, err)
	}
	defer file.Close()

	err = file.Truncate(options.SizeBytes)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to allocate image file '%s': %w", outputPath, err)
	}

	bootSector := fat32BuildBootSector(totalSectors, fatSizeSectors, volumeLabel)
	_, err = file.WriteAt(bootSector, 0)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to write boot sector: %w", err)
	}

	_, err = file.WriteAt(bootSector, int64(fat32BackupBootSector)*fat32SectorSize)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to write backup boot sector: %w", err)
	}

	fsInfo := fat32BuildFSInfoSector(totalClusters)
	_, err = file.WriteAt(fsInfo, int64(fat32FSInfoSectorNum)*fat32SectorSize)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to write FSInfo sector: %w", err)
	}

	fatTable := fat32BuildFATTable(fatSizeSectors)
	for i := uint32(0); i < fat32NumberOfFATs; i++ {
		fatOffset := int64(fat32ReservedSectors+i*fatSizeSectors) * fat32SectorSize
		_, err = file.WriteAt(fatTable, fatOffset)
		if err != nil {
			os.Remove(outputPath)
			return tracederrors.TracedErrorf("Failed to write FAT table %d: %w", i+1, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "FAT32 image created: totalSectors=%d, fatSizeSectors=%d, totalClusters=%d, dataStartSector=%d.",
		totalSectors, fatSizeSectors, totalClusters, dataStartSector)
	logging.LogInfoByCtxf(ctx, "Create FAT32 image '%s' with size '%d' bytes finished.", outputPath, options.SizeBytes)

	return nil
}

// Fat32ListFiles lists all file and directory paths in a FAT32 image at 'imagePath'.
//
// Returns paths in a find-like format starting with ".".
func Fat32ListFiles(ctx context.Context, imagePath string) ([]string, error) {
	if imagePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("imagePath")
	}

	logging.LogInfoByCtxf(ctx, "List files in FAT32 image '%s' started.", imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open FAT32 image '%s': %w", imagePath, err)
	}
	defer file.Close()

	bpb, err := fat32ReadBPB(file)
	if err != nil {
		return nil, err
	}

	fat, err := fat32ReadFAT(file, bpb)
	if err != nil {
		return nil, err
	}

	results := []string{"."}

	err = fat32ListDirectory(file, bpb, fat, bpb.rootCluster, ".", &results)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List files in FAT32 image '%s' finished. Found '%d' entries.", imagePath, len(results))

	return results, nil
}

// Fat32WriteFile writes the given content bytes to the specified path inside a FAT32 image.
//
// The 'filePath' must use forward slashes as separator (e.g. "test/a.txt").
// Parent directories are created automatically if they do not exist.
func Fat32WriteFile(ctx context.Context, imagePath string, filePath string, content []byte) error {
	if imagePath == "" {
		return tracederrors.TracedErrorEmptyString("imagePath")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	logging.LogInfoByCtxf(ctx, "Write file '%s' to FAT32 image '%s' started. Size: '%d' bytes.", filePath, imagePath, len(content))

	file, err := os.OpenFile(imagePath, os.O_RDWR, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open FAT32 image '%s': %w", imagePath, err)
	}
	defer file.Close()

	bpb, err := fat32ReadBPB(file)
	if err != nil {
		return err
	}

	fat, err := fat32ReadFAT(file, bpb)
	if err != nil {
		return err
	}

	parts := fat32SplitPath(filePath)
	if len(parts) == 0 {
		return tracederrors.TracedErrorf("Invalid file path '%s'.", filePath)
	}

	dirCluster := bpb.rootCluster
	for i := 0; i < len(parts)-1; i++ {
		existingCluster, found, err := fat32FindEntryInDirectory(file, bpb, fat, dirCluster, parts[i], true)
		if err != nil {
			return err
		}

		if found {
			dirCluster = existingCluster
		} else {
			newCluster, err := fat32AllocateCluster(fat)
			if err != nil {
				return err
			}

			err = fat32ZeroCluster(file, bpb, newCluster)
			if err != nil {
				return err
			}

			err = fat32AddDirectoryEntry(file, bpb, fat, dirCluster, parts[i], newCluster, 0, true)
			if err != nil {
				return err
			}

			dirCluster = newCluster
		}
	}

	fileName := parts[len(parts)-1]
	fileSize := uint32(len(content))

	var firstCluster uint32
	if fileSize > 0 {
		firstCluster, err = fat32WriteContent(file, bpb, fat, content)
		if err != nil {
			return err
		}
	}

	err = fat32AddDirectoryEntry(file, bpb, fat, dirCluster, fileName, firstCluster, fileSize, false)
	if err != nil {
		return err
	}

	err = fat32WriteFAT(file, bpb, fat)
	if err != nil {
		return err
	}

	err = fat32UpdateFSInfo(file, bpb, fat)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Write file '%s' to FAT32 image '%s' finished.", filePath, imagePath)

	return nil
}

// Fat32ReadFile reads the content of a file at 'filePath' from a FAT32 image.
//
// The 'filePath' must use forward slashes as separator (e.g. "test/a.txt").
func Fat32ReadFile(ctx context.Context, imagePath string, filePath string) ([]byte, error) {
	if imagePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("imagePath")
	}

	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	logging.LogInfoByCtxf(ctx, "Read file '%s' from FAT32 image '%s' started.", filePath, imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open FAT32 image '%s': %w", imagePath, err)
	}
	defer file.Close()

	bpb, err := fat32ReadBPB(file)
	if err != nil {
		return nil, err
	}

	fat, err := fat32ReadFAT(file, bpb)
	if err != nil {
		return nil, err
	}

	parts := fat32SplitPath(filePath)
	if len(parts) == 0 {
		return nil, tracederrors.TracedErrorf("Invalid file path '%s'.", filePath)
	}

	dirCluster := bpb.rootCluster
	for i := 0; i < len(parts)-1; i++ {
		cluster, found, err := fat32FindEntryInDirectory(file, bpb, fat, dirCluster, parts[i], true)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, tracederrors.TracedErrorf("Directory '%s' not found in FAT32 image '%s'.", parts[i], imagePath)
		}
		dirCluster = cluster
	}

	fileName := parts[len(parts)-1]
	fileEntry, err := fat32FindFileEntry(file, bpb, fat, dirCluster, fileName)
	if err != nil {
		return nil, err
	}

	if fileEntry == nil {
		return nil, tracederrors.TracedErrorf("File '%s' not found in FAT32 image '%s'.", filePath, imagePath)
	}

	if fileEntry.size == 0 {
		logging.LogInfoByCtxf(ctx, "Read file '%s' from FAT32 image '%s' finished. Size: 0 bytes.", filePath, imagePath)
		return []byte{}, nil
	}

	clusterChain := fat32GetClusterChain(fat, fileEntry.cluster)
	clusterBytes := fat32GetClusterSize(bpb)
	content := make([]byte, 0, fileEntry.size)

	remaining := int(fileEntry.size)
	for _, cluster := range clusterChain {
		offset := fat32ClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to read cluster %d: %w", cluster, err)
		}

		readSize := clusterBytes
		if remaining < readSize {
			readSize = remaining
		}

		content = append(content, buf[:readSize]...)
		remaining -= readSize

		if remaining <= 0 {
			break
		}
	}

	logging.LogInfoByCtxf(ctx, "Read file '%s' from FAT32 image '%s' finished. Size: '%d' bytes.", filePath, imagePath, len(content))

	return content, nil
}

// Fat32Exists checks if a file or directory exists at 'path' in a FAT32 image.
func Fat32Exists(ctx context.Context, imagePath string, path string) (bool, error) {
	if imagePath == "" {
		return false, tracederrors.TracedErrorEmptyString("imagePath")
	}

	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	files, err := Fat32ListFiles(ctx, imagePath)
	if err != nil {
		return false, err
	}

	searchPath := "./" + strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")
	for _, f := range files {
		if strings.EqualFold(f, searchPath) {
			return true, nil
		}
	}

	return false, nil
}

// --- Internal types ---

type fat32BPB struct {
	sectorSize      uint16
	sectorsPerClust uint8
	reservedSectors uint16
	numberOfFATs    uint8
	fatSizeSectors  uint32
	rootCluster     uint32
	totalSectors    uint32
	dataStartSector uint32
	totalClusters   uint32
	fsInfoSector    uint16
}

type fat32DirEntry struct {
	name        string
	cluster     uint32
	size        uint32
	isDirectory bool
}

// --- Internal functions ---

func fat32ReadBPB(file *os.File) (*fat32BPB, error) {
	bs := make([]byte, fat32SectorSize)
	_, err := file.ReadAt(bs, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read boot sector: %w", err)
	}

	sig := binary.LittleEndian.Uint16(bs[510:512])
	if sig != fat32BootSignature {
		return nil, tracederrors.TracedErrorf("Invalid boot sector signature: 0x%04X.", sig)
	}

	bpb := &fat32BPB{
		sectorSize:      binary.LittleEndian.Uint16(bs[11:13]),
		sectorsPerClust: bs[13],
		reservedSectors: binary.LittleEndian.Uint16(bs[14:16]),
		numberOfFATs:    bs[16],
		fatSizeSectors:  binary.LittleEndian.Uint32(bs[36:40]),
		rootCluster:     binary.LittleEndian.Uint32(bs[44:48]),
		totalSectors:    binary.LittleEndian.Uint32(bs[32:36]),
		fsInfoSector:    binary.LittleEndian.Uint16(bs[48:50]),
	}

	bpb.dataStartSector = uint32(bpb.reservedSectors) + uint32(bpb.numberOfFATs)*bpb.fatSizeSectors
	dataSectors := bpb.totalSectors - bpb.dataStartSector
	bpb.totalClusters = dataSectors / uint32(bpb.sectorsPerClust)

	return bpb, nil
}

func fat32ReadFAT(file *os.File, bpb *fat32BPB) ([]uint32, error) {
	fatBytes := make([]byte, bpb.fatSizeSectors*uint32(bpb.sectorSize))
	fatOffset := int64(bpb.reservedSectors) * int64(bpb.sectorSize)

	_, err := file.ReadAt(fatBytes, fatOffset)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read FAT: %w", err)
	}

	entries := len(fatBytes) / 4
	fat := make([]uint32, entries)
	for i := 0; i < entries; i++ {
		fat[i] = binary.LittleEndian.Uint32(fatBytes[i*4:i*4+4]) & 0x0FFFFFFF
	}

	return fat, nil
}

func fat32WriteFAT(file *os.File, bpb *fat32BPB, fat []uint32) error {
	fatBytes := make([]byte, bpb.fatSizeSectors*uint32(bpb.sectorSize))
	for i := 0; i < len(fat) && i*4 < len(fatBytes); i++ {
		binary.LittleEndian.PutUint32(fatBytes[i*4:i*4+4], fat[i])
	}

	for i := uint8(0); i < bpb.numberOfFATs; i++ {
		fatOffset := int64(bpb.reservedSectors)*int64(bpb.sectorSize) + int64(i)*int64(bpb.fatSizeSectors)*int64(bpb.sectorSize)
		_, err := file.WriteAt(fatBytes, fatOffset)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to write FAT table %d: %w", i+1, err)
		}
	}

	return nil
}

func fat32ClusterToOffset(bpb *fat32BPB, cluster uint32) int64 {
	return int64(bpb.dataStartSector)*int64(bpb.sectorSize) + int64(cluster-2)*int64(bpb.sectorsPerClust)*int64(bpb.sectorSize)
}

func fat32GetClusterSize(bpb *fat32BPB) int {
	return int(bpb.sectorsPerClust) * int(bpb.sectorSize)
}

func fat32GetClusterChain(fat []uint32, startCluster uint32) []uint32 {
	chain := []uint32{}
	cluster := startCluster

	for cluster >= 2 && cluster < fat32EOFMarker {
		chain = append(chain, cluster)
		if int(cluster) >= len(fat) {
			break
		}
		cluster = fat[cluster]
	}

	return chain
}

func fat32ListDirectory(file *os.File, bpb *fat32BPB, fat []uint32, dirCluster uint32, prefix string, results *[]string) error {
	clusterChain := fat32GetClusterChain(fat, dirCluster)
	clusterBytes := fat32GetClusterSize(bpb)

	for _, cluster := range clusterChain {
		offset := fat32ClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to read directory cluster %d: %w", cluster, err)
		}

		for i := 0; i < clusterBytes; i += 32 {
			entry := buf[i : i+32]

			if entry[0] == 0x00 {
				return nil
			}

			if entry[0] == 0xE5 {
				continue
			}

			if entry[11] == 0x0F {
				continue
			}

			name := fat32ParseShortName(entry)
			if name == "." || name == ".." {
				continue
			}

			attr := entry[11]
			isDir := (attr & 0x10) != 0

			entryPath := prefix + "/" + name
			*results = append(*results, entryPath)

			if isDir {
				entryCluster := uint32(binary.LittleEndian.Uint16(entry[20:22]))<<16 | uint32(binary.LittleEndian.Uint16(entry[26:28]))
				err := fat32ListDirectory(file, bpb, fat, entryCluster, entryPath, results)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func fat32ParseShortName(entry []byte) string {
	name := strings.TrimRight(string(entry[0:8]), " ")
	ext := strings.TrimRight(string(entry[8:11]), " ")

	if ext != "" {
		return strings.ToLower(name + "." + ext)
	}

	return strings.ToLower(name)
}

func fat32FindEntryInDirectory(file *os.File, bpb *fat32BPB, fat []uint32, dirCluster uint32, name string, isDir bool) (uint32, bool, error) {
	clusterChain := fat32GetClusterChain(fat, dirCluster)
	clusterBytes := fat32GetClusterSize(bpb)
	searchName := strings.ToUpper(name)

	for _, cluster := range clusterChain {
		offset := fat32ClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return 0, false, tracederrors.TracedErrorf("Failed to read directory cluster %d: %w", cluster, err)
		}

		for i := 0; i < clusterBytes; i += 32 {
			entry := buf[i : i+32]

			if entry[0] == 0x00 {
				return 0, false, nil
			}

			if entry[0] == 0xE5 {
				continue
			}

			if entry[11] == 0x0F {
				continue
			}

			attr := entry[11]
			entryIsDir := (attr & 0x10) != 0

			if isDir != entryIsDir {
				continue
			}

			entryName := strings.TrimRight(string(entry[0:8]), " ")
			entryExt := strings.TrimRight(string(entry[8:11]), " ")

			fullName := entryName
			if entryExt != "" {
				fullName = entryName + "." + entryExt
			}

			if strings.ToUpper(fullName) == searchName {
				entryCluster := uint32(binary.LittleEndian.Uint16(entry[20:22]))<<16 | uint32(binary.LittleEndian.Uint16(entry[26:28]))
				return entryCluster, true, nil
			}
		}
	}

	return 0, false, nil
}

func fat32FindFileEntry(file *os.File, bpb *fat32BPB, fat []uint32, dirCluster uint32, name string) (*fat32DirEntry, error) {
	clusterChain := fat32GetClusterChain(fat, dirCluster)
	clusterBytes := fat32GetClusterSize(bpb)
	searchName := strings.ToUpper(name)

	for _, cluster := range clusterChain {
		offset := fat32ClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to read directory cluster %d: %w", cluster, err)
		}

		for i := 0; i < clusterBytes; i += 32 {
			entry := buf[i : i+32]

			if entry[0] == 0x00 {
				return nil, nil
			}

			if entry[0] == 0xE5 {
				continue
			}

			if entry[11] == 0x0F {
				continue
			}

			attr := entry[11]
			if (attr & 0x10) != 0 {
				continue
			}

			entryName := strings.TrimRight(string(entry[0:8]), " ")
			entryExt := strings.TrimRight(string(entry[8:11]), " ")

			fullName := entryName
			if entryExt != "" {
				fullName = entryName + "." + entryExt
			}

			if strings.ToUpper(fullName) == searchName {
				entryCluster := uint32(binary.LittleEndian.Uint16(entry[20:22]))<<16 | uint32(binary.LittleEndian.Uint16(entry[26:28]))
				entrySize := binary.LittleEndian.Uint32(entry[28:32])

				return &fat32DirEntry{
					name:        fullName,
					cluster:     entryCluster,
					size:        entrySize,
					isDirectory: false,
				}, nil
			}
		}
	}

	return nil, nil
}

func fat32AllocateCluster(fat []uint32) (uint32, error) {
	for i := uint32(2); i < uint32(len(fat)); i++ {
		if fat[i] == 0 {
			fat[i] = fat32EOFMarker
			return i, nil
		}
	}

	return 0, tracederrors.TracedErrorf("FAT32 image is full, no free clusters available.")
}

func fat32ZeroCluster(file *os.File, bpb *fat32BPB, cluster uint32) error {
	offset := fat32ClusterToOffset(bpb, cluster)
	zeros := make([]byte, fat32GetClusterSize(bpb))

	_, err := file.WriteAt(zeros, offset)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to zero cluster %d: %w", cluster, err)
	}

	return nil
}

func fat32WriteContent(file *os.File, bpb *fat32BPB, fat []uint32, content []byte) (uint32, error) {
	clusterBytes := fat32GetClusterSize(bpb)
	remaining := content
	var firstCluster uint32
	var prevCluster uint32

	for len(remaining) > 0 {
		cluster, err := fat32AllocateCluster(fat)
		if err != nil {
			return 0, err
		}

		if firstCluster == 0 {
			firstCluster = cluster
		}

		if prevCluster != 0 {
			fat[prevCluster] = cluster
		}

		writeSize := clusterBytes
		if len(remaining) < writeSize {
			writeSize = len(remaining)
		}

		buf := make([]byte, clusterBytes)
		copy(buf, remaining[:writeSize])

		offset := fat32ClusterToOffset(bpb, cluster)
		_, err = file.WriteAt(buf, offset)
		if err != nil {
			return 0, tracederrors.TracedErrorf("Failed to write content to cluster %d: %w", cluster, err)
		}

		remaining = remaining[writeSize:]
		prevCluster = cluster
	}

	return firstCluster, nil
}

func fat32AddDirectoryEntry(file *os.File, bpb *fat32BPB, fat []uint32, dirCluster uint32, name string, fileCluster uint32, fileSize uint32, isDir bool) error {
	shortName := fat32FormatShortName(name)
	clusterChain := fat32GetClusterChain(fat, dirCluster)
	clusterBytes := fat32GetClusterSize(bpb)

	for _, cluster := range clusterChain {
		offset := fat32ClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to read directory cluster %d: %w", cluster, err)
		}

		for i := 0; i < clusterBytes; i += 32 {
			if buf[i] == 0x00 || buf[i] == 0xE5 {
				entry := fat32BuildDirectoryEntry(shortName, fileCluster, fileSize, isDir)
				copy(buf[i:i+32], entry)

				_, err = file.WriteAt(buf, offset)
				if err != nil {
					return tracederrors.TracedErrorf("Failed to write directory entry: %w", err)
				}

				return nil
			}
		}
	}

	newCluster, err := fat32AllocateCluster(fat)
	if err != nil {
		return err
	}

	lastCluster := clusterChain[len(clusterChain)-1]
	fat[lastCluster] = newCluster

	err = fat32ZeroCluster(file, bpb, newCluster)
	if err != nil {
		return err
	}

	offset := fat32ClusterToOffset(bpb, newCluster)
	buf := make([]byte, clusterBytes)
	entry := fat32BuildDirectoryEntry(shortName, fileCluster, fileSize, isDir)
	copy(buf[0:32], entry)

	_, err = file.WriteAt(buf, offset)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write directory entry in new cluster: %w", err)
	}

	return nil
}

func fat32BuildDirectoryEntry(shortName [11]byte, cluster uint32, fileSize uint32, isDir bool) []byte {
	entry := make([]byte, 32)

	copy(entry[0:11], shortName[:])

	if isDir {
		entry[11] = 0x10
	} else {
		entry[11] = 0x20
	}

	binary.LittleEndian.PutUint16(entry[20:22], uint16(cluster>>16))
	binary.LittleEndian.PutUint16(entry[26:28], uint16(cluster&0xFFFF))

	if !isDir {
		binary.LittleEndian.PutUint32(entry[28:32], fileSize)
	}

	return entry
}

func fat32FormatShortName(name string) [11]byte {
	var result [11]byte
	for i := range result {
		result[i] = ' '
	}

	upper := strings.ToUpper(name)
	dotIndex := strings.LastIndex(upper, ".")

	if dotIndex >= 0 {
		baseName := upper[:dotIndex]
		ext := upper[dotIndex+1:]

		n := len(baseName)
		if n > 8 {
			n = 8
		}
		copy(result[0:n], []byte(baseName[:n]))

		e := len(ext)
		if e > 3 {
			e = 3
		}
		copy(result[8:8+e], []byte(ext[:e]))
	} else {
		n := len(upper)
		if n > 8 {
			n = 8
		}
		copy(result[0:n], []byte(upper[:n]))
	}

	return result
}

func fat32SplitPath(path string) []string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}

func fat32UpdateFSInfo(file *os.File, bpb *fat32BPB, fat []uint32) error {
	fsInfoOffset := int64(bpb.fsInfoSector) * int64(bpb.sectorSize)
	fs := make([]byte, bpb.sectorSize)

	_, err := file.ReadAt(fs, fsInfoOffset)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to read FSInfo sector: %w", err)
	}

	var freeCount uint32
	var nextFree uint32
	for i := uint32(2); i < uint32(len(fat)); i++ {
		if fat[i] == 0 {
			freeCount++
			if nextFree == 0 {
				nextFree = i
			}
		}
	}

	binary.LittleEndian.PutUint32(fs[488:492], freeCount)
	binary.LittleEndian.PutUint32(fs[492:496], nextFree)

	_, err = file.WriteAt(fs, fsInfoOffset)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write FSInfo sector: %w", err)
	}

	return nil
}

func fat32BuildBootSector(totalSectors uint32, fatSizeSectors uint32, volumeLabel [11]byte) []byte {
	bs := make([]byte, fat32SectorSize)

	// Jump boot code
	bs[0] = 0xEB
	bs[1] = 0x58
	bs[2] = 0x90

	// OEM Name
	copy(bs[3:11], []byte("MSWIN4.1"))

	// Bytes per sector
	binary.LittleEndian.PutUint16(bs[11:13], fat32SectorSize)

	// Sectors per cluster
	bs[13] = fat32SectorsPerCluster

	// Reserved sectors
	binary.LittleEndian.PutUint16(bs[14:16], fat32ReservedSectors)

	// Number of FATs
	bs[16] = fat32NumberOfFATs

	// Root entry count (0 for FAT32)
	binary.LittleEndian.PutUint16(bs[17:19], 0)

	// Total sectors 16 (0 for FAT32)
	binary.LittleEndian.PutUint16(bs[19:21], 0)

	// Media type
	bs[21] = fat32MediaType

	// FAT size 16 (0 for FAT32)
	binary.LittleEndian.PutUint16(bs[22:24], 0)

	// Sectors per track
	binary.LittleEndian.PutUint16(bs[24:26], 63)

	// Number of heads
	binary.LittleEndian.PutUint16(bs[26:28], 255)

	// Hidden sectors
	binary.LittleEndian.PutUint32(bs[28:32], 0)

	// Total sectors 32
	binary.LittleEndian.PutUint32(bs[32:36], totalSectors)

	// FAT size 32
	binary.LittleEndian.PutUint32(bs[36:40], fatSizeSectors)

	// Ext flags (mirroring enabled)
	binary.LittleEndian.PutUint16(bs[40:42], 0)

	// FS version
	binary.LittleEndian.PutUint16(bs[42:44], 0)

	// Root cluster
	binary.LittleEndian.PutUint32(bs[44:48], fat32RootDirCluster)

	// FSInfo sector
	binary.LittleEndian.PutUint16(bs[48:50], fat32FSInfoSectorNum)

	// Backup boot sector
	binary.LittleEndian.PutUint16(bs[50:52], fat32BackupBootSector)

	// Drive number
	bs[64] = 0x80

	// Reserved
	bs[65] = 0x00

	// Boot signature
	bs[66] = 0x29

	// Volume serial number
	binary.LittleEndian.PutUint32(bs[67:71], 0x12345678)

	// Volume label
	copy(bs[71:82], volumeLabel[:])

	// File system type
	copy(bs[82:90], []byte("FAT32   "))

	// Boot sector signature
	binary.LittleEndian.PutUint16(bs[510:512], fat32BootSignature)

	return bs
}

func fat32BuildFSInfoSector(totalClusters uint32) []byte {
	fs := make([]byte, fat32SectorSize)

	binary.LittleEndian.PutUint32(fs[0:4], fat32FSInfoSignature1)

	binary.LittleEndian.PutUint32(fs[484:488], fat32FSInfoSignature2)

	// Free cluster count (total clusters minus 1 for root dir)
	binary.LittleEndian.PutUint32(fs[488:492], totalClusters-1)

	// Next free cluster
	binary.LittleEndian.PutUint32(fs[492:496], 3)

	binary.LittleEndian.PutUint32(fs[508:512], fat32FSInfoSignature3)

	return fs
}

func fat32BuildFATTable(fatSizeSectors uint32) []byte {
	fat := make([]byte, int(fatSizeSectors)*fat32SectorSize)

	// Entry 0: Media type marker
	binary.LittleEndian.PutUint32(fat[0:4], 0x0FFFFF00|uint32(fat32MediaType))

	// Entry 1: End of chain marker
	binary.LittleEndian.PutUint32(fat[4:8], fat32ClusterInUse)

	// Entry 2: Root directory cluster (end of chain)
	binary.LittleEndian.PutUint32(fat[8:12], fat32EOFMarker)

	return fat
}

func fat32CalculateFATSize(totalSectors uint32) uint32 {
	dataSectors := totalSectors - fat32ReservedSectors
	entries := (dataSectors * fat32SectorSize) / (fat32SectorsPerCluster*fat32SectorSize + fat32NumberOfFATs*4)
	fatSizeBytes := (entries + 2) * 4
	fatSizeSectors := (fatSizeBytes + fat32SectorSize - 1) / fat32SectorSize
	return fatSizeSectors
}

func fat32PadVolumeLabel(label string) [11]byte {
	var result [11]byte
	for i := range result {
		result[i] = ' '
	}
	if label != "" {
		n := len(label)
		if n > 11 {
			n = 11
		}
		copy(result[:n], []byte(label[:n]))
	} else {
		copy(result[:], []byte("NO NAME    "))
	}
	return result
}
