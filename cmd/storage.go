package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getStoragePool(cmd *cobra.Command) *rest.StoragePool {
	var pool *rest.StoragePool
	var err error
	switch {
	case cmd.Flags().Changed("id"):
		id, _ := cmd.Flags().GetString("id")
		pool, err = restClient.GetStoragePool(id)
	case cmd.Flags().Changed("name"):
		name, _ := cmd.Flags().GetString("name")
		pool, err = restClient.GetStoragePoolByName(name)
	default:
		cmd.Usage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pool
}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: `storage operations`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var storageBrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "list storage pool files",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("path", cmd.Flags().Lookup("path"))
		viper.BindPFlag("details", cmd.Flags().Lookup("details"))
		viper.BindPFlag("recursive", cmd.Flags().Lookup("recursive"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		files, err := pool.Browse(restClient, viper.GetString("path"), viper.GetBool("recursive"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !viper.GetBool("details") {
			paths := []string{}
			for _, file := range files {
				if !file.IsDir {
					paths = append(paths, file.Path)
				}
			}
			fmt.Println(formatString(paths))
		} else {
			fmt.Println(formatString(files))
		}
	},
}

var storageDownloadCmd = &cobra.Command{
	Use:    "download [file]",
	Short:  "download a file from a storage pool",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("output", cmd.Flags().Lookup("output"))
		viper.BindPFlag("progress-bar", cmd.Flags().Lookup("progress-bar"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		resp, err := pool.Download(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Create the file
		var out *os.File
		switch output := viper.GetString("output"); output {
		case "-":
			out = os.Stdout
		case "":
			out, err = os.Create(path.Base(args[0]))
		default:
			out, err = os.Create(output)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer out.Close()

		// Write the body to file
		if viper.GetBool("progress-bar") && viper.GetString("output") != "-" {
			bar := progressbar.DefaultBytes(
				resp.ContentLength,
				"downloading",
			)
			io.Copy(io.MultiWriter(out, bar), resp.Body)
			fmt.Println("")
		} else {
			_, err = io.Copy(out, resp.Body)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var storageConvertDiskCmd = &cobra.Command{
	Use:   "convert-disk",
	Short: "convert or copy a disk",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("src-storage")
		cmd.MarkFlagRequired("src-filename")
		cmd.MarkFlagRequired("dest-storage")
		cmd.MarkFlagRequired("dest-filename")
		viper.BindPFlag("src-storage", cmd.Flags().Lookup("src-storage"))
		viper.BindPFlag("src-filename", cmd.Flags().Lookup("src-filename"))
		viper.BindPFlag("dest-storage", cmd.Flags().Lookup("dest-storage"))
		viper.BindPFlag("dest-filename", cmd.Flags().Lookup("dest-filename"))
		viper.BindPFlag("dest-format", cmd.Flags().Lookup("dest-format"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		srcPool, err := restClient.GetStoragePoolByName(viper.GetString("src-storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		destPool, err := restClient.GetStoragePoolByName(viper.GetString("dest-storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Println("Converting Disk")
		}
		handleTask(srcPool.ConvertDisk(restClient, viper.GetString("src-filename"), destPool.ID, viper.GetString("dest-filename"), viper.GetString("dest-format")))
	},
}

var storageCopyFileCmd = &cobra.Command{
	Use:   "copy-file",
	Short: "copy a storage pool file",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("src-storage")
		cmd.MarkFlagRequired("src-filename")
		cmd.MarkFlagRequired("dest-storage")
		cmd.MarkFlagRequired("dest-filename")
		viper.BindPFlag("src-storage", cmd.Flags().Lookup("src-storage"))
		viper.BindPFlag("src-filename", cmd.Flags().Lookup("src-filename"))
		viper.BindPFlag("dest-storage", cmd.Flags().Lookup("dest-storage"))
		viper.BindPFlag("dest-filename", cmd.Flags().Lookup("dest-filename"))
		viper.BindPFlag("dest-format", cmd.Flags().Lookup("dest-format"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleTask(restClient.CopyFile(viper.GetString("src-storage"), viper.GetString("src-filename"), viper.GetString("dest-storage"), viper.GetString("dest-filename")))
	},
}

var storageCopyURLCmd = &cobra.Command{
	Use:   "copy-url",
	Short: "copy a url to the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		cmd.MarkFlagRequired("url")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("url", cmd.Flags().Lookup("url"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Println("\nDownloading " + viper.GetString("url"))
		}
		handleTask(pool.CopyURL(restClient, viper.GetString("url"), viper.GetString("filename")))
	},
}

var storageCreateDiskCmd = &cobra.Command{
	Use:   "create-disk",
	Short: "create a disk in the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("disk-format", cmd.Flags().Lookup("disk-format"))
		viper.BindPFlag("disk-size", cmd.Flags().Lookup("disk-size"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Println("\nConverting Disk " + viper.GetString("filename"))
		}
		handleTask(pool.CreateDisk(restClient, viper.GetString("filename"), viper.GetString("disk-format"), uint(viper.GetInt("disk-size"))))
	},
}

var storageCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new storage pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var sp rest.StoragePool
		err = unmarshal(data, &sp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := sp.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var storageDeleteFileCmd = &cobra.Command{
	Use:   "delete-file [file]",
	Short: "delete a file from the storage pool",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		err := pool.DeleteFile(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

var storageDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete storage pool",
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		err := pool.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var storageStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "disable a storage pool",
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		err := pool.Stop(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var storageStartCmd = &cobra.Command{
	Use:   "start",
	Short: "re-enable a storage pool",
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		err := pool.Start(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var storageDiskInfoCmd = &cobra.Command{
	Use:   "disk-info [filename]",
	Short: "get information for a disk in a storage pool",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		info, err := pool.DiskInfo(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString((info)))
	},
}

var storageGetIDCmd = &cobra.Command{
	Use:   "get-id [name]",
	Short: "get storage pool id from name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := restClient.GetStoragePoolByName(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(pool.ID)
	},
}

var storageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get storage pool details",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		fmt.Println(formatString(pool))
	},
}

var storageGrowDiskCmd = &cobra.Command{
	Use:   "grow-disk",
	Short: "grow a disk in the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("disk-size", cmd.Flags().Lookup("disk-size"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		handleTask(pool.GrowDisk(restClient, viper.GetString("filename"), uint(viper.GetInt("disk-size"))))
	},
}

var storageListCmd = &cobra.Command{
	Use:   "list",
	Short: "list storage pools",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListStoragePools(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(pools))
		} else {
			list := []map[string]string{}
			for _, pool := range pools {
				var info = map[string]string{"id": pool.ID, "name": pool.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var storageMoveFileCmd = &cobra.Command{
	Use:   "move-file",
	Short: "move a storage pool file",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("src-storage")
		cmd.MarkFlagRequired("src-filename")
		cmd.MarkFlagRequired("dest-storage")
		cmd.MarkFlagRequired("dest-filename")
		viper.BindPFlag("src-storage", cmd.Flags().Lookup("src-storage"))
		viper.BindPFlag("src-filename", cmd.Flags().Lookup("src-filename"))
		viper.BindPFlag("dest-storage", cmd.Flags().Lookup("dest-storage"))
		viper.BindPFlag("dest-filename", cmd.Flags().Lookup("dest-filename"))
		viper.BindPFlag("dest-format", cmd.Flags().Lookup("dest-format"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleTask(restClient.MoveFile(viper.GetString("src-storage"), viper.GetString("src-filename"), viper.GetString("dest-storage"), viper.GetString("dest-filename")))
	},
}

var storageUploadCmd = &cobra.Command{
	Use:   "upload [file]",
	Short: "upload a file to a storage pool",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("storageId")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("dest-filename", cmd.Flags().Lookup("dest-filename"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pool := getStoragePool(cmd)
		var filename string
		if cmd.Flags().Changed("dest-filename") {
			filename = viper.GetString("dest-filename")
		} else {
			filename = path.Base(args[0])
		}
		err := pool.Upload(restClient, args[0], filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func initIDFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	cmd.Flags().StringP("name", "n", "", "Storage Pool Name")
}

func init() {
	RootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(storageCreateCmd)
	storageCmd.AddCommand(storageDeleteCmd)
	initIDFlags(storageDeleteCmd)

	storageCmd.AddCommand(storageListCmd)
	addListFlags(storageListCmd)

	storageCmd.AddCommand(storageStopCmd)
	initIDFlags(storageStopCmd)
	storageCmd.AddCommand(storageStartCmd)
	initIDFlags(storageStartCmd)

	storageCmd.AddCommand(storageBrowseCmd)
	initIDFlags(storageBrowseCmd)
	storageBrowseCmd.Flags().String("path", "", "path inside the storage pool to browse")
	storageBrowseCmd.Flags().Bool("details", false, "detailed directory listing")
	storageBrowseCmd.Flags().Bool("recursive", false, "recursively list files")

	storageCmd.AddCommand(storageDownloadCmd)
	storageDownloadCmd.Flags().StringP("output", "o", "", "output file")
	storageDownloadCmd.Flags().Bool("progress-bar", false, "show a progress bar with --wait")
	initIDFlags(storageDownloadCmd)

	storageCmd.AddCommand(storageConvertDiskCmd)
	storageConvertDiskCmd.Flags().String("src-storage", "", "Source storage pool name")
	storageConvertDiskCmd.Flags().String("src-filename", "", "Source filename")
	storageConvertDiskCmd.Flags().String("dest-storage", "", "Destination storage pool name")
	storageConvertDiskCmd.Flags().String("dest-filename", "", "Destination filename")
	storageConvertDiskCmd.Flags().String("dest-format", "qcow2", "Destination file format")
	addTaskFlags(storageConvertDiskCmd)

	storageCmd.AddCommand(storageCopyFileCmd)
	storageCopyFileCmd.Flags().String("src-storage", "", "Source storage pool id")
	storageCopyFileCmd.Flags().String("src-filename", "", "path to the file in the source storage pool")
	storageCopyFileCmd.Flags().String("dest-storage", "", "Destination storage pool id")
	storageCopyFileCmd.Flags().String("dest-filename", "", "path to the file in the destination storage pool")
	addTaskFlags(storageCopyFileCmd)

	storageCmd.AddCommand(storageCopyURLCmd)
	initIDFlags(storageCopyURLCmd)
	storageCopyURLCmd.Flags().String("filename", "", "filename for the disk")
	storageCopyURLCmd.Flags().String("url", "", "url to download")
	addTaskFlags(storageCopyURLCmd)

	storageCmd.AddCommand(storageCreateDiskCmd)
	initIDFlags(storageCreateDiskCmd)
	storageCreateDiskCmd.Flags().String("filename", "", "filename for the disk")
	storageCreateDiskCmd.Flags().String("disk-format", "qcow2", "disk format ()")
	storageCreateDiskCmd.Flags().Int("disk-size", 25, "size of the disk in GB")
	addTaskFlags(storageCreateDiskCmd)

	storageCmd.AddCommand(storageDeleteFileCmd)
	initIDFlags(storageDeleteFileCmd)

	storageCmd.AddCommand(storageDiskInfoCmd)
	initIDFlags(storageDiskInfoCmd)

	storageCmd.AddCommand(storageGetIDCmd)

	storageCmd.AddCommand(storageGetCmd)
	initIDFlags(storageGetCmd)

	storageCmd.AddCommand(storageGrowDiskCmd)
	initIDFlags(storageGrowDiskCmd)
	storageGrowDiskCmd.Flags().String("filename", "", "filename for the disk")
	storageGrowDiskCmd.Flags().Int("disk-size", 0, "size to add in GB")
	addTaskFlags(storageGrowDiskCmd)

	storageCmd.AddCommand(storageMoveFileCmd)
	storageMoveFileCmd.Flags().String("src-storage", "", "Source storage pool id")
	storageMoveFileCmd.Flags().String("src-filename", "", "path to the file in the source storage pool")
	storageMoveFileCmd.Flags().String("dest-storage", "", "Destination storage pool id")
	storageMoveFileCmd.Flags().String("dest-filename", "", "path to the file in the destination storage pool")
	addTaskFlags(storageMoveFileCmd)

	storageCmd.AddCommand(storageUploadCmd)
	initIDFlags(storageUploadCmd)
	storageUploadCmd.Flags().String("dest-filename", "", "path to the file in the destination storage pool")

}
