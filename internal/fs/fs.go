package fs

import (
	"io/ioutil"
	"os"
	"path"
)

func NewFileSystem(root string) *FileSystem {
	return &FileSystem{
		Root: root,
	}
}

type FileSystem struct {
	Root string
}

type FileInfo struct {
	Name  string
	IsDir bool
}

func (s *FileSystem) List(p string) ([]*FileInfo, error) {
	files, err := ioutil.ReadDir(path.Join(s.Root, p))
	if err != nil {
		return nil, err
	}

	var infos []*FileInfo
	for _, f := range files {
		infos = append(infos, &FileInfo{
			Name:  f.Name(),
			IsDir: f.IsDir(),
		})
	}

	return infos, nil
}

func (s *FileSystem) Remove(p string, name string) error {

	return os.RemoveAll(path.Join(s.Root, p, name))
}

func (s *FileSystem) Mkdir(p string) error {
	return os.MkdirAll(path.Join(s.Root, p), 0755)
}

func (s *FileSystem) CreateFile(p string, name string, data []byte) error {
	if err := s.Mkdir(path.Join(s.Root, p)); err != nil {
		return err
	}

	fp, err := os.OpenFile(path.Join(s.Root, p, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if _, err := fp.Write(data); err != nil {
		return err
	}

	return nil
}
