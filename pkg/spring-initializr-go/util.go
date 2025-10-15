package springinitializr

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ensureDir(root *os.Root, dir string, perm os.FileMode) error {
	if dir == "" || dir == "." {
		return nil
	}
	dir = filepath.Clean(dir)
	parts := strings.Split(dir, string(os.PathSeparator))
	cur := ""
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		cur = filepath.Join(cur, p)
		err := root.Mkdir(cur, perm)
		if err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
	}
	return nil
}

func unzipProject(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	root, err := os.OpenRoot(dest)
	if err != nil {
		return fmt.Errorf("could not open directory %q: %w", dest, err)
	}
	defer root.Close()

	for _, f := range r.File {
		name := filepath.Clean(f.Name)

		if name == "." || name == "" {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := ensureDir(root, name, f.Mode().Perm()); err != nil {
				return err
			}
			continue
		}

		parent := filepath.Dir(name)
		if parent != "." && parent != "/" {
			if err := ensureDir(root, parent, 0o755); err != nil {
				return err
			}
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := root.Create(name)
		if err != nil {
			rc.Close()
			return fmt.Errorf("error creating %q in project dir: %w", name, err)
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

	}

	return nil
}
