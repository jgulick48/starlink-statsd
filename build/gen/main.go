// Command gen regenerates the Go protobuf bindings from proto/dish.protoset
// without protoc, using the protogen/gengo libraries directly.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const modPath = "github.com/jgulick48/starlink-statsd/pkg/"

var wkt = map[string]string{
	"google/protobuf/timestamp.proto":  "google.golang.org/protobuf/types/known/timestamppb;timestamppb",
	"google/protobuf/duration.proto":   "google.golang.org/protobuf/types/known/durationpb;durationpb",
	"google/protobuf/wrappers.proto":   "google.golang.org/protobuf/types/known/wrapperspb;wrapperspb",
	"google/protobuf/empty.proto":      "google.golang.org/protobuf/types/known/emptypb;emptypb",
	"google/protobuf/any.proto":        "google.golang.org/protobuf/types/known/anypb;anypb",
	"google/protobuf/struct.proto":     "google.golang.org/protobuf/types/known/structpb;structpb",
	"google/protobuf/field_mask.proto": "google.golang.org/protobuf/types/known/fieldmaskpb;fieldmaskpb",
	"google/protobuf/descriptor.proto": "google.golang.org/protobuf/types/descriptorpb;descriptorpb",
}

func goPackage(pkg string) string {
	parts := strings.Split(pkg, ".")
	if len(parts) > 2 && parts[0] == "SpaceX" && parts[1] == "API" {
		parts = parts[2:]
	}
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	if len(parts) == 0 || parts[0] == "" {
		parts = []string{"root"}
	}
	return modPath + "spacex.com/api/" + strings.Join(parts, "/") + ";" + parts[len(parts)-1]
}

func main() {
	in := flag.String("in", "proto/dish.protoset", "input FileDescriptorSet")
	out := flag.String("out", "pkg", "output directory")
	flag.Parse()
	if err := run(*in, *out); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run(in, out string) error {
	b, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	set := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(b, set); err != nil {
		return err
	}
	// Topologically sort (stable, deterministic).
	byName := map[string]*descriptorpb.FileDescriptorProto{}
	for _, f := range set.File {
		byName[f.GetName()] = f
	}
	var ordered []*descriptorpb.FileDescriptorProto
	seen := map[string]bool{}
	var visit func(n string) error
	visit = func(n string) error {
		if seen[n] {
			return nil
		}
		seen[n] = true
		f, ok := byName[n]
		if !ok {
			return fmt.Errorf("missing dependency %s", n)
		}
		for _, d := range f.Dependency {
			if err := visit(d); err != nil {
				return err
			}
		}
		ordered = append(ordered, f)
		return nil
	}
	for _, f := range set.File {
		if err := visit(f.GetName()); err != nil {
			return err
		}
	}

	req := &pluginpb.CodeGeneratorRequest{}
	for _, f := range ordered {
		if gp, ok := wkt[f.GetName()]; ok {
			if f.Options == nil {
				f.Options = &descriptorpb.FileOptions{}
			}
			f.Options.GoPackage = proto.String(gp)
		} else if strings.HasPrefix(f.GetName(), "google/protobuf/") {
			return fmt.Errorf("unmapped well-known type %s", f.GetName())
		} else {
			if f.Options == nil {
				f.Options = &descriptorpb.FileOptions{}
			}
			f.Options.GoPackage = proto.String(goPackage(f.GetPackage()))
			req.FileToGenerate = append(req.FileToGenerate, f.GetName())
		}
		req.ProtoFile = append(req.ProtoFile, f)
	}

	gen, err := protogen.Options{}.New(req)
	if err != nil {
		return err
	}
	for _, f := range gen.Files {
		if f.Generate {
			internal_gengo.GenerateFile(gen, f)
		}
	}
	resp := gen.Response()
	if resp.Error != nil {
		return fmt.Errorf("%s", resp.GetError())
	}
	if err := os.RemoveAll(filepath.Join(out, "spacex.com")); err != nil {
		return err
	}
	for _, rf := range resp.File {
		name := strings.TrimPrefix(rf.GetName(), modPath)
		p := filepath.Join(out, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(rf.GetContent()), 0o644); err != nil {
			return err
		}
	}
	fmt.Printf("generated %d files\n", len(resp.File))
	return nil
}
