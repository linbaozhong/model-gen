package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
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

			err := filepath.Walk(path, func(filename string, f os.FileInfo, _ error) error {
				if f.IsDir() && filename != path {
					return filepath.SkipDir
				}
				if filepath.Ext(filename) == ".go" {
					if strings.Contains(filename, "_table.go") {
						return nil
					}
					return handleFile(filename)
				}
				return nil
			})
			if err != nil {
				fmt.Println(err.Error())
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

	rootCmd.PersistentFlags().StringVarP(&tagName, "tag", "t", GORM_TAG, "orm tag name (default is xorm)")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "./models", "models path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug true")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if path == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}
		path = home
	}
}
