package export

import (
	"fmt"
	"strings"

	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/config"
	"google.golang.org/protobuf/compiler/protogen"
)

func Export(gen *protogen.Plugin, conf []config.ServiceConfig) {
	g := gen.NewGeneratedFile(
		"../../grpc-rest-gateway/output.go",
		protogen.GoImportPath("github.com/HomayoonAlimohammadi/mini-grpc-gateway"),
	)

	g.P(("package main"))
	g.P(fmt.Sprintf("// %+v", conf))

	exportImports(g)

	exportMain(g, conf)

	exportServiceHandlers(g, conf)

	exportWriteError(g)
}

func exportImports(g *protogen.GeneratedFile) {
	g.P("import (")
	g.P("\"log\"")
	g.P("\"fmt\"")
	g.P("\"net/http\"")
	g.P("\"time\"")
	g.P("\"google.golang.org/grpc/credentials/insecure\"")
	g.P("\"github.com/HomayoonAlimohammadi/mini-grpc-gateway/pb/post\"")
	g.P("\"github.com/gorilla/mux\"")
	g.P("\"google.golang.org/grpc\"")
	g.P("\"google.golang.org/protobuf/encoding/protojson\"")
	g.P(")")
}

func exportMain(g *protogen.GeneratedFile, svcConf []config.ServiceConfig) {
	g.P("func main() {")
	g.P("r := mux.NewRouter()")

	exportMethodHandlers(g, svcConf)

	g.P("srv := &http.Server {")
	g.P("Handler: r,")
	g.P("Addr: \":8000\",")
	g.P("WriteTimeout: 15 * time.Second,")
	g.P("ReadTimeout: 15 * time.Second,")
	g.P("}")

	g.P("log.Print(\"serving on :8000\")")

	g.P("log.Fatal(srv.ListenAndServe())")

	g.P("}")
}

func exportMethodHandlers(g *protogen.GeneratedFile, svcConf []config.ServiceConfig) {
	for _, sc := range svcConf {
		svcHVarName := getSvcHandlerInitializerVarName(sc.GoName)
		initalizerName := getSvcHandlerInitializerName(sc.GoName)

		g.P(fmt.Sprintf("%s, err := %s()", svcHVarName, initalizerName))
		g.P("if err != nil {")
		g.P(fmt.Sprintf("log.Fatalf(\"failed to get %s handler: %%v\", err)", sc.GoName))
		g.P("}")

		for _, mc := range sc.Methods {
			methodHandlerName := getMethodHandlerName(mc.GoName)
			g.P(fmt.Sprintf("r.HandleFunc(\"%s\", %s.%s)",
				mc.Options.Url, svcHVarName, methodHandlerName))
		}
	}
}

func exportMethodsHandlers(g *protogen.GeneratedFile, methodConf config.MethodConfig, svcConf config.ServiceConfig) {
	exportInputMaker(g, methodConf, svcConf)

	svcHandlertructName := getSvcHandlerStructName(svcConf.GoName)
	handlerName := getMethodHandlerName(methodConf.GoName)

	g.P(fmt.Sprintf("func (sh *%s) %s(w http.ResponseWriter, r *http.Request) {", svcHandlertructName, handlerName))

	g.P("ctx := r.Context()")
	g.P(fmt.Sprintf("in := %s()", getMethodRequestMakerName(methodConf.GoName)))

	g.P(fmt.Sprintf("resp, err := sh.client.%s(ctx, in)", methodConf.GoName))
	g.P("if err != nil {")
	g.P("writeError(")
	g.P(fmt.Sprintf("w, fmt.Errorf(\"failed to call %s: %%w\", err),", methodConf.GoName))
	g.P("http.StatusInternalServerError,")
	g.P(")")
	g.P("return")
	g.P("}")

	g.P("b, err := protojson.Marshal(resp)")
	g.P("if err != nil {")
	g.P("writeError(")
	g.P("w, fmt.Errorf(\"failed to convert response to json: %w\", err),")
	g.P("http.StatusInternalServerError,")
	g.P(")")
	g.P("return")
	g.P("}")

	g.P("w.Write(b)")

	g.P("}")
}

func exportInputMaker(g *protogen.GeneratedFile, mc config.MethodConfig, sc config.ServiceConfig) {
	g.P(fmt.Sprintf("func %s() *%s.%s {", getMethodRequestMakerName(mc.GoName), sc.Package, mc.In.GoIdent.GoName))

	g.P(fmt.Sprintf("return &%s.%s{}", sc.Package, mc.In.GoIdent.GoName))

	g.P("}")
}

func exportServiceHandlers(g *protogen.GeneratedFile, svcConf []config.ServiceConfig) {
	for _, sc := range svcConf {
		structName := getSvcHandlerStructName(sc.GoName)
		initializerName := getSvcHandlerInitializerName(sc.GoName)

		g.P(fmt.Sprintf("type %s struct {", structName))
		g.P(fmt.Sprintf("client %s.%sClient", sc.Package, sc.GoName))
		g.P("}")

		g.P(fmt.Sprintf("func %s() (*%s, error) {", initializerName, structName))

		g.P(fmt.Sprintf("conn, err := grpc.Dial(\":%d\", grpc.WithTransportCredentials(insecure.NewCredentials()))", sc.GPRCPort))
		g.P("if err != nil {")
		g.P(fmt.Sprintf("return nil, fmt.Errorf(\"failed to dial %d: %%w\", err)", sc.GPRCPort))
		g.P("}")

		g.P(fmt.Sprintf("return &%s{", structName))
		g.P(fmt.Sprintf("client: %s.New%sClient(conn),", sc.Package, sc.GoName))
		g.P("}, nil")

		g.P("}")

		for _, mc := range sc.Methods {
			exportMethodsHandlers(g, mc, sc)
		}
	}
}

func unexportedGoName(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}

func getSvcHandlerStructName(goName string) string {
	return unexportedGoName(goName)
}

func getSvcHandlerInitializerName(goName string) string {
	return fmt.Sprintf("Get%sHandler", goName)
}

func getSvcHandlerInitializerVarName(goName string) string {
	return fmt.Sprintf("%sHandler", goName)
}

func getMethodHandlerName(goName string) string {
	return fmt.Sprintf("%sHandler", goName)
}

func getMethodRequestMakerName(goName string) string {
	return fmt.Sprintf("make%sRequest", goName)
}

func exportWriteError(g *protogen.GeneratedFile) {
	g.P("func writeError(w http.ResponseWriter, err error, status int) {")
	g.P("w.WriteHeader(status)")
	g.P("w.Write([]byte(err.Error()))")
	g.P("}")
}
