package school_course_data

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/go-redis/redis"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func convertToMask(class *Class, day Day, begin uint16, end uint16) {
	class.WeekMask |= 1 << day

	class.DayMasks[day][0] |= createMask(0, 63, int64(begin), int64(end))
	class.DayMasks[day][1] |= createMask(64, 127, int64(begin), int64(end))
	class.DayMasks[day][2] |= createMask(128, 191, int64(begin), int64(end))
	class.DayMasks[day][3] |= createMask(192, 255, int64(begin), int64(end))
}

type UCRAPI struct {
	BaseAPI
}

func NewUCRAPI(redis *redis.Client) *UCRAPI {
	return &UCRAPI{BaseAPI: BaseAPI{Redis: redis}}
}

func (a *UCRAPI) Terms(maxResults int64) ([]string, error) {
	request, err := http.NewRequest("GET", "https://registrationssb.ucr.edu/StudentRegistrationSsb/ssb/classSearch/getTerms", nil)

	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Set("offset", "1")
	query.Set("max", strconv.FormatInt(maxResults, 10))
	request.URL.RawQuery = query.Encode()

	response, err := a.makeRequest(request, 12*time.Hour)

	if err != nil {
		return nil, err
	}

	var cap uint64

	if maxResults >= 0 {
		cap = uint64(maxResults)
	} else {
		cap = 0
	}

	terms := make([]string, 0, cap)

	_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, parserErr error) {
		if parserErr != nil {
			err = ParseError
			return
		}

		if dataType != jsonparser.Object {
			err = ParseError
			return
		}

		code, parserErr := jsonparser.GetString(value, "code")

		if parserErr != nil {
			err = ParseError
			return
		}

		terms = append(terms, code)
	})

	if err != nil {
		return nil, err
	}

	return terms, nil
}

// Returns courses for given term in ascending sorted order
func (a *UCRAPI) Courses(term string, maxResults int64) ([]string, error) {
	request, err := http.NewRequest("GET", "https://registrationssb.ucr.edu/StudentRegistrationSsb/ssb/classSearch/get_subjectcoursecombo", nil)

	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Set("term", term)
	query.Set("searchTerm", "")
	query.Set("offset", "1")
	query.Set("max", strconv.FormatInt(maxResults, 10))
	request.URL.RawQuery = query.Encode()

	response, err := a.makeRequest(request, 1*time.Hour)

	if err != nil {
		return nil, err
	}

	var cap uint64

	if maxResults >= 0 {
		cap = uint64(maxResults)
	} else {
		cap = 0
	}

	courses := make([]string, 0, cap)

	_, err = jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, parserErr error) {
		if parserErr != nil {
			parserErr = ParseError
			return
		}

		if dataType != jsonparser.Object {
			err = ParseError
			return
		}

		code, parserErr := jsonparser.GetString(value, "code")

		if parserErr != nil {
			err = ParseError
			return
		}

		courses = append(courses, code)
	})

	if err != nil {
		return nil, ParseError
	}

	// TODO: Cache sorted slice instead of sorting every time
	sort.Strings(courses)

	return courses, nil
}

