package demo_test

import (
	"encoding/json"
	"fmt"

	"github.com/goinbox/gohttp/httpserver"
	"github.com/goinbox/gohttp/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gdemo/controller/api"
	"gdemo/controller/api/demo"
	"gdemo/logic/factory"
	"gdemo/perror"
	"gdemo/test"
)

var _ = Describe("Demo API", Ordered, func() {
	var runner *test.ApiControllerRunner
	var id int64

	BeforeEach(func() {
		r := router.NewRouter()
		r.MapRouteItems(new(demo.Controller))

		runner = &test.ApiControllerRunner{
			Server: httpserver.NewServer(r),
			App:    factory.DefaultLogicFactory.AppLogic().ListAllApps(test.Context())[0],
		}
	})

	Context("add", func() {
		It("success", func() {
			content, err := runner.Run("/Demo/Add", `
{
  "Name": "demo-suite-test"
}
`)
			Expect(err).To(BeNil())
			var resp struct {
				api.BaseResponse
				Data *demo.AddResponse
			}
			_ = json.Unmarshal(content, &resp)
			Expect(resp.Errno).To(Equal(perror.Success))
			id = resp.Data.ID
		})
	})

	Context("edit", func() {
		It("success", func() {
			content, err := runner.Run("/Demo/Edit", fmt.Sprintf(`
{
  "ID": %d,
  "Status": 0
}
`, id))
			Expect(err).To(BeNil())
			var resp struct {
				api.BaseResponse
				Data *demo.EditResponse
			}
			_ = json.Unmarshal(content, &resp)
			Expect(resp.Errno).To(Equal(perror.Success))
			Expect(resp.Data.RowsAffected).To(BeEquivalentTo(1))
		})
	})

	Context("index", func() {
		It("success", func() {
			content, err := runner.Run("/Demo/Index", fmt.Sprintf(`
{
  "IDs": [%d]
}
`, id))
			Expect(err).To(BeNil())
			var resp struct {
				api.BaseResponse
				Data *demo.IndexResponse
			}
			_ = json.Unmarshal(content, &resp)
			Expect(resp.Errno).To(Equal(perror.Success))
			Expect(resp.Data.Total).To(BeEquivalentTo(1))
		})
	})

	Context("del", func() {
		It("success", func() {
			content, err := runner.Run("/Demo/Del", fmt.Sprintf(`
{
  "IDs": [%d]
}
`, id))
			Expect(err).To(BeNil())
			var resp struct {
				api.BaseResponse
				Data *demo.DelResponse
			}
			_ = json.Unmarshal(content, &resp)
			Expect(resp.Errno).To(Equal(perror.Success))
			Expect(resp.Data.RowsAffected).To(BeEquivalentTo(1))
		})
	})
})
