package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nhassl3/article-saver-bot/pkg/e"
	"github.com/nhassl3/article-saver-bot/pkg/storage"
)

const defaultPerm = 0774

var ErrNoSavedPage = errors.New("no saved pages")

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap("can't create file", err)
	}
	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(page)
	return e.WrapIfErr("can't encode file", err)
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, e.Wrap("can't read directory", err)
	}

	if len(files) == 0 {
		return nil, e.Wrap("Zero value", ErrNoSavedPage)
	}

	// 0-9
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(len(files))

	file := files[n]

	// open & decode
	return s.DecodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) Remove(page *storage.Page) (err error) {
	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't determine file name", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	err = os.Remove(path)
	return e.WrapIfErr(fmt.Sprintf("can't remove file %s", path), err)
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't determine file name", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("can't check if file exists %s", path), err)
	}
	return true, nil
}

func (s Storage) DecodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't open file", err)
	}
	defer func() { _ = f.Close() }()

	var page storage.Page

	if err = gob.NewDecoder(f).Decode(&page); err != nil {
		return nil, e.Wrap("can't decode file", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
