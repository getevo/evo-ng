package ng

import (
	"bufio"
	"fmt"
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo-ng/lib/proc"
	zglob "github.com/mattn/go-zglob"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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
	defer f.Close()
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
				fmt.Println(mod.ImportPath)
				mod.Version = ""
			}
			matches, err := zglob.Glob(filepath.Join(mod.ImportPath, "**/*"))
			if err != nil {
				panic("unable to open " + mod.ImportPath)
			}

			var dir = filepath.Join(file.WorkingDir(), "vendor", in)

			var skipList = [...]string{
				mod.ImportPath + ".",
				mod.ImportPath + "vendor",
				mod.ImportPath + "examples",
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
					panic(err)
				}
			}

			break
		}
	}

}

func buildModVendorList(copyPattern []string, mod *ImportMod) map[string]bool {
	vendorList := map[string]bool{}

	for _, pat := range copyPattern {
		matches, err := zglob.Glob(filepath.Join(mod.Dir, pat))
		if err != nil {
			fmt.Println("Error! glob match failure:", err)
			os.Exit(1)
		}

		for _, m := range matches {
			vendorList[m] = false
		}
	}

	return vendorList
}

func pkgModPath(importPath, version string) string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		// the default GOPATH for go v1.11
		goPath = filepath.Join(os.Getenv("HOME"), "go")
	}

	var normPath string

	for _, char := range importPath {
		if unicode.IsUpper(char) {
			normPath += "!" + string(unicode.ToLower(char))
		} else {
			normPath += string(char)
		}
	}

	return filepath.Join(goPath, "pkg", "mod", fmt.Sprintf("%s@%s", normPath, version))
}

func copyFile(src, dst string) (int64, error) {
	srcStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if srcStat.IsDir() {
		file.MakePath(dst)
		fmt.Println("====>", dst)

		return 0, nil
	}
	if !srcStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
