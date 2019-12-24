package school_course_data

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"net/http"
	"time"
)

type IAPI interface {
	Terms(int64) ([]string, error)
	Courses(string, int64) ([]string, error)
	Classes(string, string, int64) ([]*Class, error)
}

type BaseAPI struct {
	Redis *redis.Client
}

func (a *BaseAPI) cacheRequest(request *http.Request, data []byte, ttl time.Duration) {
	if a.Redis != nil {
		a.Redis.Set(request.URL.String(), data, ttl)
	}
}

func (a *BaseAPI) getFromCache(request *http.Request) (data []byte, ok bool) {
	if a.Redis == nil {
		return nil, false
	}

	data, err := a.Redis.Get(request.URL.String()).Bytes()

	if err != nil {
		return nil, false
	}

	return data, true
}

func (a *BaseAPI) makeRequest(request *http.Request, ttl time.Duration) ([]byte, error) {
	data, ok := a.getFromCache(request)

	// Exists in cache
	if ok {
		return data, nil
	}

	// Does not exist in cache
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	trimmed, err := trimJSON(raw)

	if err != nil {
		return nil, err
	}

	a.cacheRequest(request, trimmed, ttl)

	return trimmed, nil
}
