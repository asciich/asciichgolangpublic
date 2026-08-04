package diskimagecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/storage/diskimageutils"
)

func NewListPartitionsCmd() *cobra.Command {
	const short = "List partitions of the diskimage --image-path."

	cmd := &cobra.Command{
		Use:   "list-partitions",
		Short: short,
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