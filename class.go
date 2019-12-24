package school_course_data

type Day uint

const (
	Sunday    Day = 0
	Monday    Day = 1
	Tuesday   Day = 2
	Wednesday Day = 3
	Thursday  Day = 4
	Friday    Day = 5
	Saturday  Day = 6
)

type Class struct {
	CourseName string        `json:"course"`
	CRN        string        `json:"crn"`
	Type       string        `json:"type"`
	LinkID     string        `json:"link_id"`
	Units      int           `json:"units"`
	Seats      *Availability `json:"seats"`
	Waitlist   *Availability `json:"waitlist"`
	WeekMask   uint8         `json:"week_mask"`
	DayMasks   [7][4]uint64  `json:"day_masks"`
	Schedule [][]*MeetingTime `json:"schedule"`
}

type Availability struct {
	AvailableUnreserved int `json:"available_unreserved"`
	CapacityUnreserved  int `json:"capacity_unreserved"`

	AvailableReserved int `json:"available_reserved"`
	CapacityReserved  int `json:"capacity_reserved"`
}

type MeetingTime struct {
	Begin uint16 `json:"begin"`
	End   uint16 `json:"end"`
}
