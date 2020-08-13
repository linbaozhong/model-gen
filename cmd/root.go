package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	XORM_TAG = "xorm"
	GORM_TAG = "gorm"
)

var (
	// Used for flags.
	tagName string
	path    string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "model-gen",
		Short: "model 访问字典生成器",
		Long:  `model 访问字典生成器.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Mkdir(path+"/table", os.ModePerm)

			if err := writeBaseFile(path + "/table/base_sorm.go"); err != nil {
				showError(err.Error())
				return
			}

			err := filepath.Walk(path, func(filename string, f os.FileInfo, _ error) error {
				if f.IsDir() && filename != path {
					return filepath.SkipDir
				}
				if filepath.Ext(filename) == ".go" {
					if strings.Contains(filename, "_table.go") || strings.Contains(filename, "_sorm.go") {
						return nil
					}
					_, e := handleFile(filename)
					return e
				}
				return nil
			})
			if err != nil {
				showError(err.Error())
			}
			if err := writeBaseFile(path + "/table"); err != nil {
				showError(err.Error())
			}
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&tagName, "tag", "t", GORM_TAG, "orm tag name")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "./models", "models path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug true")

}

func showError(msg interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println("Error:", msg, file, line)
	os.Exit(1)
}

func initConfig() {
	if path == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			showError(err)
		}
		path = home
	}
}
