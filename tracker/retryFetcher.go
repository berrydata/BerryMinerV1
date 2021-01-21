package tracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/berry-data/BerryMiner/util"
)

var retryFetchLog = util.NewLogger("tracker", "FetchWithRetries")

//FetchRequest holds info for a request
type FetchRequest struct {
	queryURL string
	timeout  time.Duration
}

func fetchWithRetries(req *FetchRequest) ([]byte, error) {
	return _recFetch(req, time.Now().Add(req.timeout))
}

func batchFetchWithRetries(reqs []*FetchRequest) ([][]byte, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	res := make([][]byte, len(reqs))

	//A potential optimization is to use go routines for each sub-API request
	for i := 0; i < len(reqs); i++ {
		req := reqs[i]
		data, err := _recFetch(req, time.Now().Add(req.timeout))

		//in this case, one failure means all fail
		if err != nil {
			retryFetchLog.Warn("Batch request failure, ignoring that part of the request: %v\n", err)
			res[i] = nil
		} else {
			res[i] = data
		}
	}
	return res, nil
}

func _recFetch(req *FetchRequest, expiration time.Time) ([]byte, error) {
	retryFetchLog.Debug("Fetch request will expire at: %v (timeout: %v)", expiration, req.timeout)

	now := time.Now()
	client := http.Client{
		Timeout: expiration.Sub(now),
	}

	r, err := client.Get(req.queryURL)
	if err != nil {
		//log local non-timeout errors for now
		retryFetchLog.Warn("Problem fetching data from: %s. %v", req.queryURL, err)
		now := time.Now()
		if now.After(expiration) {
			retryFetchLog.Error("Timeout expired, not retrying query and passing error up")
			return nil, err
		}
		//FIXME: should this be configured as fetch error sleep duration?
		time.Sleep(500 * time.Millisecond)

		//try again
		retryFetchLog.Warn("Trying fetch again...")
		return _recFetch(req, expiration)
	}

	data, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		retryFetchLog.Warn("Response from fetching  %s. Response code: %d, payload: %s", req.queryURL, r.StatusCode, data)
		//log local non-timeout errors for now
		now := time.Now()
		if now.After(expiration) {
			return nil, fmt.Errorf("giving up fetch request after request timeout: %d", r.StatusCode)
		}
		//FIXME: should this be configured as fetch error sleep duration?
		time.Sleep(500 * time.Millisecond)

		//try again
		return _recFetch(req, expiration)
	}
	return data, nil
}
