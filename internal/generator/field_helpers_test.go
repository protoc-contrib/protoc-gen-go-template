package generator_test

import (
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/protoc-contrib/protoc-gen-go-template/internal/generator"
)

// scalarField returns a FieldDescriptorProto for a primitive scalar.
func scalarField(name string, t descriptorpb.FieldDescriptorProto_Type, repeated bool) *descriptorpb.FieldDescriptorProto {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	if repeated {
		label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	}
	return &descriptorpb.FieldDescriptorProto{
		Name:  proto.String(name),
		Type:  t.Enum(),
		Label: label.Enum(),
	}
}

// messageField returns a FieldDescriptorProto of message type referring to
// the given fully-qualified message type name (e.g. ".demo.Bar").
func messageField(name, typeName string, repeated bool) *descriptorpb.FieldDescriptorProto {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	if repeated {
		label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	}
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		Label:    label.Enum(),
		TypeName: proto.String(typeName),
	}
}

// renderWithData executes a tiny template using the funcmap against `data`.
func renderWithData(src string, data any) (string, error) {
	tmpl, err := template.New("t").Funcs(generator.ProtoHelpersFuncMap).Parse(src)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

var _ = Describe("descriptor-walking helpers", func() {
	DescribeTable("goType on scalars",
		func(t descriptorpb.FieldDescriptorProto_Type, repeated bool, want string) {
			f := scalarField("x", t, repeated)
			out, err := renderWithData(`{{ goType "" . }}`, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(want))
		},
		Entry("int32", descriptorpb.FieldDescriptorProto_TYPE_INT32, false, "int32"),
		Entry("repeated int32", descriptorpb.FieldDescriptorProto_TYPE_INT32, true, "[]int32"),
		Entry("bool", descriptorpb.FieldDescriptorProto_TYPE_BOOL, false, "bool"),
		Entry("string", descriptorpb.FieldDescriptorProto_TYPE_STRING, false, "string"),
		Entry("repeated string", descriptorpb.FieldDescriptorProto_TYPE_STRING, true, "[]string"),
		Entry("float32", descriptorpb.FieldDescriptorProto_TYPE_FLOAT, false, "float32"),
		Entry("float64", descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, false, "float64"),
		Entry("uint64", descriptorpb.FieldDescriptorProto_TYPE_UINT64, false, "uint64"),
	)

	It("goType with package prefixes message types", func() {
		f := messageField("b", ".demo.Bar", false)
		out, err := renderWithData(`{{ goType "pkg" . }}`, f)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal("*pkg.Bar"))
	})

	It("goType on repeated messages emits slice-of-pointer", func() {
		f := messageField("b", ".demo.Bar", true)
		out, err := renderWithData(`{{ goType "" . }}`, f)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal("[]*Bar"))
	})

	DescribeTable("isFieldMessage / isFieldRepeated",
		func(f *descriptorpb.FieldDescriptorProto, wantMsg, wantRepeated string) {
			out, err := renderWithData(`{{ isFieldMessage . }}|{{ isFieldRepeated . }}`, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(wantMsg + "|" + wantRepeated))
		},
		Entry("scalar singular", scalarField("x", descriptorpb.FieldDescriptorProto_TYPE_STRING, false), "false", "false"),
		Entry("scalar repeated", scalarField("x", descriptorpb.FieldDescriptorProto_TYPE_STRING, true), "false", "true"),
		Entry("message singular", messageField("x", ".demo.Bar", false), "true", "false"),
		Entry("message repeated", messageField("x", ".demo.Bar", true), "true", "true"),
	)

	It("haskellType renders primitives", func() {
		f := scalarField("x", descriptorpb.FieldDescriptorProto_TYPE_STRING, false)
		out, err := renderWithData(`{{ haskellType "" . }}`, f)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("Text"))
	})
})
