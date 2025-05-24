package rc

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/lonelyday/rsync/config"

	"github.com/sirupsen/logrus"
)

// copyFile copies a file from src to dst and sets its modification time.
func copyFile(src, dst string, fi os.FileInfo) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, config.FilePerm)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return os.Chtimes(dst, fi.ModTime(), fi.ModTime())
}

// isValidPath checks if the given path exists and is a directory.
func isValidPath(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("path (%s) doesn't exist", path)
		}
		return fmt.Errorf("error with provided path (%s): %w", path, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("provided path %s isn't a directory", path)
	}
	return nil
}

// delMissing deletes files/dirs in dst that are not present in src.
func delMissing(src, dst map[string]bool) error {
	var errs []error
	for dp := range dst {
		if !src[dp] {
			fullPath := filepath.Join(*config.DstF, dp)
			logrus.Infof("Deleted: %s", fullPath)
			if err := os.RemoveAll(fullPath); err != nil {
				errs = append(errs, fmt.Errorf("failed to delete %s: %w", fullPath, err))
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// getPaths returns all relative paths under root.
func getPaths(root string) (map[string]bool, error) {
	paths := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logrus.Errorf("Error walking path %s: %v", path, err)
			return nil
		}
		if path == root {
			return nil
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			logrus.Errorf("Error getting relative path for %s: %v", path, err)
			return nil
		}
		paths[relPath] = true
		return nil
	})
	return paths, err
}

// Sync synchronizes files and directories from source to destination.
func Sync() error {
	srcPaths := make(map[string]bool)
	if err := isValidPath(*config.SrcF); err != nil {
		return err
	}
	if err := isValidPath(*config.DstF); err != nil {
		return err
	}
	if err := filepath.WalkDir(*config.SrcF, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(*config.SrcF, path)
		if err != nil {
			return err
		}
		if *config.DeleteMissing && (path != *config.SrcF) {
			srcPaths[relPath] = true
		}
		targetPath := filepath.Join(*config.DstF, relPath)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(targetPath, config.FolderPerm)
		}
		dstInfo, err := os.Stat(targetPath)
		if err == nil {
			if !dstInfo.IsDir() && dstInfo.ModTime().Equal(info.ModTime()) {
				return nil
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if err = copyFile(path, targetPath, info); err != nil {
			logrus.Errorf("Failed to Sync: %s (%v))", path, err)
		} else {
			logrus.Infof("Synced: %s", path)
		}
		return nil
	}); err != nil {
		return err
	}
	if *config.DeleteMissing {
		dstPaths, err := getPaths(*config.DstF)
		if err != nil {
			return err
		}
		return delMissing(srcPaths, dstPaths)
	}
	return nil
}
