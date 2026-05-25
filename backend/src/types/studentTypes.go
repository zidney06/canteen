package types

type GetStudent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsBlocked bool   `json:"is_blocked"`
}

type StudentData struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type GetStudentByIdType struct {
	Message string      `json:"message"`
	Data    StudentData `json:"data"`
}
