package wrappa_test

import (
	"github.com/concourse/concourse/atc/wrappa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concurrent Request Limiting", func() {
	Describe("ConcurrentRequestLimitFlag#UnmarshalFlag", func() {
		var crl wrappa.ConcurrentRequestLimitFlag

		BeforeEach(func() {
			crl = wrappa.ConcurrentRequestLimitFlag{}
		})

		It("parses the API action and limit", func() {
			crl.UnmarshalFlag("ListAllJobs=3")
			Expect(crl.Action).To(Equal("ListAllJobs"), "wrong action")
			Expect(crl.Limit).To(Equal(3), "wrong limit")
		})

		It("returns an error when the flag has no equals sign", func() {
			err := crl.UnmarshalFlag("banana")
			Expect(err).To(
				MatchError(
					"invalid concurrent request limit " +
					"'banana': value must be an assignment",
				),
			)
		})

		It("returns an error when the flag has multiple equals signs", func() {
			err := crl.UnmarshalFlag("foo=bar=baz")
			Expect(err).To(
				MatchError(
					"invalid concurrent request limit " +
					"'foo=bar=baz': value must be an assignment",
				),
			)
		})

		It("returns an error when the limit is not an integer", func() {
			err := crl.UnmarshalFlag("ListAllJobs=foo")
			Expect(err).To(
				MatchError(
					"invalid concurrent request limit " +
					"'ListAllJobs=foo': limit must be an integer",
				),
			)
		})
	})
})
