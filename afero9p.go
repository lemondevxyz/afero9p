package afero9p

import (
	"fmt"
	"net"
	"time"
	"sync"
	"os"

	"aqwari.net/net/styx"
	"github.com/spf13/afero"
)

type Server struct {
	Fs      afero.Fs
	Options ServerOptions
}

func (s Server) StyxServer() *styx.Server {
	return &styx.Server{
		WriteTimeout: s.Options.WriteTimeout,
		IdleTimeout:  s.Options.IdleTimeout,
		MaxSize:      s.Options.MaxSize,
		Auth:         s.Options.Auth,
		OpenAuth:     s.Options.OpenAuth,
		ErrorLog:     s.Options.ErrorLog,
		TraceLog:     s.Options.TraceLog,
		Handler: styx.HandlerFunc(func(sesh *styx.Session) {
			for sesh.Next() {
				switch t := sesh.Request().(type) {
				case styx.Tchmod:
					t.Rchmod(s.Fs.Chmod(t.Path(), t.Mode))
				case styx.Tchown:
					t.Rchown(s.Fs.Chown(t.Path(), t.Uid, t.Gid))
				case styx.Tcreate:
					file, err := s.Fs.OpenFile(t.Name, os.O_CREATE, t.Mode)
					if err == nil {
						file.Close()
					}

					t.Rcreate(nil, err)
				case styx.Topen:
					//file, err := s.Fs.OpenFile(t.Path(), t.Flag|os.O_RDWR, 0755)
					var file afero.File
					var err error
					info, err := s.Fs.Stat(t.Path())
					if err != nil {
						t.Rerror(err.Error())
						return
					}

					if info.IsDir() {
						file, err = s.Fs.OpenFile(t.Path(), t.Flag, 0755)
					} else {
						file, err = s.Fs.OpenFile(t.Path(), os.O_RDWR|t.Flag, 0755)
					}
					t.Ropen(file, err)
				case styx.Tremove:
					t.Rremove(s.Fs.Remove(t.Path()))
				case styx.Trename:
					t.Rrename(s.Fs.Rename(t.OldPath, t.NewPath))
				case styx.Tstat:
					t.Rstat(s.Fs.Stat(t.Path()))
				case styx.Tsync:
					f, err := s.Fs.Open(t.Path())
					if err != nil {
						t.Rsync(err)
						return
					}

					t.Rsync(f.Sync())
					f.Close()
				case styx.Ttruncate:
					f, err := s.Fs.OpenFile(t.Path(), os.O_WRONLY, 0755)
					fmt.Println(t.Path(), f, err)
					if err != nil {
						t.Rtruncate(err)
						return
					}

					err = f.Truncate(t.Size)
					fmt.Println(t.Size, err)
					t.Rtruncate(err)
					f.Close()
				case styx.Tutimes:
					t.Rutimes(s.Fs.Chtimes(t.Path(), t.Atime, t.Mtime))
				case styx.Twalk:
					t.Rwalk(s.Fs.Stat(t.Path()))
				}
			}
		}),
	}
}

type ServerOptions struct {
	// inherhited from styx.Server
	Listener     net.Listener
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MaxSize      int64
	// TLSConfig    *tls.Config
	Auth     styx.AuthFunc
	OpenAuth styx.AuthOpenFunc
	ErrorLog styx.Logger
	TraceLog styx.Logger
}

func NewServer(options ServerOptions, fs afero.Fs) error {
	return Server{Fs: &openFs{
		Fs: fs,
		cache: map[string]afero.File{},
		mtx: sync.Mutex{},}, Options: options}.StyxServer().Serve(options.Listener)
}
