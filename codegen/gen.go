package codegen

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strconv"
	"strings"

	"github.com/YiCodes/gocode"
)

type ParseContext struct {
	imports     map[string]string
	webAPIList  []*webAPI
	packageName string
	pkg         *types.Package
}

func NewParseContext() *ParseContext {
	c := &ParseContext{}
	c.imports = make(map[string]string)
	c.imports["github.com/YiCodes/goweb/web"] = "web"
	c.imports["net/http"] = "http"

	return c
}

func (c *ParseContext) addWebAPI(api *webAPI) {
	c.webAPIList = append(c.webAPIList, api)
}

func (c *ParseContext) getPackageNameByPath(path string, packageName string) string {
	p, ok := c.imports[path]

	if ok {
		return p
	}

	c.imports[path] = packageName

	return packageName
}

type webAPIParameter struct {
	name        string
	isPointer   bool
	paramType   string
	isHTTPParam bool
	tempVarName string
}

type webAPI struct {
	name string

	params  []*webAPIParameter
	returns []*webAPIParameter

	paramCount  int
	returnCount int
}

func getHandleFuncDecl(file *ast.File) []*ast.FuncDecl {
	var funcList []*ast.FuncDecl

	for _, d := range file.Decls {
		f, ok := d.(*ast.FuncDecl)

		if ok {
			c := f.Name.Name[0]

			if c >= 'A' && c <= 'Z' {
				funcList = append(funcList, f)
			}
		}
	}

	return funcList
}

func Compile(parseContext *ParseContext, dest string) error {
	outFile, err := os.Create(dest)

	if err != nil {
		return err
	}

	defer outFile.Close()

	codeWriter := gocode.NewCodeWriter(outFile)
	codeWriter.WriteLine("package ", parseContext.pkg.Name())
	codeWriter.WriteLine()
	codeWriter.Write("import")
	codeWriter.BeginBlock("(")

	for k, v := range parseContext.imports {
		if strings.HasSuffix(k, v) {
			codeWriter.WriteLine(`"`, k, `"`)
		} else {
			codeWriter.WriteLine(v, " ", `"`, k, `"`)
		}
	}

	codeWriter.EndBlock(")")
	codeWriter.WriteLine()

	codeWriter.Write("func ConfigServeMuxHandler(mux *http.ServeMux, opts *web.HandlerSetupOptions)")
	codeWriter.BeginBlock("{")
	codeWriter.WriteLine("routeMap := opts.Route")
	codeWriter.WriteLine("msgCodec := opts.Codec")
	codeWriter.WriteLine("onRequestError := opts.OnRequestError")
	codeWriter.WriteLine("onResponseError := opts.OnResponseError")

	for _, api := range parseContext.webAPIList {
		codeWriter.Write(`mux.HandleFunc(routeMap.GetRoute("`)
		codeWriter.Write(api.name)
		codeWriter.Write(`"), func(w http.ResponseWriter, r *http.Request)`)
		codeWriter.BeginBlock("{")

		for _, apiParam := range api.params {
			if apiParam.isHTTPParam {
				continue
			}

			codeWriter.Write("var ")
			codeWriter.Write(apiParam.tempVarName)

			if apiParam.isPointer {
				codeWriter.WriteLine(" = new(", apiParam.paramType, ")")
			} else {
				codeWriter.WriteLine(" ", apiParam.paramType)
			}
		}

		if api.paramCount > 0 || api.returnCount > 0 {
			codeWriter.WriteLine("var err error")
		}

		if api.paramCount > 0 {
			codeWriter.Write("err = msgCodec.Decode(r")

			for _, apiParam := range api.params {
				if apiParam.isHTTPParam {
					continue
				}

				codeWriter.Write(", ")

				if !apiParam.isPointer {
					codeWriter.Write("&")
				}

				codeWriter.Write(apiParam.tempVarName)
			}

			codeWriter.WriteLine(")")

			codeWriter.Write("if err != nil")
			codeWriter.BeginBlock("{")
			codeWriter.WriteLine("onRequestError(r, w, err)")
			codeWriter.WriteLine("return")
			codeWriter.EndBlock("}")
		}

		if api.returnCount > 0 {
			for i, r := range api.returns {
				if i > 0 {
					codeWriter.Write(", ")
				}

				codeWriter.Write(r.tempVarName)
			}

			codeWriter.Write(" := ")
		}

		codeWriter.Write(api.name)
		codeWriter.Write("(")

		for i, apiParam := range api.params {
			if i > 0 {
				codeWriter.Write(", ")
			}

			codeWriter.Write(apiParam.tempVarName)
		}

		codeWriter.WriteLine(")")

		if api.returnCount > 0 {
			codeWriter.Write("err = msgCodec.Encode(w")

			for _, r := range api.returns {
				codeWriter.Write(", ")

				if !r.isPointer {
					codeWriter.Write("&")
				}

				codeWriter.Write(r.tempVarName)
			}

			codeWriter.WriteLine(")")

			codeWriter.Write("if err != nil")
			codeWriter.BeginBlock("{")
			codeWriter.WriteLine("onResponseError(r, err)")
			codeWriter.WriteLine("return")
			codeWriter.EndBlock("}")
		}

		codeWriter.EndBlock("})")
	}

	codeWriter.EndBlock("}")

	return nil
}

