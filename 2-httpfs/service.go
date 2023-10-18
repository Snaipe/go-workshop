package main

import (
	"archive/tar"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Service struct {
	RootDir string
}

func processFileError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return c.NoContent(http.StatusNotFound)
	case errors.Is(err, os.ErrPermission):
		return c.NoContent(http.StatusForbidden)
	}
	return err
}

func (svc *Service) GetPath(c echo.Context) error {
	path := filepath.Join(svc.RootDir, c.Request().URL.Path)
	slog.Info("GetPath", "path", path)

	f, err := os.Open(path)
	if err != nil {
		slog.Error("open", "path", path, "error", err)
		return processFileError(c, err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		slog.Error("stat", "path", path, "error", err)
		return processFileError(c, err)
	}

	if !finfo.IsDir() {
		c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", finfo.Size()))
		return c.Stream(http.StatusOK, mime.TypeByExtension(filepath.Ext(path)), f)
	}

	resp := bufio.NewWriter(c.Response())
	defer resp.Flush()

	tw := tar.NewWriter(resp)
	defer tw.Close()

	orig := path
	err = filepath.WalkDir(orig, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("dir stat", "err", err)
			return err
		}
		slog.Info("visiting", "path", path)
		finfo, err := os.Lstat(path)
		switch {
		case errors.Is(err, os.ErrPermission):
			return nil
		case err != nil:
			slog.Error("stat", "err", err)
			return err
		}

		var linkname string
		if (finfo.Mode() & fs.ModeType) == fs.ModeSymlink {
			linkname, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(finfo, linkname)
		if err != nil {
			slog.Error("tar.FileInfoHeader", "err", err)
			return err
		}

		hdr.Name = filepath.Join(".", path[len(orig):])

		switch finfo.Mode() & fs.ModeType {
		case fs.ModeSymlink:
			slog.Info("symlink", "target", hdr.Linkname)
			fallthrough
		case fs.ModeDir:
			if err := tw.WriteHeader(hdr); err != nil {
				slog.Error("tar.WriteHeader", "err", err)
				return err
			}

			return nil
		}

		f, err := os.Open(path)
		switch {
		case errors.Is(err, os.ErrPermission):
			return nil
		case err != nil:
			slog.Error("open", "err", err)
			return err
		}
		defer f.Close()

		if err := tw.WriteHeader(hdr); err != nil {
			slog.Error("tar.WriteHeader", "err", err)
			return err
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			slog.Error("copy", "err", err)
		}
		return err
	})
	if err != nil {
		slog.Error("walk", "err", err)
		return processFileError(c, err)
	}

	return nil
}

func (svc *Service) PutPath(c echo.Context) error {
	path := filepath.Join(svc.RootDir, c.Request().URL.Path)
	slog.Info("PutPath", "path", path)

	f, err := os.Create(path)
	switch {
	case errors.Is(err, os.ErrPermission):
		return c.NoContent(http.StatusForbidden)
	case err != nil:
		slog.Error("create", "path", path, "error", err)
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, c.Request().Body); err != nil {
		return err
	}

	return f.Close()
}
