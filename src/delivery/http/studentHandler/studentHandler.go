package studentHandler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"src/logic/managers/studentManager"
)

type StudentHandler struct {
	logger *logrus.Entry
	studentManager.StudentManager
}

func (sh *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {}
