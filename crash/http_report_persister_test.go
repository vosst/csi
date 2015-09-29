package crash

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
)

func TestHttpCrashReportPersisterChecksForReachability(t *testing.T) {
	u, _ := url.Parse("http://localhost:9090")
	mrm := &MockReachabilityMonitor{}
	mrm.On("CheckHostReachability", "localhost:9090").Return(NotReachable)

	persister := HttpCrashReportPersister{*u, mrm, &http.Client{}}

	f, _ := os.Open("test_data/test.crash")

	report, err := ParseCrashReport(NewLineReader{f})
	assert.Nil(t, err)

	err = persister.Persist(report)
	assert.NotNil(t, err)

	mrm.AssertExpectations(t)
}

func TestHttpCrashReportPersisterDoesNotSendViaWWAN(t *testing.T) {
	u, _ := url.Parse("http://localhost:9090")
	mrm := &MockReachabilityMonitor{}
	mrm.On("CheckHostReachability", "localhost:9090").Return(IsReachable | IsWWAN)

	persister := HttpCrashReportPersister{*u, mrm, &http.Client{}}

	f, _ := os.Open("test_data/test.crash")

	report, err := ParseCrashReport(NewLineReader{f})
	assert.Nil(t, err)

	err = persister.Persist(report)
	assert.NotNil(t, err)

	mrm.AssertExpectations(t)
}

func TestHttpCrashReportPersisterSendsValidBSON(t *testing.T) {
	go func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			bytes, _ := ioutil.ReadAll(r.Body)

			report := make(map[string]interface{})
			if err := bson.Unmarshal(bytes, report); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Write([]byte(`EMPTY`))
			}
		}))

		server.Close()
	}()

	u, _ := url.Parse("http://localhost:9090")
	mrm := &MockReachabilityMonitor{}
	mrm.On("CheckHostReachability", "localhost:9090").Return(IsReachable)

	persister := HttpCrashReportPersister{*u, mrm, &http.Client{}}

	f, _ := os.Open("test_data/test.crash")

	report, err := ParseCrashReport(NewLineReader{f})
	assert.Nil(t, err)

	err = persister.Persist(report)
	assert.Nil(t, err)
}
