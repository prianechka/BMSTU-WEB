package server

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"src/configs/backend"
	"src/db/roomRepo"
	"src/db/studentRepo"
	"src/db/thingRepo"
	"src/db/userRepo"
	"src/delivery/http/authHandler"
	"src/delivery/http/roomHandler"
	"src/delivery/http/studentHandler"
	"src/delivery/http/thingHandler"
	"src/docs"
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/controllers/userController"
	"src/logic/managers/appManager"
	"src/logic/managers/authManager"
	"src/logic/managers/roomManager"
	"src/logic/managers/studentManager"
	"src/logic/managers/thingManager"
	"src/middleware"
	utils "src/utils/connection"
)

var serverType = os.Getenv("SERVER_TYPE")

type Server struct {
	config *configs.ServerConfig
	logger *logrus.Entry
}

func CreateServer(config *configs.ServerConfig, logger *logrus.Entry) *Server {
	return &Server{config: config, logger: logger}
}

// Start
// @securityDefinitions.apikey JWT-Token
func (s *Server) Start() error {
	if serverType == "mirror" {
		docs.SwaggerInfo.BasePath = "/mirror1"
	}

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

	RoomManager := roomManager.CreateNewRoomManager(RoomController)
	StudentManager := studentManager.CreateNewStudentManager(RoomController, StudentController, UserController, ThingController)
	ThingManager := thingManager.CreateNewThingManager(RoomController, StudentController, ThingController)
	AuthManager := authManager.CreateNewAuthManager(UserController)
	AppManager := appManager.AppManager{}

	StudentHandler := studentHandler.CreateNewStudentHandler(s.logger, *StudentManager)
	AuthHandler := authHandler.CreateNewAuthHandler(s.logger, *AuthManager, AppManager)
	ThingHandler := thingHandler.CreateNewThingHandler(s.logger, *ThingManager)
	RoomHandler := roomHandler.CreateNewRoomHandler(s.logger, *RoomManager)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/students", StudentHandler.GetAllStudents).Methods("GET")
	router.HandleFunc("/students", StudentHandler.AddNewStudent).Methods("POST")
	router.HandleFunc("/students/{stud-number}", StudentHandler.ChangeStudentGroup).Methods("PUT")
	router.HandleFunc("/students/{stud-number}", StudentHandler.ViewStudentInfo).Methods("GET")
	router.HandleFunc("/student-live-acts/{stud-number}", StudentHandler.TransferStudent).Methods("POST")
	router.HandleFunc("/student-live-acts/{stud-number}", StudentHandler.ViewStudentLivingHistory).Methods("GET")
	router.HandleFunc("/student-things-acts/{mark-number}", StudentHandler.TransferThingFromToStudents).Methods("POST")
	router.HandleFunc("/student-things-acts/{mark-number}", ThingHandler.ViewThingHistory).Methods("GET")

	router.HandleFunc("/things", ThingHandler.GetThings).Methods("GET")
	router.HandleFunc("/things/{mark-number}", ThingHandler.GetThing).Methods("GET")
	router.HandleFunc("/things", ThingHandler.AddNewThing).Methods("POST")
	router.HandleFunc("/things/{mark-number}", ThingHandler.TransferThingBetweenRooms).Methods("PATCH")
	router.HandleFunc("/login", AuthHandler.Authorize).Methods("POST")
	router.HandleFunc("/rooms", RoomHandler.GetAllRooms).Methods("GET")
	router.HandleFunc("/rooms/{room-id}", RoomHandler.GetRoom).Methods("GET")

	accessRouter := middleware.CheckAccess(router)
	upgradedRouter := middleware.Panic(accessRouter)

	return http.ListenAndServe(s.config.PortToStart, upgradedRouter)
}
