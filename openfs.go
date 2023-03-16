package afero9p

import (
	"os"
	"sync"

	"github.com/spf13/afero"
)

type openFs struct {
	afero.Fs
	cache map[string]afero.File
	mtx   sync.Mutex
}

type openFile struct {
	afero.File
	name   string
	fs    *openFs
	closed bool
	mtx    sync.Mutex
}

func (f *openFs) cacheExists(name string) afero.File {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	file, _ := f.cache[name]
	return file
}

func (f *openFs) returnFile(name string, file afero.File, err error) (afero.File, error) {
	if err != nil {
		return nil, err
	}

	return &openFile{file, name, f, false, sync.Mutex{}}, nil
}

func (f *openFs) deleteCache(name string) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	delete(f.cache, name)
}

func (f *openFile) Close() error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if !f.closed {
		f.closed = true
		f.fs.deleteCache(f.name)
		f.File.Close()

		return nil
	}

	return nil
}

func (f *openFs) Open(name string) (afero.File, error) {
	if fi := f.cacheExists(name); fi != nil {
		return &openFile{fi, name, f, false, sync.Mutex{}}, nil
	}

	fi, err := f.Fs.Open(name)
	return f.returnFile(name, fi, err)
}

func (f *openFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if fi := f.cacheExists(name); fi != nil {
		return fi, nil
	}

	fi, err := f.Fs.OpenFile(name, flag, perm)
	return fi, err
}

func (f *openFs) Create(name string) (afero.File, error) {
	if fi := f.cacheExists(name); fi != nil {
		return fi, nil
	}

	fi, err := f.Fs.Create(name)
	return f.returnFile(name, fi, err)
}
