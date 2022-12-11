package objects

type Student struct {
	id            int
	accID         int
	name          string
	surname       string
	studentGroup  string
	studentNumber string
	roomID        int
}

type StudentDTO struct {
	name          string
	surname       string
	studentGroup  string
	studentNumber string
}

type StudentResponseDTO struct {
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	StudentGroup  string `json:"studentGroup"`
	StudentNumber string `json:"studentNumber"`
	RoomID        int    `json:"roomID"`
}

func NewStudentWithParams(id, accID int, name, surname, studentGroup, studentNumber string, roomID int) Student {
	return Student{
		id:            id,
		accID:         accID,
		name:          name,
		surname:       surname,
		studentGroup:  studentGroup,
		studentNumber: studentNumber,
		roomID:        roomID,
	}
}

func NewEmptyStudent() Student {
	return Student{id: None}
}

func (s *Student) GetID() int {
	return s.id
}

func (s *Student) GetAccID() int {
	return s.accID
}

func (s *Student) GetName() string {
	return s.name
}

func (s *Student) GetSurname() string {
	return s.surname
}

func (s *Student) GetStudentGroup() string {
	return s.studentGroup
}

func (s *Student) GetStudentNumber() string {
	return s.studentNumber
}

func (s *Student) GetRoomID() int {
	return s.roomID
}

func (s *Student) SetRoomID(id int) {
	s.roomID = id
}

func (s *Student) SetGroup(group string) {
	s.studentGroup = group
}

func (s *Student) SetName(name string) {
	s.name = name
}

func (s *Student) SetSurname(surname string) {
	s.surname = surname
}

func NewStudentDTO(name, surname, group, studNumber string) StudentDTO {
	return StudentDTO{
		name:          name,
		surname:       surname,
		studentGroup:  group,
		studentNumber: studNumber,
	}
}

func (s *StudentDTO) GetName() string {
	return s.name
}

func (s *StudentDTO) GetSurname() string {
	return s.surname
}

func (s *StudentDTO) GetStudentGroup() string {
	return s.studentGroup
}

func (s *StudentDTO) GetStudentNumber() string {
	return s.studentNumber
}

func CreateStudentResponse(students []Student) []StudentResponseDTO {
	newArray := make([]StudentResponseDTO, Empty)
	for _, tmpStudent := range students {
		newArray = append(newArray, StudentResponseDTO{
			Name:          tmpStudent.GetName(),
			Surname:       tmpStudent.GetSurname(),
			StudentGroup:  tmpStudent.GetStudentGroup(),
			StudentNumber: tmpStudent.GetStudentNumber(),
			RoomID:        tmpStudent.GetRoomID(),
		})
	}
	return newArray
}

func CreateStudentResponseSingle(student Student) StudentResponseDTO {
	return StudentResponseDTO{
		Name:          student.GetName(),
		Surname:       student.GetSurname(),
		StudentGroup:  student.GetStudentGroup(),
		StudentNumber: student.GetStudentNumber(),
		RoomID:        student.GetRoomID(),
	}
}
