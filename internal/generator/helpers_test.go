package generator_test

import (
	"fmt"
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/protoc-contrib/protoc-gen-go-template/internal/generator"
)

// render invokes a template that calls `name` against args, passed in as
// `.` so each argument is reachable as `(index . N)`.
func render(name string, args ...any) (string, error) {
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = fmt.Sprintf("(index . %d)", i)
	}
	src := fmt.Sprintf("{{ %s %s }}", name, strings.Join(placeholders, " "))
	tmpl, err := template.New("t").Funcs(generator.ProtoHelpersFuncMap).Parse(src)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, args); err != nil {
		return "", err
	}
	return sb.String(), nil
}

var _ = Describe("funcmap", func() {
	DescribeTable("naming helpers",
		func(name, input, want string) {
			got, err := render(name, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want))
		},
		Entry("camelCase hello_world", "camelCase", "hello_world", "helloWorld"),
		Entry("camelCase foo_bar_baz", "camelCase", "foo_bar_baz", "fooBarBaz"),
		Entry("camelCase single char", "camelCase", "a", "A"),
		Entry("lowerCamelCase hello_world", "lowerCamelCase", "hello_world", "helloWorld"),
		Entry("lowerCamelCase foo_bar", "lowerCamelCase", "foo_bar", "fooBar"),
		Entry("lowerCamelCase single char", "lowerCamelCase", "a", "a"),
		Entry("snakeCase HelloWorld", "snakeCase", "HelloWorld", "hello_world"),
		Entry("snakeCase FooBarBaz", "snakeCase", "FooBarBaz", "foo_bar_baz"),
		Entry("kebabCase HelloWorld", "kebabCase", "HelloWorld", "hello-world"),
		Entry("kebabCase foo_bar", "kebabCase", "foo_bar", "foo-bar"),
		Entry("upperFirst", "upperFirst", "hello", "Hello"),
		Entry("lowerFirst", "lowerFirst", "HELLO", "hELLO"),
		Entry("upperCase", "upperCase", "hello", "HELLO"),
	)

	DescribeTable("arithmetic helpers",
		func(name string, a, b, want int) {
			got, err := render(name, a, b)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(fmt.Sprint(want)))
		},
		Entry("add", "add", 2, 3, 5),
		Entry("add negative", "add", -4, 1, -3),
		Entry("subtract", "subtract", 10, 4, 6),
		Entry("multiply", "multiply", 6, 7, 42),
		Entry("multiply by zero", "multiply", 5, 0, 0),
		Entry("divide", "divide", 20, 5, 4),
		Entry("divide exact", "divide", 9, 3, 3),
	)

	It("rejects division by zero at template-execution time", func() {
		_, err := render("divide", 1, 0)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("divide"))
	})

	DescribeTable("string predicates",
		func(name string, args []any, want string) {
			got, err := render(name, args...)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want))
		},
		Entry("contains present", "contains", []any{"ell", "hello"}, "true"),
		Entry("contains absent", "contains", []any{"xyz", "hello"}, "false"),
		Entry("trimstr strips prefix+suffix", "trimstr", []any{"/", "/foo/bar/"}, "foo/bar"),
	)

	It("renders json and prettyjson", func() {
		out, err := render("json", map[string]int{"a": 1})
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(`{"a":1}`))

		out, err = render("prettyjson", map[string]int{"a": 1})
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("\n  \"a\": 1"))
	})

	It("splits strings into a list of items", func() {
		src := `{{ range $i, $v := splitArray "," "a,b,,c" }}{{if $i}},{{end}}{{$v}}{{end}}`
		tmpl, err := template.New("t").Funcs(generator.ProtoHelpersFuncMap).Parse(src)
		Expect(err).NotTo(HaveOccurred())
		var sb strings.Builder
		Expect(tmpl.Execute(&sb, nil)).To(Succeed())
		// splitArray drops empty segments, so "a,b,,c" → ["a","b","c"].
		Expect(sb.String()).To(Equal("a,b,c"))
	})

	DescribeTable("proto type helpers",
		func(name, input, want string) {
			got, err := render(name, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want))
		},
		Entry("shortType strips package prefix", "shortType", ".foo.bar.Baz", "Baz"),
		Entry("shortType passthrough when no dots", "shortType", "Baz", "Baz"),
		Entry("namespacedFlowType dollar-delimits", "namespacedFlowType", ".foo.bar.Baz", "foo$bar$Baz"),
		Entry("goNormalize snake→PascalCamel", "goNormalize", "foo_bar_baz", "fooBarBaz"),
		Entry("goNormalize rewrites id→ID", "goNormalize", "user_id", "userID"),
		Entry("lowerGoNormalize", "lowerGoNormalize", "foo_bar_baz", "fooBarBaz"),
	)

	DescribeTable("jsSuffixReserved",
		func(input, want string) {
			got, err := render("jsSuffixReserved", input)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want))
		},
		Entry("appends _ to reserved", "class", "class_"),
		Entry("appends _ to another reserved", "function", "function_"),
		Entry("passthrough non-reserved", "userName", "userName"),
	)

	It("replaceDict substitutes from a dict", func() {
		src := `{{ replaceDict "hello world" (dict "hello" "hi" "world" "earth") }}`
		tmpl, err := template.New("t").Funcs(generator.ProtoHelpersFuncMap).Parse(src)
		Expect(err).NotTo(HaveOccurred())
		var sb strings.Builder
		Expect(tmpl.Execute(&sb, nil)).To(Succeed())
		Expect(sb.String()).To(Equal("hi earth"))
	})

	It("registers the expected set of helpers", func() {
		wantKeys := []string{
			"camelCase", "lowerCamelCase", "snakeCase", "kebabCase",
			"upperFirst", "lowerFirst", "upperCase",
			"add", "subtract", "multiply", "divide",
			"contains", "trimstr", "json", "prettyjson",
			"goType", "goPkg", "jsType", "haskellType",
			"httpVerb", "httpPath", "httpBody",
			"first", "last", "splitArray",
		}
		for _, k := range wantKeys {
			Expect(generator.ProtoHelpersFuncMap).To(HaveKey(k), "missing helper %q", k)
		}
	})
})
