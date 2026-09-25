package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListFiles returns every file under the exercise's answers workspace.
func (a *App) ListFiles(exID string) ([]*FileInfo, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return nil, err
	}
	work, err := a.eng.WorkDir(ex)
	if err != nil {
		return nil, err
	}
	out := make([]*FileInfo, 0)
	err = filepath.WalkDir(work, func(p string, de os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if p == work {
			return nil
		}
		rel, rerr := filepath.Rel(work, p)
		if rerr != nil {
			return nil
		}
		relSlash := filepath.ToSlash(rel)
		if strings.HasPrefix(relSlash, ".") {
			return nil // skip hidden entries (node_modules, .git etc)
		}
		fi, ferr := de.Info()
		if ferr != nil {
			return nil
		}
		out = append(out, &FileInfo{
			Path: relSlash,
			Name: de.Name(),
			Size: fi.Size(),
			Dir:  de.IsDir(),
		})
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, err
}

// ReadFile returns the UTF-8 content of one file in the exercise workspace.
func (a *App) ReadFile(exID, rel string) (string, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return "", err
	}
	work, err := a.eng.WorkDir(ex)
	if err != nil {
		return "", err
	}
	p, err := safeJoin(work, rel)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile saves UTF-8 content into the exercise workspace (creating
// missing parent directories). Refuses paths escaping the workspace.
func (a *App) WriteFile(exID, rel, content string) error {
	ex, err := a.findExercise(exID)
	if err != nil {
		return err
	}
	work, err := a.eng.WorkDir(ex)
	if err != nil {
		return err
	}
	p, err := safeJoin(work, rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(content), 0o644)
}