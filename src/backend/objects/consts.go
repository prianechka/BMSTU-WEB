package objects

const (
	Null                           = 0
	None                           = -1
	NotLiving                      = 0
	Ret                            = 0
	Get                            = 1
	EmptyString                    = ""
	Empty                          = 0
	DefaultPageSize                = 25
	DefaultPage                    = 0
	Max                            = 10000
	LoginInputMessage              = "Введите логин: "
	PasswordInputMessage           = "Введите пароль: "
	MarkInputMessage               = "Введите маркировочный номер: "
	TypeInputMessage               = "Введите тип вещи: "
	NameInputMessage               = "Введите имя студента: "
	SurnameInputMessage            = "Введите фамилию студента: "
	GroupInputMessage              = "Введите номер группы студента: "
	NumberInputMessage             = "Введите номер студенческого билета: "
	IDRoomInputMessage             = "Введите ID комнаты студента: "
	RoomIDInputMessage             = "Введите ID комнаты: "
	MarkNumInputMessage            = "Введите маркировочный номер вещи: "
	AuthOK                         = "Авторизация прошла успешно!"
	TransferThingOK                = "Вещь успешно перемещена!"
	EvicStudentOK                  = "Студент успешно выселен!"
	SettleStudentOK                = "Студент успешно заселён!"
	GiveThingOK                    = "Вещь передана!"
	ReturnThingOK                  = "Вещь отдана обратно!"
	AddOK                          = "Операция успешно проведена!"
	StudentNotFoundErrorString     = "Студент не найден!"
	ReturnThingErrorString         = "Вещь и так была не у студента!"
	ThingNotFound                  = "Вещь не найдена!"
	GiveThingErrorString           = "Вещь уже у другого студента!"
	SettleStudentErrorString       = "Студент уже живёт в другой комнате!"
	EvicStudentErrorString         = "Студент уже нигде не живёт!"
	UserNotFoundErrorString        = "Пользователь не найден!"
	UserAddErrorString             = "Пользователя не удалось добавить!"
	StudentAlreadyExistErrorString = "Студент с таким же студенческим билетом уже существует!"
	UserAlreadyExistErrorString    = "Пользователь с таким логином уже существует!"
	StudentAddErrorString          = "Студента не удалось добавить!"
	StudentChangeErrorString       = "Не удалось изменить данные студента!"
	StudentChangeOKString          = "Данные о студенте успешно обновлены!"
	LoginErrorString               = "Такого логина не существует!"
	PasswordErrorString            = "Пароль введен неверно!"
	RoomNotFoundErrorString        = "Комната не найдена!"
	ThingNotInRoomErrorString      = "Вещь не находится в этой комнате!"
	ThingAlreadyInRoomErrorString  = "Вещь уже находится в этой комнате!"
	DeleteRoomErrorString          = "Не удалось удалить комнату!"
	DeleteThingErrorString         = "Не удалось удалить вещь!"
	UniqueMarkNumberErrorString    = "Вещь с таким же уникальным номером уже есть в базе"
	SettleErrorString              = "Не удалось заселить студента!"
	EvicErrorString                = "Не удалось выселить студента!"
	BadTransferErrorString         = "Не удалось осуществить действие"
	DBConnectErrorString           = "Не получилось подключиться к Базе данных!"
	IntErrorString                 = "Неправильно введено целое число!"
	InternalServerErrorString      = "Проблемы на стороне сервера"
	NotEnoughParamErrorString      = "Не указан ID студента для поиска!"
	MustBeIntErrorString           = "Параметр обязательно должен быть числом!"
	EmptyParamsErrorString         = "Параметр не должен быть пустой"
	WrongParamsErrorString         = "Параметры указаны неверно!"
)

type TransferDirection int