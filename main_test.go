package semaphore_test

import (
	"testing"

	"github.com/ForestEckhardt/semaphore"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

func TestUnitFlagParser(t *testing.T) {
	suite := spec.New("flag-parser", spec.Report(report.Terminal{}))
	suite("FlagParser", testFlagParser)
	suite.Run(t)
}

func testFlagParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		flagParser semaphore.FlagParser
	)

	it.Before(func() {
		flagParser = semaphore.NewFlagParser()
	})

	it("parses a flag while escape spaces inside single and double quotes", func() {
		flags, err := flagParser.ParseFlags(`-buildmode=default -ldflags="-X main.variable=some-value" -ldflags='-X main.variable=some-other-value' -tags=paketo`)
		Expect(err).NotTo(HaveOccurred())

		Expect(flags).To(Equal([]string{
			`-buildmode=default`,
			`-ldflags="-X main.variable=some-value"`,
			`-ldflags='-X main.variable=some-other-value'`,
			`-tags=paketo`,
		}))
	})

	it("parses a flag while escape spaces inside single and double quotes while omitting additional white space", func() {
		flags, err := flagParser.ParseFlags(`-buildmode=default -ldflags="-X main.variable=some-value"     -ldflags='-X main.variable=some-other-value' -tags=paketo     `)
		Expect(err).NotTo(HaveOccurred())

		Expect(flags).To(Equal([]string{
			`-buildmode=default`,
			`-ldflags="-X main.variable=some-value"`,
			`-ldflags='-X main.variable=some-other-value'`,
			`-tags=paketo`,
		}))
	})

	context("failure cases", func() {
		context("when there is an unclosed \"", func() {
			it("returns an error", func() {
				_, err := flagParser.ParseFlags(`-buildmode=default -ldflags="-X main.variable=some-value -ldflags='-X main.variable=some-other-value' -tags=paketo`)
				Expect(err).To(MatchError("expected closing \" before end of input"))
			})
		})

		context("when there is an unclosed '", func() {
			it("returns an error", func() {
				_, err := flagParser.ParseFlags(`-buildmode=default -ldflags="-X main.variable=some-value" -ldflags='-X main.variable=some-other-value -tags=paketo`)
				Expect(err).To(MatchError("expected closing ' before end of input"))
			})
		})
	})

}
