package diskimageutils

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"path"
	"strings"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	fat32MinImageSize = 33 * 1024 * 1024
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

	d, err := diskfs.Create(outputPath, options.SizeBytes, diskfs.Raw, diskfs.SectorSizeDefault)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create disk image '%s': %w", outputPath, err)
	}

	spec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeFat32,
		VolumeLabel: options.VolumeLabel,
	}

	_, err = d.CreateFilesystem(spec)
	if err != nil {
		os.Remove(outputPath)
		return tracederrors.TracedErrorf("Failed to create FAT32 filesystem on '%s': %w", outputPath, err)
	}

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

	fs, err := fat32OpenFilesystem(imagePath)
	if err != nil {
		return nil, err
	}

	results := []string{"."}

	err = fat32ListDirectoryRecursive(fs, "/", ".", &results)
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

	fs, err := fat32OpenFilesystemReadWrite(imagePath)
	if err != nil {
		return err
	}

	normalizedPath := fat32NormalizePath(filePath)

	// Create parent directories
	dir := path.Dir(normalizedPath)
	if dir != "/" && dir != "." {
		err = fat32MkdirAll(fs, dir)
		if err != nil {
			return err
		}
	}

	file, err := fs.OpenFile(normalizedPath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open file '%s' for writing in FAT32 image '%s': %w", filePath, imagePath, err)
	}

	_, err = file.Write(content)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write content to file '%s' in FAT32 image '%s': %w", filePath, imagePath, err)
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

	fs, err := fat32OpenFilesystem(imagePath)
	if err != nil {
		return nil, err
	}

	normalizedPath := fat32NormalizePath(filePath)

	file, err := fs.OpenFile(normalizedPath, os.O_RDONLY)
	if err != nil {
		return nil, tracederrors.TracedErrorf("File '%s' not found in FAT32 image '%s'.", filePath, imagePath)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read file '%s' from FAT32 image '%s': %w", filePath, imagePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Read file '%s' from FAT32 image '%s' finished. Size: '%d' bytes.", filePath, imagePath, len(content))

	return content, nil
}

// Fat32Exists checks if a file or directory exists at 'path' in a FAT32 image.
func Fat32Exists(ctx context.Context, imagePath string, filePath string) (bool, error) {
	if imagePath == "" {
		return false, tracederrors.TracedErrorEmptyString("imagePath")
	}

	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	files, err := Fat32ListFiles(ctx, imagePath)
	if err != nil {
		return false, err
	}

	searchPath := "./" + strings.TrimPrefix(strings.TrimPrefix(filePath, "/"), "./")
	for _, f := range files {
		if strings.EqualFold(f, searchPath) {
			return true, nil
		}
	}

	return false, nil
}

