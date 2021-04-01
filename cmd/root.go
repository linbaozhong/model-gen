package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	module  string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "model-gen",
		Short: "model 访问字典生成器",
		Long: `model 访问字典生成器.
在需要生成的struct上增加注释 //tablename [表名]
				`,
		Run: func(cmd *cobra.Command, args []string) {
			//module
			p := path[:1]
			if p == "." || p == "\\" {
				p = path[strings.Index(path, "\\")+1:]
				if len(p) > 0 {
					module += "/" + p
				}
			} else {
				module += "/" + path
			}

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
					return handleFile(module, filename)
				}
				return nil
			})
			if err != nil {
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

	rootCmd.PersistentFlags().StringVarP(&tagName, "tag", "t", XORM_TAG, "ORM名称.支持:xorm,gorm")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "./models", "models路径")
	rootCmd.PersistentFlags().StringVarP(&module, "module", "m", "", "module名称")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "true:调试模式")

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
