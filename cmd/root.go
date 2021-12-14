package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

const (
	XORM_TAG  = "xorm"
	GORM_TAG  = "gorm"
	Separator = string(os.PathSeparator)
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
			//
			dir, e := ioutil.ReadDir(path)
			if e != nil {
				showError(e.Error())
				return
			}
			//module_path
			module_path := module
			p := path[:1]
			if p == "." {
				pos := strings.Index(path, Separator)
				if pos == -1 {
					p = path[1:]
				} else {
					p = path[pos+1:]
				}

				if len(p) > 0 {
					module_path += "/" + strings.Replace(p, Separator, "/", -1)
				}
			} else {
				module_path += "/" + strings.Replace(path, Separator, "/", -1)
			}
			//module
			pos := strings.Index(module, "/")
			if pos > 0 {
				module = module[:pos]
			}

			if e = os.Mkdir(filepath.Join(path, "dao"), os.ModePerm); e != nil && !os.IsExist(e) {
				showError(e.Error())
				//return
			}
			if e = writeDaoBaseFile(filepath.Join(path, "dao", "a_base.go"), module_path); e != nil {
				showError(e.Error())
			}
			if e = os.Mkdir(filepath.Join(path, "table"), os.ModePerm); e != nil && !os.IsExist(e) {
				showError(e.Error())
				//return
			}
			if e = writeBaseFile(filepath.Join(path, "table", "base_sorm.go")); e != nil {
				showError(e.Error())
				//return
			}
			if e = writeBuildFile(filepath.Join(path, "table", "build_sorm.go")); e != nil {
				showError(e.Error())
			}

			for _, f := range dir {
				if f.IsDir() {
					continue
				}
				var filename = f.Name()
				if filepath.Ext(filename) == ".go" {
					if strings.Contains(filename, "_table.go") || strings.Contains(filename, "_sorm.go") {
						continue
					}
					handleFile(module, module_path, filepath.Join(path, f.Name()))
				}
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

//
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
