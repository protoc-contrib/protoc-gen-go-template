package generator_test

import (
	"github.com/protoc-contrib/protoc-gen-template/internal/generator"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options.Set", func() {
	var opts *generator.Options

	BeforeEach(func() {
		opts = &generator.Options{}
	})

	It("accepts template_dir", func() {
		Expect(opts.Set("template_dir", "/tmp/tmpl")).To(Succeed())
		Expect(opts.TemplateDir).To(Equal("/tmp/tmpl"))
	})

	It("accepts destination_dir", func() {
		Expect(opts.Set("destination_dir", "gen")).To(Succeed())
		Expect(opts.DestinationDir).To(Equal("gen"))
	})

	DescribeTable("boolean flags",
		func(name string, getter func(*generator.Options) bool) {
			Expect(opts.Set(name, "true")).To(Succeed())
			Expect(getter(opts)).To(BeTrue())

			opts2 := &generator.Options{}
			Expect(opts2.Set(name, "false")).To(Succeed())
			Expect(getter(opts2)).To(BeFalse())

			opts3 := &generator.Options{}
			Expect(opts3.Set(name, "t")).To(Succeed())
			Expect(getter(opts3)).To(BeTrue())
		},
		Entry("debug", "debug", func(o *generator.Options) bool { return o.Debug }),
		Entry("all", "all", func(o *generator.Options) bool { return o.All }),
		Entry("single-package-mode", "single-package-mode", func(o *generator.Options) bool { return o.SinglePackageMode }),
		Entry("file-mode", "file-mode", func(o *generator.Options) bool { return o.FileMode }),
	)

	It("rejects unknown options", func() {
		err := opts.Set("nope", "1")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown plugin option"))
	})

	It("rejects malformed booleans", func() {
		err := opts.Set("debug", "yeah")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid value"))
	})
})
