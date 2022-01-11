package ng

import (
	"bufio"
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo-ng/lib/proc"
	zglob "github.com/mattn/go-zglob"
	"go/build"
	"golang.org/x/mod/module"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ImportMod struct {
	ImportPath string
	Version    string
	Dir        string          // full path, $GOPATH/pkg/mod/
	Pkgs       []string        // sub-pkg import paths
	VendorList map[string]bool // files to vendor
}

func CopyModule(in string) {
	f, err := os.Open(file.WorkingDir() + "/go.mod")
	if err != nil {
		proc.Die("unable to open go.mod")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	run("go", "mod", "vendor")
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var mod *ImportMod
	var replace = map[string]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		s := strings.Fields(line)
		if len(s) == 0 {
			continue
		}
		if s[0] == "replace" && len(s) == 4 {
			replace[s[1]] = s[3]
			continue
		}
		if s[0] == in && len(s) >= 2 {

			mod = &ImportMod{
				ImportPath: s[0],
				Version:    s[1],
				Pkgs:       []string{},
			}

			if v, ok := replace[mod.ImportPath]; ok {
				mod.ImportPath = strings.Replace(v, "/", "/", -1)
				mod.Version = ""
			} else {
				mod.ImportPath, err = GetModulePath(mod.ImportPath, mod.Version)
				if err != nil {
					evo.Panic(err)
				}
			}

			matches, err := zglob.Glob(filepath.Join(mod.ImportPath, "**/*"))
			if err != nil {
				evo.Panic("unable to open " + mod.ImportPath)
			}

			var dir = filepath.Join(file.WorkingDir(), "vendor", in)

			var skipList = [...]string{
				".", "vendor", "examples",
			}
			for i, item := range skipList {
				skipList[i] = strings.TrimRight(mod.ImportPath, "/") + "/" + item
			}

			var st = len(mod.ImportPath)
			for _, match := range matches {
				var shouldSkip = false
				for _, skip := range skipList {
					if strings.HasPrefix(match, skip) {
						shouldSkip = true
						break
					}
				}
				if shouldSkip {
					continue
				}
				var _, err = copyFile(match, filepath.Join(dir, match[st:]))
				if err != nil {
					evo.Panic(err)
				}
			}

			break
		}
	}

}

func GetModulePath(name, version string) (string, error) {
	// first we need GOMODCACHE

	cache, ok := os.LookupEnv("GOMODCACHE")
	if !ok {
		if os.Getenv("GOPATH") == "" {
			cache = path.Join(build.Default.GOPATH, "pkg", "mod")
		} else {
			cache = path.Join(os.Getenv("GOPATH"), "pkg", "mod")
		}

	}

	// then we need to escape path
	escapedPath, err := module.EscapePath(name)
	if err != nil {
		return "", err
	}

	// version also
	escapedVersion, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}

	return path.Join(cache, escapedPath+"@"+escapedVersion), nil
}

func copyFile(src, dst string) (int64, error) {
	srcStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if srcStat.IsDir() {
		_ = file.MakePath(dst)
		return 0, nil
	}
	if !srcStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			return
		}
	}(srcFile)

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			return
		}
	}(dstFile)

	return io.Copy(dstFile, srcFile)
}
