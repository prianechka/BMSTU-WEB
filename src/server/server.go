package server

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"src/configs"
	"src/db/roomRepo"
	"src/db/studentRepo"
	"src/db/thingRepo"
	"src/db/userRepo"
	"src/delivery/http/authHandler"
	"src/delivery/http/studentHandler"
	"src/delivery/http/thingHandler"
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/controllers/userController"
	"src/logic/managers/appManager"
	"src/logic/managers/authManager"
	"src/logic/managers/studentManager"
	"src/logic/managers/thingManager"
	utils "src/utils/connection"
)

type Server struct {
	config *configs.ServerConfig
	logger *logrus.Entry
}

func CreateServer(config *configs.ServerConfig, logger *logrus.Entry) *Server {
	return &Server{config: config, logger: logger}
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	router := r.PathPrefix("/api/v1/").Subrouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	roomDB := utils.NewPgSQLConnection(s.config.ConnParams)
	studentDB := utils.NewPgSQLConnection(s.config.ConnParams)
	thingDB := utils.NewPgSQLConnection(s.config.ConnParams)
	userDB := utils.NewPgSQLConnection(s.config.ConnParams)

	roomRepository := roomRepo.PgRoomRepo{Conn: roomDB}
	studentRepository := studentRepo.PgStudentRepo{Conn: studentDB}
	thingRepository := thingRepo.PgThingRepo{Conn: thingDB}
	userRepository := userRepo.PgUserRepo{Conn: userDB}

	RoomController := roomController.RoomController{Repo: &roomRepository}
	StudentController := studentController.StudentController{Repo: &studentRepository}
	ThingController := thingController.ThingController{Repo: &thingRepository}
	UserController := userController.UserController{Repo: &userRepository}

	StudentManager := studentManager.CreateNewStudentManager(RoomController, StudentController, UserController, ThingController)
	ThingManager := thingManager.CreateNewThingManager(RoomController, StudentController, ThingController)
	AuthManager := authManager.CreateNewAuthManager(UserController)
	AppManager := appManager.AppManager{}

	StudentHandler := studentHandler.CreateNewStudentHandler(s.logger, *StudentManager)
	AuthHandler := authHandler.CreateNewAuthHandler(s.logger, *AuthManager, AppManager)
	ThingHandler := thingHandler.CreateNewThingHandler(s.logger, *ThingManager)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/students/", StudentHandler.GetAllStudents).Methods("GET")
	router.HandleFunc("/students/", StudentHandler.AddNewStudent).Methods("POST")
	router.HandleFunc("/students/{stud-number}/", StudentHandler.ChangeStudentGroup).Methods("PUT")
	router.HandleFunc("/students/{stud-number}/rooms/", StudentHandler.SettleStudent).Methods("POST")
	router.HandleFunc("/students/{stud-number}/rooms/", StudentHandler.EvicStudent).Methods("DELETE")

	router.HandleFunc("/students/{stud-number}/things/", ThingHandler.ViewStudentThings).Methods("GET")
	router.HandleFunc("/students/{stud-number}/things/{mark-number}/", StudentHandler.GiveStudentThing).Methods("POST")
	router.HandleFunc("/students/{stud-number}/things/{mark-number}/", StudentHandler.ReturnThingFromStudent).Methods("DELETE")

	router.HandleFunc("/things/", ThingHandler.ViewAllThings).Methods("GET")
	router.HandleFunc("/things/", ThingHandler.AddNewThing).Methods("POST")
	router.HandleFunc("/things/free/", ThingHandler.ViewFreeThings).Methods("GET")
	router.HandleFunc("/things/{mark-number}/", ThingHandler.TransferThingBetweenRooms).Methods("PATCH")
	router.HandleFunc("/auth", AuthHandler.Authorize).Methods("POST")

	return http.ListenAndServe(s.config.PortToStart, router)
}
