package qs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Doc struct {
	Cate string `yaml:"cate"`
	Xxx  []Xxx  `yaml:"xxx"`
}

type Xxx struct {
	Name string     `yaml:"name"`
	Qs   [][]string `yaml:"qs"`
}

type Docs []Doc

func NewDocs(fp string) Docs {
	var docs Docs
	if PathExists(fp) {
		err := filepath.WalkDir(fp, func(path string, de fs.DirEntry, err error) error {
			if !de.IsDir() {
				f, err := Load(path)
				if err != nil {
					return err
				}
				d := yaml.NewDecoder(bytes.NewReader(f))
				for {
					spec := new(Doc)
					spec.Cate = de.Name()
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
			return nil
		})
		if err != nil {
			return nil
		}
	}
	return docs
}

// GetNames Get All Names
// func (d Docs) GetNames() (names []string) {
// 	for _, doc := range d {
// 		names = append(names, doc.Name)
// 	}
// 	return
// }
//
// func (d Docs) GetNameByCate(cate string) (names []string) {
// 	for _, doc := range d {
// 		if doc.Cate == cate {
// 			names = append(names, doc.Name)
// 		}
// 	}
// 	return
// }
//
//
// func (d Docs) IsHitName(query string) bool {
// 	return lo.ContainsBy(d, func(item Doc) bool {
// 		return strings.EqualFold(item.Name, query)
// 	})
// }
//
// func (d Docs) GetQsByName(name string) []Xxx {
// 	for _, doc := range d {
// 		if strings.ToLower(doc.Name) == strings.ToLower(name) {
// 			return doc.Xxx
// 		}
// 	}
// 	return nil
// }
//
// func (d Docs) SearchQs(query string) (qs []Xxx) {
// 	for _, doc := range d {
// 		for _, xxx := range doc.Xxx {
// 			qsLower := strings.ToLower(xxx.Qs)
// 			asLower := strings.ToLower(xxx.As)
// 			query = strings.ToLower(query)
// 			if strings.Contains(qsLower, query) || strings.Contains(asLower, query) {
// 				qs = append(qs, xxx)
// 			}
// 		}
// 	}
// 	return
// }

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
