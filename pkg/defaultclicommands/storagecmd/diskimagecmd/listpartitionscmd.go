package diskimagecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/storage/diskimageutils"
)

func NewListPartitionsCmd() *cobra.Command {
	const short = "List partitions of the diskimage --image-path."
	const long = short + `

Example usage:

  $ a-helper storage disk-image list-partitions --image-path=rpi.img
  Partition table type: MBR
  Disk image: rpi.img
  Partitions found: 2

    Partition 1:
      Start LBA:  16384
      Size LBA:   1048576
      Size Bytes: 536870912
      Bootable:   false
      FileSystem: FAT32 (LBA)
      Type Code:  0x0C

    Partition 2:
      Start LBA:  1064960
      Size LBA:   4620288
      Size Bytes: 2365587456
      Bootable:   false
      FileSystem: Linux
      Type Code:  0x83`

	cmd := &cobra.Command{
		Use:   "list-partitions",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			imagePath, err := cmd.Flags().GetString("image-path")
			if err != nil {
				return err
			}

			if imagePath == "" {
				return fmt.Errorf("flag '--image-path' is required but empty")
			}

			partitionTable, err := diskimageutils.ReadPartitionTable(cmd.Context(), imagePath)
			if err != nil {
				return err
			}

			fmt.Printf("Partition table type: %s\n", partitionTable.TableType)
			fmt.Printf("Disk image: %s\n", partitionTable.DiskImagePath)
			fmt.Printf("Partitions found: %d\n\n", len(partitionTable.Partitions))

			for _, p := range partitionTable.Partitions {
				fmt.Printf("  Partition %d:\n", p.Index)
				fmt.Printf("    Start LBA:  %d\n", p.StartLBA)
				fmt.Printf("    Size LBA:   %d\n", p.SizeLBA)
				fmt.Printf("    Size Bytes: %d\n", p.SizeBytes)
				fmt.Printf("    Bootable:   %t\n", p.Bootable)
				if p.FileSystem != "" {
					fmt.Printf("    FileSystem: %s\n", p.FileSystem)
				}
				if p.TypeCode != 0 {
					fmt.Printf("    Type Code:  0x%02X\n", p.TypeCode)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().String("image-path", "", "Path of the image")

	return cmd
}
