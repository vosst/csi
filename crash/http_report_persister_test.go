package crash

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vosst/csi/machine"
	"labix.org/v2/mgo/bson"
)

func TestHttpCrashReportPersisterSendsValidBSON(t *testing.T) {
	go func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			bytes, _ := ioutil.ReadAll(r.Body)

			report := make(map[string]interface{})
			if err := bson.Unmarshal(bytes, report); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}))

		server.Close()
	}()

	mi := &machine.MockIdentifier{}
	mi.On("Identify").Return([]byte{42, 42, 42}, nil)

	u, _ := url.Parse("http://localhost:9090")

	persister := HttpReportPersister{*u, mi, &http.Client{}}

	f, _ := os.Open("test_data/test.crash")

	report, err := ParseReport(NewLineReader{f})
	assert.Nil(t, err)

	err = persister.Persist(report)
	assert.Nil(t, err)
}
