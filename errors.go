package school_course_data

import "errors"

var ParseError = errors.New("error parsing JSON response")
var ServerError = errors.New("server error")
var CourseNotFound = errors.New("course not found")
var TimeParseError = errors.New("unable to parse time")