// Fat32DeleteFile deletes a file at 'filePath' from a FAT32 image.
//
// The 'filePath' must use forward slashes as separator (e.g. "test/a.txt").
// This operates at the raw level by marking the directory entry as deleted (0xE5)
// and freeing the associated FAT cluster chain.
func Fat32DeleteFile(ctx context.Context, imagePath string, filePath string) error {
	if imagePath == "" {
		return tracederrors.TracedErrorEmptyString("imagePath")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	logging.LogInfoByCtxf(ctx, "Delete file '%s' from FAT32 image '%s' started.", filePath, imagePath)

	file, err := os.OpenFile(imagePath, os.O_RDWR, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open FAT32 image '%s': %w", imagePath, err)
	}
	defer file.Close()

	bpb, err := fat32DeleteReadBPB(file)
	if err != nil {
		return err
	}

	fat, err := fat32DeleteReadFAT(file, bpb)
	if err != nil {
		return err
	}

	parts := fat32DeleteSplitPath(filePath)
	if len(parts) == 0 {
		return tracederrors.TracedErrorf("Invalid file path '%s'.", filePath)
	}

	// Navigate to parent directory
	dirCluster := bpb.rootCluster
	for i := 0; i < len(parts)-1; i++ {
		cluster, found, err := fat32DeleteFindDirectory(file, bpb, fat, dirCluster, parts[i])
		if err != nil {
			return err
		}
		if !found {
			return tracederrors.TracedErrorf("Directory '%s' not found in FAT32 image '%s'.", parts[i], imagePath)
		}
		dirCluster = cluster
	}

	// Find and delete the file entry
	fileName := parts[len(parts)-1]
	fileCluster, err := fat32DeleteMarkEntry(file, bpb, fat, dirCluster, fileName)
	if err != nil {
		return err
	}

	// Free the cluster chain
	if fileCluster >= 2 && fileCluster < 0x0FFFFFF8 {
		fat32DeleteFreeClusterChain(fat, fileCluster)
	}

	// Write updated FAT back
	err = fat32DeleteWriteFAT(file, bpb, fat)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete file '%s' from FAT32 image '%s' finished.", filePath, imagePath)

	return nil
}

// --- Internal types and functions for delete ---

type fat32DeleteBPB struct {
	sectorSize      uint16
	sectorsPerClust uint8
	reservedSectors uint16
	numberOfFATs    uint8
	fatSizeSectors  uint32
	rootCluster     uint32
	dataStartSector uint32
}

func fat32DeleteReadBPB(file *os.File) (*fat32DeleteBPB, error) {
	bs := make([]byte, 512)
	_, err := file.ReadAt(bs, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read boot sector: %w", err)
	}

	sig := binary.LittleEndian.Uint16(bs[510:512])
	if sig != 0xAA55 {
		return nil, tracederrors.TracedErrorf("Invalid boot sector signature: 0x%04X.", sig)
	}

	bpb := &fat32DeleteBPB{
		sectorSize:      binary.LittleEndian.Uint16(bs[11:13]),
		sectorsPerClust: bs[13],
		reservedSectors: binary.LittleEndian.Uint16(bs[14:16]),
		numberOfFATs:    bs[16],
		fatSizeSectors:  binary.LittleEndian.Uint32(bs[36:40]),
		rootCluster:     binary.LittleEndian.Uint32(bs[44:48]),
	}

	bpb.dataStartSector = uint32(bpb.reservedSectors) + uint32(bpb.numberOfFATs)*bpb.fatSizeSectors

	return bpb, nil
}

func fat32DeleteReadFAT(file *os.File, bpb *fat32DeleteBPB) ([]uint32, error) {
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

func fat32DeleteWriteFAT(file *os.File, bpb *fat32DeleteBPB, fat []uint32) error {
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

func fat32DeleteClusterToOffset(bpb *fat32DeleteBPB, cluster uint32) int64 {
	return int64(bpb.dataStartSector)*int64(bpb.sectorSize) + int64(cluster-2)*int64(bpb.sectorsPerClust)*int64(bpb.sectorSize)
}

func fat32DeleteGetClusterSize(bpb *fat32DeleteBPB) int {
	return int(bpb.sectorsPerClust) * int(bpb.sectorSize)
}

func fat32DeleteGetClusterChain(fat []uint32, startCluster uint32) []uint32 {
	chain := []uint32{}
	cluster := startCluster

	for cluster >= 2 && cluster < 0x0FFFFFF8 {
		chain = append(chain, cluster)
		if int(cluster) >= len(fat) {
			break
		}
		cluster = fat[cluster]
	}

	return chain
}

func fat32DeleteFindDirectory(file *os.File, bpb *fat32DeleteBPB, fat []uint32, dirCluster uint32, name string) (uint32, bool, error) {
	clusterChain := fat32DeleteGetClusterChain(fat, dirCluster)
	clusterBytes := fat32DeleteGetClusterSize(bpb)
	searchName := strings.ToUpper(name)

	for _, cluster := range clusterChain {
		offset := fat32DeleteClusterToOffset(bpb, cluster)
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

			// Skip LFN entries
			if entry[11] == 0x0F {
				continue
			}

			attr := entry[11]
			if (attr & 0x10) == 0 {
				continue // Not a directory
			}

			entryName := fat32DeleteGetEntryName(entry)
			if strings.ToUpper(entryName) == searchName {
				entryCluster := uint32(binary.LittleEndian.Uint16(entry[20:22]))<<16 | uint32(binary.LittleEndian.Uint16(entry[26:28]))
				return entryCluster, true, nil
			}
		}
	}

	return 0, false, nil
}

// fat32DeleteMarkEntry finds the file entry in the directory and marks it as deleted (0xE5).
// Also marks any associated LFN entries as deleted.
// Returns the start cluster of the deleted file.
func fat32DeleteMarkEntry(file *os.File, bpb *fat32DeleteBPB, fat []uint32, dirCluster uint32, name string) (uint32, error) {
	clusterChain := fat32DeleteGetClusterChain(fat, dirCluster)
	clusterBytes := fat32DeleteGetClusterSize(bpb)
	searchName := strings.ToUpper(name)

	for _, cluster := range clusterChain {
		offset := fat32DeleteClusterToOffset(bpb, cluster)
		buf := make([]byte, clusterBytes)

		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return 0, tracederrors.TracedErrorf("Failed to read directory cluster %d: %w", cluster, err)
		}

		for i := 0; i < clusterBytes; i += 32 {
			entry := buf[i : i+32]

			if entry[0] == 0x00 {
				return 0, tracederrors.TracedErrorf("File '%s' not found for deletion.", name)
			}

			if entry[0] == 0xE5 {
				continue
			}

			// Skip LFN entries for now, we'll handle them when we find the matching short entry
			if entry[11] == 0x0F {
				continue
			}

			attr := entry[11]
			if (attr & 0x10) != 0 {
				continue // Skip directories
			}

			entryName := fat32DeleteGetEntryName(entry)
			if strings.ToUpper(entryName) == searchName {
				fileCluster := uint32(binary.LittleEndian.Uint16(entry[20:22]))<<16 | uint32(binary.LittleEndian.Uint16(entry[26:28]))

				// Mark preceding LFN entries as deleted
				fat32DeleteMarkLFNEntries(buf, i)

				// Mark the short name entry as deleted
				buf[i] = 0xE5

				// Write modified buffer back
				_, err = file.WriteAt(buf, offset)
				if err != nil {
					return 0, tracederrors.TracedErrorf("Failed to write deleted entry: %w", err)
				}

				return fileCluster, nil
			}
		}
	}

	return 0, tracederrors.TracedErrorf("File '%s' not found for deletion.", name)
}

// fat32DeleteMarkLFNEntries marks all LFN entries preceding the short entry at position shortEntryOffset as deleted.
func fat32DeleteMarkLFNEntries(buf []byte, shortEntryOffset int) {
	// Walk backwards from the short entry and mark all consecutive LFN entries
	for j := shortEntryOffset - 32; j >= 0; j -= 32 {
		if buf[j] == 0xE5 || buf[j] == 0x00 {
			break
		}
		if buf[j+11] != 0x0F {
			break // Not an LFN entry
		}
		buf[j] = 0xE5
	}
}

func fat32DeleteFreeClusterChain(fat []uint32, startCluster uint32) {
	cluster := startCluster

	for cluster >= 2 && cluster < 0x0FFFFFF8 {
		if int(cluster) >= len(fat) {
			break
		}
		next := fat[cluster]
		fat[cluster] = 0
		cluster = next
	}
}

func fat32DeleteGetEntryName(entry []byte) string {
	name := strings.TrimRight(string(entry[0:8]), " ")
	ext := strings.TrimRight(string(entry[8:11]), " ")

	if ext != "" {
		return name + "." + ext
	}

	return name
}

func fat32DeleteSplitPath(filePath string) []string {
	filePath = strings.TrimPrefix(filePath, "/")
	filePath = strings.TrimPrefix(filePath, "./")
	filePath = strings.TrimSuffix(filePath, "/")

	if filePath == "" {
		return nil
	}

	return strings.Split(filePath, "/")
}

func fat32OpenFilesystemReadWrite(imagePath string) (filesystem.FileSystem, error) {
	d, err := diskfs.Open(imagePath, diskfs.WithOpenMode(diskfs.ReadWriteExclusive))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open disk image '%s' for read-write: %w", imagePath, err)
	}

	fs, err := d.GetFilesystem(0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get FAT32 filesystem from '%s': %w", imagePath, err)
	}

	return fs, nil
}

