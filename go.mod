module model-gen

go 1.14

require (
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/vetcher/go-astra v1.2.0
    github.com/linbaozhong/model-gen v0.0.0-20200813120133-874935090ec5

)
replace (
    github.com/linbaozhong/model-gen v0.0.0-20200813120133-874935090ec5 => ./ // indirect
)