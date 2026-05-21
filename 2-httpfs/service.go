package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Service struct {
	RootDir string
}

func (svc *Service) GetPath(c echo.Context) (err error) {
	path := filepath.Join(svc.RootDir, c.Request().URL.Path)
	slog.Info("GetPath", "path", path)

	defer func() {
		if err != nil {
			slog.Error("GetPath", "path", path, "err", err)
		}
	}()

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return c.File(path)
	}

	var buf bytes.Buffer
	tarw := tar.NewWriter(&buf)

	root := path
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			// Skip files we can't access
			return nil
		}

		link := ""
		if (d.Type() & fs.ModeSymlink) != 0 {
			var err error
			link, err = os.Readlink(path)
			if err != nil {
				return nil
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		hdr.Name = rel

		if d.Type().IsRegular() {
			f, err := os.Open(path)
			if err != nil {
				// Skip files we can't read
				return nil
			}

			if err := tarw.WriteHeader(hdr); err != nil {
				return err
			}

			if _, err := io.Copy(tarw, f); err != nil {
				return err
			}
		} else {
			if err := tarw.WriteHeader(hdr); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := tarw.Close(); err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/x-tar", buf.Bytes())
}

func (svc *Service) PutPath(c echo.Context) error {
	path := filepath.Join(svc.RootDir, c.Request().URL.Path)
	slog.Info("PutPath", "path", path)

	f, err := os.Create(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return c.String(http.StatusNotFound, "File not found")
	case errors.Is(err, os.ErrPermission):
		return c.String(http.StatusForbidden, "Permission denied")
	case err != nil:
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, c.Request().Body)
	if err != nil {
		return err
	}

	return nil
}
