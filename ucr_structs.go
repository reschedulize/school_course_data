package school_course_data

type ClassSearchResult struct {
	Success    bool                    `json:"success"`
	TotalCount int                     `json:"totalCount"`
	Data       []ClassSearchResultData `json:"data"`
}

type ClassSearchResultData struct {
	//ID                      int         `json:"id"`
	//Term                    string      `json:"term"`
	//TermDesc                string      `json:"termDesc"`
	CourseReferenceNumber string `json:"courseReferenceNumber"`
	//PartOfTerm              string      `json:"partOfTerm"`
	//CourseNumber            string      `json:"courseNumber"`
	//Subject                 string      `json:"subject"`
	//SubjectDescription      string      `json:"subjectDescription"`
	//SequenceNumber          string      `json:"sequenceNumber"`
	//CampusDescription       string      `json:"campusDescription"`
	ScheduleTypeDescription string `json:"scheduleTypeDescription"`
	//CourseTitle             string      `json:"courseTitle"`
	CreditHours       int `json:"creditHours"`
	MaximumEnrollment int `json:"maximumEnrollment"`
	Enrollment        int `json:"enrollment"`
	SeatsAvailable    int `json:"seatsAvailable"`
	WaitCapacity      int `json:"waitCapacity"`
	WaitCount         int `json:"waitCount"`
	WaitAvailable     int `json:"waitAvailable"`
	//CrossList               interface{} `json:"crossList"`
	//CrossListCapacity       interface{} `json:"crossListCapacity"`
	//CrossListCount          interface{} `json:"crossListCount"`
	//CrossListAvailable      interface{} `json:"crossListAvailable"`
	//CreditHourHigh          int         `json:"creditHourHigh"`
	//CreditHourLow           int         `json:"creditHourLow"`
	//CreditHourIndicator     string      `json:"creditHourIndicator"`
	//OpenSection             bool        `json:"openSection"`
	LinkIdentifier  string `json:"linkIdentifier"`
	IsSectionLinked bool   `json:"isSectionLinked"`
	SubjectCourse   string `json:"subjectCourse"`
	//Faculty                 []struct {
	//	BannerID              string      `json:"bannerId"`
	//	Category              interface{} `json:"category"`
	//	Class                 string      `json:"class"`
	//	CourseReferenceNumber string      `json:"courseReferenceNumber"`
	//	DisplayName           string      `json:"displayName"`
	//	EmailAddress          string      `json:"emailAddress"`
	//	PrimaryIndicator      bool        `json:"primaryIndicator"`
	//	Term                  string      `json:"term"`
	//} `json:"faculty"`
	MeetingsFaculty     []MeetingsFaculty `json:"meetingsFaculty"`
	ReservedSeatSummary *struct {
		//Class                       string `json:"class"`
		//CourseReferenceNumber       string `json:"courseReferenceNumber"`
		MaximumEnrollmentReserved   int `json:"maximumEnrollmentReserved"`
		MaximumEnrollmentUnreserved int `json:"maximumEnrollmentUnreserved"`
		SeatsAvailableReserved      int `json:"seatsAvailableReserved"`
		SeatsAvailableUnreserved    int `json:"seatsAvailableUnreserved"`
		//TermCode                    string `json:"termCode"`
		WaitAvailableReserved   int `json:"waitAvailableReserved"`
		WaitAvailableUnreserved int `json:"waitAvailableUnreserved"`
		WaitCapacityReserved    int `json:"waitCapacityReserved"`
		WaitCapacityUnreserved  int `json:"waitCapacityUnreserved"`
	} `json:"reservedSeatSummary"`
}

type MeetingsFaculty struct {
	//Category              string        `json:"category"`
	//Class                 string        `json:"class"`
	//CourseReferenceNumber string        `json:"courseReferenceNumber"`
	//Faculty               []interface{} `json:"faculty"`
	MeetingTime struct {
		BeginTime string `json:"beginTime"`
		//Building              string  `json:"building"`
		//BuildingDescription   string  `json:"buildingDescription"`
		//Campus                string  `json:"campus"`
		//CampusDescription     string  `json:"campusDescription"`
		//Category              string  `json:"category"`
		//Class                 string  `json:"class"`
		//CourseReferenceNumber string  `json:"courseReferenceNumber"`
		//CreditHourSession     float64 `json:"creditHourSession"`
		//EndDate               string  `json:"endDate"`
		EndTime string `json:"endTime"`
		Friday  bool   `json:"friday"`
		//HoursWeek             float64 `json:"hoursWeek"`
		//MeetingScheduleType   string  `json:"meetingScheduleType"`
		Monday bool `json:"monday"`
		//Room                  string  `json:"room"`
		Saturday bool `json:"saturday"`
		//StartDate             string  `json:"startDate"`
		Sunday bool `json:"sunday"`
		//Term                  string  `json:"term"`
		Thursday  bool `json:"thursday"`
		Tuesday   bool `json:"tuesday"`
		Wednesday bool `json:"wednesday"`
	} `json:"meetingTime"`
	Term string `json:"term"`
}