func (a *UCRAPI) Classes(term string, course string, maxResults int64) ([]*Class, error) {
	classes, err := a.rawClasses(term, course, maxResults)

	if err != nil {
		return nil, err
	}

	output := make([]*Class, len(classes))

	for i, c := range classes {
		class := Class{
			CourseName: c.SubjectCourse,
			CRN:        c.CourseReferenceNumber,
			Type:       c.ScheduleTypeDescription,
			Units:      c.CreditHours,
			Schedule:   make([][]*MeetingTime, 7),
		}

		if c.IsSectionLinked {
			class.LinkID = c.LinkIdentifier
		} else {
			class.LinkID = "X0"
		}

		if c.ReservedSeatSummary == nil {
			class.Seats = &Availability{
				AvailableUnreserved: c.SeatsAvailable,
				CapacityUnreserved:  c.MaximumEnrollment,
				AvailableReserved:   0,
				CapacityReserved:    0,
			}

			class.Waitlist = &Availability{
				AvailableUnreserved: c.WaitAvailable,
				CapacityUnreserved:  c.WaitCapacity,
				AvailableReserved:   0,
				CapacityReserved:    0,
			}
		} else {
			class.Seats = &Availability{
				AvailableUnreserved: c.ReservedSeatSummary.SeatsAvailableUnreserved,
				CapacityUnreserved:  c.ReservedSeatSummary.MaximumEnrollmentUnreserved,
				AvailableReserved:   c.ReservedSeatSummary.SeatsAvailableReserved,
				CapacityReserved:    c.ReservedSeatSummary.MaximumEnrollmentReserved,
			}

			class.Waitlist = &Availability{
				AvailableUnreserved: c.ReservedSeatSummary.WaitAvailableUnreserved,
				CapacityUnreserved:  c.ReservedSeatSummary.WaitCapacityUnreserved,
				AvailableReserved:   c.ReservedSeatSummary.WaitAvailableReserved,
				CapacityReserved:    c.ReservedSeatSummary.WaitCapacityReserved,
			}
		}

		for _, mf := range c.MeetingsFaculty {
			var begin uint16
			var end uint16

			var err error

			begin, err = parseTime(mf.MeetingTime.BeginTime)
			end, err = parseTime(mf.MeetingTime.EndTime)

			if err == nil {
				if mf.MeetingTime.Sunday {
					convertToMask(&class, Sunday, begin, end)
				}

				if mf.MeetingTime.Monday {
					convertToMask(&class, Monday, begin, end)
				}

				if mf.MeetingTime.Tuesday {
					convertToMask(&class, Tuesday, begin, end)
				}

				if mf.MeetingTime.Wednesday {
					convertToMask(&class, Wednesday, begin, end)
				}

				if mf.MeetingTime.Thursday {
					convertToMask(&class, Thursday, begin, end)
				}

				if mf.MeetingTime.Friday {
					convertToMask(&class, Friday, begin, end)
				}

				if mf.MeetingTime.Saturday {
					convertToMask(&class, Saturday, begin, end)
				}

				// OLD
				if mf.MeetingTime.Sunday {
					class.Schedule[Sunday] = append(class.Schedule[Sunday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Monday {
					class.Schedule[Monday] = append(class.Schedule[Monday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Tuesday {
					class.Schedule[Tuesday] = append(class.Schedule[Tuesday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Wednesday {
					class.Schedule[Wednesday] = append(class.Schedule[Wednesday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Thursday {
					class.Schedule[Thursday] = append(class.Schedule[Thursday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Friday {
					class.Schedule[Friday] = append(class.Schedule[Friday], &MeetingTime{Begin: begin, End: end})
				}

				if mf.MeetingTime.Saturday {
					class.Schedule[Saturday] = append(class.Schedule[Saturday], &MeetingTime{Begin: begin, End: end})
				}
			}
		}

		output[i] = &class
	}

	return output, nil
}

func (a *UCRAPI) rawClasses(term string, course string, max int64) (courses []ClassSearchResultData, err error) {
	secondRequest, err := http.NewRequest("GET", "https://registrationssb.ucr.edu/StudentRegistrationSsb/ssb/searchResults/searchResults", nil)

	if err != nil {
		return
	}

	query := secondRequest.URL.Query()
	query.Add("txt_term", term)
	query.Add("txt_subjectcoursecombo", course)
	query.Add("pageOffset", "0")
	query.Add("pageMaxSize", strconv.FormatInt(max, 10))
	secondRequest.URL.RawQuery = query.Encode()

	response, ok := a.getFromCache(secondRequest)

	if !ok {
		jar, err := cookiejar.New(nil)

		if err != nil {
			return nil, err
		}

		// First request
		form := make(url.Values)
		form.Add("term", term)

		firstRequest, err := http.NewRequest("POST", "https://registrationssb.ucr.edu/StudentRegistrationSsb/ssb/term/search?mode=search", strings.NewReader(form.Encode()))

		if err != nil {
			return nil, err
		}

		firstRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{Jar: jar}

		res, err := client.Do(firstRequest)
		defer res.Body.Close()

		if err != nil {
			return nil, err
		}

		// Second request
		res, err = client.Do(secondRequest)
		defer res.Body.Close()

		if err != nil {
			return nil, err
		}

		raw, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		response, err = trimJSON(raw)

		if err != nil {
			return nil, err
		}

		// Check that data is good before caching
		success, err := jsonparser.GetBoolean(response, "success")

		if err != nil {
			return nil, ParseError
		}

		if !success {
			return nil, ServerError
		}

		count, err := jsonparser.GetInt(response, "totalCount")

		if err != nil {
			return nil, ParseError
		}

		if count == 0 {
			return nil, CourseNotFound
		}

		// Cache response
		go a.cacheRequest(secondRequest, response, 1*time.Hour)
	}

	// TODO: Use jsonparser instead of builtin json unmarshal
	var result ClassSearchResult
	err = json.Unmarshal(response, &result)

	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
