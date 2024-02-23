package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

const tpl = `{{ $prevLevel2 := "" }}

{{range .}}
{{if ne $prevLevel2 .Cate}}
## {{.Cate}}
{{ $prevLevel2 = .Cate }}
{{end}}

### {{.Name}}

{{range .Xxx}}
- {{.Qs}}
{{- end}}

{{end}}
`

// mdCmd represents the md command
var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		docs := NewDocs(cfgFile)

		tmpl := template.Must(template.New("").Parse(tpl))

		file, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		err = tmpl.Execute(file, docs)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Markdown output has been written to output.md")
	},
}

var (
	cfgFile    string
	targetFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(mdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "src/data/qs.yml", "config file (default is src/data/qs.yml)")
	rootCmd.PersistentFlags().StringVar(&targetFile, "target", "qs.md", "target file (default is qs.md)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("src/data/")
		viper.SetConfigName("qs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

type Doc struct {
	Name string `yaml:"name,omitempty"`
	Cate string `yaml:"cate,omitempty"`
	Xxx  []Xxx  `yaml:"xxx,omitempty"`
}

type Xxx struct {
	Qs string `yaml:"qs,omitempty"`
	As string `yaml:"as,omitempty"`
}

type Docs []Doc

var once sync.Once

func NewDocs(fp string) Docs {
	var docs Docs
	once.Do(func() {
		if PathExists(fp) {
			f, err := Load(fp)
			if err != nil {
				return
			}
			d := yaml.NewDecoder(bytes.NewReader(f))
			for {
				spec := new(Doc)
				if err := d.Decode(&spec); err != nil {
					// break the loop in case of EOF
					if errors.Is(err, io.EOF) {
						break
					}
					panic(err)
				}
				if spec != nil {
					docs = append(docs, *spec)
				}
			}
		}
	})
	return docs
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

// Load reads data saved under given name.
func Load(name string) ([]byte, error) {
	// p := c.path(name)
	if _, err := os.Stat(name); err != nil {
		return nil, err
	}
	return os.ReadFile(name)
}
