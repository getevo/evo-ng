package ng

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/internal/file"
	"github.com/getevo/evo-ng/internal/regex"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Package struct {
	IsLocal   bool
	Structs   []Struct
	Name      string
	Path      string
	Imports   []string
	Comments  []string
	Functions []Function
	Mod       Mod
}

type Mod struct {
	Version string
	Commit  string
	Date    time.Time
}

type Struct struct {
	Name      string
	Fields    []Field
	Functions []Function
}

type Field struct {
	Name string
	Type string
	Tag  string
}

type Function struct {
	Name     string
	Extend   Type
	Extended bool
	Input    []Type
	Result   []Type
}

type Type struct {
	IsPtr  bool
	Var    string
	Struct string
	Pkg    string
}

func ParsePackage(path string) *Package {
	var result = Package{}
	var err = result.FindPackage(path)
	if err != nil {
		evo.Panic(err)
	}
	if result.Path == "" {
		panic(evo.STrace("unable to find "+path, 0))
	}
	files, err := ioutil.ReadDir(result.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		src, err := ioutil.ReadFile(result.Path + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, "", string(src), parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		/*		conf := types.Config{Importer: importer.Default()}
				pkg, err := conf.Check("cmd", fs, []*ast.File{file}, nil)
				scope := pkg.Scope()*/

		result.Name = file.Name.String()
		for _, comment := range file.Comments {
			for _, item := range comment.List {
				result.Comments = append(result.Comments, item.Text)
			}
		}

		for _, decl := range file.Scope.Objects {
			if typeSpec, ok := decl.Decl.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					st := Struct{
						Name: decl.Name,
					}
					if structType.Fields != nil {

						for _, fieldTypeSpec := range structType.Fields.List {

							if len(fieldTypeSpec.Names) > 0 {
								var f = Field{
									Name: fieldTypeSpec.Names[0].Name,
								}
								if fieldTypeSpec.Tag != nil {
									f.Tag = fieldTypeSpec.Tag.Value
								}
								if ident, ok := fieldTypeSpec.Type.(*ast.Ident); ok {
									f.Type = ident.Name
								}

								st.Fields = append(st.Fields, f)
							}

						}
					}
					result.Structs = append(result.Structs, st)
				}

			}
		}

		for _, decl := range file.Decls {
			switch t := decl.(type) {
			// That's a func decl !
			case *ast.FuncDecl:
				var fn = Function{}
				fn.Name = t.Name.Name
				if t.Type != nil {
					if t.Type.Params != nil {
						for _, param := range t.Type.Params.List {
							typ := getType(param.Type, result.Name)
							fn.Input = append(fn.Input, typ)
						}
					}
					if t.Type.Results != nil {
						for _, param := range t.Type.Results.List {
							typ := getType(param.Type, result.Name)
							fn.Result = append(fn.Result, typ)
						}
					}
				}
				if t.Recv != nil && len(t.Recv.List) == 1 {
					fn.Extended = true
					fn.Extend = getType(t.Recv.List[0].Type, result.Name)
					fn.Extend.Pkg = result.Name
					for index, s := range result.Structs {
						if s.Name == fn.Extend.Struct {
							result.Structs[index].Functions = append(result.Structs[index].Functions, fn)
						}
					}
				} else {
					result.Functions = append(result.Functions, fn)
				}

			}
		}

	}
	return &result
}