// --- Internal functions ---

func fat32OpenFilesystem(imagePath string) (filesystem.FileSystem, error) {
	d, err := diskfs.Open(imagePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to open disk image '%s': %w", imagePath, err)
	}

	fs, err := d.GetFilesystem(0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get FAT32 filesystem from '%s': %w", imagePath, err)
	}

	return fs, nil
}

func fat32ListDirectoryRecursive(fs filesystem.FileSystem, dirPath string, prefix string, results *[]string) error {
	entries, err := fs.ReadDir(dirPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to read directory '%s': %w", dirPath, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if name == "." || name == ".." {
			continue
		}

		// Skip volume label entries (they appear as root-level non-dir, size 0, no extension)
		if dirPath == "/" && !entry.IsDir() && entry.Size() == 0 && !strings.Contains(name, ".") {
			continue
		}

		entryPath := prefix + "/" + name
		*results = append(*results, entryPath)

		if entry.IsDir() {
			subDirPath := path.Join(dirPath, name)
			err := fat32ListDirectoryRecursive(fs, subDirPath, entryPath, results)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fat32MkdirAll(fs filesystem.FileSystem, dirPath string) error {
	parts := strings.Split(strings.Trim(dirPath, "/"), "/")
	current := ""

	for _, part := range parts {
		current = current + "/" + part
		err := fs.Mkdir(current)
		if err != nil {
			// Ignore error if directory already exists
			if !strings.Contains(err.Error(), "already exists") {
				return tracederrors.TracedErrorf("Failed to create directory '%s': %w", current, err)
			}
		}
	}

	return nil
}

func fat32NormalizePath(filePath string) string {
	filePath = strings.TrimPrefix(filePath, "./")
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	return filePath
}