func isHttpRequest(t types.Type) bool {
	switch inst := t.(type) {
	case *types.Named:
		obj := inst.Obj()

		if obj.Pkg().Path() == "net/http" && obj.Name() == "Request" {
			return true
		}
	}

	return false
}

func (c *ParseContext) getTypeName(t types.Type) (string, error) {
	switch inst := t.(type) {
	case *types.Named:
		obj := inst.Obj()

		if c.pkg == obj.Pkg() {
			return obj.Name(), nil
		}

		packageName := c.getPackageNameByPath(obj.Pkg().Path(), obj.Pkg().Name())
		return fmt.Sprintf("%s.%s", packageName, obj.Name()), nil

	case *types.Basic:
		return inst.String(), nil
	}

	return "", fmt.Errorf("this type not supported")
}

func Parse(parseContext *ParseContext, srcFiles []string) error {
	srcFileSet := token.NewFileSet()

	var srcAstFiles []*ast.File

	for _, sf := range srcFiles {
		f, err := parser.ParseFile(srcFileSet, sf, nil, 0)

		if err != nil {
			return err
		}

		srcAstFiles = append(srcAstFiles, f)
	}

	typesInfo := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	var typesConf types.Config
	typesConf.Importer = importer.Default()
	pkg, err := typesConf.Check("", srcFileSet, srcAstFiles, &typesInfo)

	if err != nil {
		return err
	}

	parseContext.pkg = pkg

	for _, srcAstFile := range srcAstFiles {
		for _, f := range getHandleFuncDecl(srcAstFile) {
			o := typesInfo.Defs[f.Name].(*types.Func)

			if !o.Exported() {
				continue
			}

			api := &webAPI{}
			api.name = f.Name.Name

			if f.Type.Params != nil {
				for i, p := range f.Type.Params.List {
					typeInfo := typesInfo.TypeOf(p.Type)
					typeInfo = typeInfo.Underlying()

					webAPIParam := &webAPIParameter{}
					webAPIParam.name = p.Names[0].Name
					webAPIParam.tempVarName = "a" + strconv.Itoa(i)

					switch instTypeInfo := typeInfo.(type) {
					case *types.Pointer:
						webAPIParam.isPointer = true

						if isHttpRequest(instTypeInfo.Elem()) {
							webAPIParam.isHTTPParam = true
							webAPIParam.tempVarName = "r"
						}

						elemTypeName, err := parseContext.getTypeName(instTypeInfo.Elem())

						if err != nil {
							return fmt.Errorf("%v(%v)", err, srcFileSet.Position(p.Type.Pos()))
						}

						webAPIParam.paramType = elemTypeName
					case *types.Interface:
						t, ok := typesInfo.TypeOf(p.Type).(*types.Named)

						if ok {
							if t.Obj().Name() == "ResponseWriter" &&
								t.Obj().Pkg().Path() == "net/http" {
								webAPIParam.tempVarName = "w"
								webAPIParam.isHTTPParam = true
							}
						} else {
							return fmt.Errorf("interface not supported(%v)", srcFileSet.Position(p.Type.Pos()))
						}

					case *types.Basic:
						webAPIParam.isPointer = false
						webAPIParam.paramType = instTypeInfo.Name()
					case *types.Struct:
						typeName, err := parseContext.getTypeName(typesInfo.TypeOf(p.Type))

						if err != nil {
							return fmt.Errorf("%v(%v)", err, srcFileSet.Position(p.Type.Pos()))
						}

						webAPIParam.paramType = typeName
					default:
						return fmt.Errorf("this type not supported(%v)", srcFileSet.Position(p.Type.Pos()))
					}

					if !webAPIParam.isHTTPParam {
						api.paramCount++
					}
					api.params = append(api.params, webAPIParam)
				}
			}

			if f.Type.Results != nil {
				for i, r := range f.Type.Results.List {
					typeInfo := typesInfo.TypeOf(r.Type).Underlying()

					webAPIParam := &webAPIParameter{}
					webAPIParam.tempVarName = "r" + strconv.Itoa(i)

					switch typeInfo.(type) {
					case *types.Pointer:
						webAPIParam.isPointer = true
					}

					api.returnCount++
					api.returns = append(api.returns, webAPIParam)
				}
			}

			parseContext.addWebAPI(api)
		}
	}

	return nil
}