func (p *Package) FindPackage(path string) error {
	var chunks = strings.Split(path, "@")
	var repo = ""
	var tag = ""
	var branch = ""
	if len(chunks) == 2 {
		repo = chunks[0]
		tag = chunks[1]
		chunks = strings.Split(tag, ":")
		if len(chunks) == 2 {
			tag = chunks[0]
			branch = chunks[1]
		}
		p.IsLocal = false
	} else {
		repo = chunks[0]
		tag = ""
		branch = ""
	}

	if file.IsDirExist("./" + path) {
		p.Path = "./" + path
		p.IsLocal = true
		return nil
	}

	if file.IsDirExist("./vendor/" + path) {
		p.Path = "./vendor/" + path
		p.IsLocal = false
		return nil
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	if tag == "branch" || tag == "dev" {
		file.MakePath(gopath + "/src/" + filepath.Dir(repo))
		if !file.IsDirExist(gopath + "/src/" + repo) {
			fmt.Println("Cloning", repo, "into", file.WorkingDir()+"/git/"+repo)
			repository, err := git.PlainClone(gopath+"/src/"+repo, false, &git.CloneOptions{
				URL:      "https://" + repo,
				Progress: os.Stdout,
			})
			if err != nil {
				file.Remove(gopath + "/src/" + repo)
				panic(err)
			}
			w, err := repository.Worktree()
			if err != nil {
				file.Remove(gopath + "/src/" + repo)
				panic(err)
			}
			err = repository.Fetch(&git.FetchOptions{
				RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
			})
			if err != nil {
				file.Remove(gopath + "/src/" + repo)
				panic(err)
			}
			err = w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
				Force:  true,
			})
			if err != nil {
				file.Remove(gopath + "/src/" + repo)
				panic(err)
			}

		}
		run("go", "mod", "edit", "-replace", repo+"="+gopath+"/src/"+repo)
		p.Path = gopath + "/src/" + repo
		p.IsLocal = false
		return nil
	} else if regex.VersionedPackage.MatchString(repo) {
		//versioned + no tag
		p.IsLocal = false
		run("go", "get", "-d", path)
		var data, err = ioutil.ReadFile("./go.mod")
		if err != nil {
			panic(err)
		}
		f, err := modfile.Parse("go.mod", data, nil)
		if err != nil {
			panic(err)
		}

		for _, item := range f.Require {
			if item.Mod.Path == repo {
				p.Path = filepath.Join(gopath, "pkg", "mod", item.Mod.Path+"@"+item.Mod.Version)
				return nil
			}

		}

	}

	return fmt.Errorf("unable to find package %s", path)
}

func getType(input ast.Expr, pkg string) Type {
	var t = Type{}

	switch exp := input.(type) {
	case *ast.StarExpr:
		if x, ok := exp.X.(*ast.Ident); ok {
			if x.Obj != nil {
				t = Type{
					IsPtr:  true,
					Var:    "",
					Struct: x.Obj.Name,
					Pkg:    pkg,
				}
			} else {
				t = Type{
					IsPtr:  true,
					Var:    "",
					Struct: x.Name,
					Pkg:    pkg,
				}
			}

		} else if x, ok := exp.X.(*ast.SelectorExpr); ok {
			t = getType(x, pkg)
			t.IsPtr = true
		}
	case *ast.Ident:
		t = Type{
			IsPtr:  false,
			Var:    "",
			Struct: exp.Name,
		}
	case *ast.SelectorExpr:

		if x, ok := exp.X.(*ast.Ident); ok {
			t = Type{
				IsPtr:  false,
				Var:    "",
				Struct: exp.Sel.Name,
				Pkg:    x.Name,
			}
		}
	default:

	}
	return t
}

func (p Package) HasFunction(in Function) bool {
	for _, fn := range p.Functions {
		if fn.Name == in.Name {
			if len(in.Input) != len(fn.Input) || len(in.Result) != len(fn.Result) {
				return false
			}

			for k, src := range in.Input {
				var dst = fn.Input[k]
				if src.IsPtr != dst.IsPtr || src.Pkg != dst.Pkg || src.Struct != dst.Struct {
					return false
				}
			}
			for k, src := range in.Result {
				var dst = fn.Result[k]
				if src.IsPtr != dst.IsPtr || src.Pkg != dst.Pkg || src.Struct != dst.Struct {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (s Struct) HasFunction(in Function) bool {
	for _, fn := range s.Functions {
		if fn.Name == in.Name && fn.Extend.IsPtr == in.Extend.IsPtr {
			if len(in.Input) != len(fn.Input) || len(in.Result) != len(fn.Result) {
				return false
			}

			for k, src := range in.Input {
				var dst = fn.Input[k]
				if src.IsPtr != dst.IsPtr || src.Pkg != dst.Pkg || src.Struct != dst.Struct {
					return false
				}
			}
			for k, src := range in.Result {
				var dst = fn.Result[k]
				if src.IsPtr != dst.IsPtr || src.Pkg != dst.Pkg || src.Struct != dst.Struct {
					return false
				}
			}
			return true
		}
	}
	return false
}
