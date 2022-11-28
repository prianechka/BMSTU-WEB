package objects

const (
	None                       = -1
	NotLiving                  = 0
	Ret                        = 0
	Get                        = 1
	EmptyString                = ""
	Empty                      = 0
	LoginInputMessage          = "Введите логин: "
	PasswordInputMessage       = "Введите пароль: "
	MarkInputMessage           = "Введите маркировочный номер: "
	TypeInputMessage           = "Введите тип вещи: "
	NameInputMessage           = "Введите имя студента: "
	SurnameInputMessage        = "Введите фамилию студента: "
	GroupInputMessage          = "Введите номер группы студента: "
	NumberInputMessage         = "Введите номер студенческого билета: "
	IDRoomInputMessage         = "Введите ID комнаты студента: "
	RoomIDInputMessage         = "Введите ID комнаты: "
	MarkNumInputMessage        = "Введите маркировочный номер вещи: "
	AuthOK                     = "\nАвторизация прошла успешно!"
	EvicStudentOK              = "\nСтудент успешно выселен!\n"
	SettleStudentOK            = "\nСтудент успешно заселён!\n"
	GiveThingOK                = "\nВещь передана!\n"
	ReturnThingOK              = "\nВещь отдана обратно!\n"
	AddOK                      = "\nОперация успешно проведена!"
	StudentNotFoundErrorString = "\nСтудент не найден!\n"
	ReturnThingErrorString     = "\nВещь и так была не у студента!\n"
	ThingNotFound              = "\nВещь не найдена!\n"
	GiveThingErrorString       = "\nВещь уже у другого студента!\n"
	SettleStudentErrorString   = "\nСтудент уже живёт в другой комнате!\n"
	EvicStudentErrorString     = "\nСтудент уже нигде не живёт! \n"
	UserNotFoundErrorString    = "\nПользователь не найден!"
	UserAddErrorString         = "Пользователя не удалось добавить!"
	StudentAddErrorString      = "Студента не удалось добавить!"
	StudentChangeErrorString   = "Не удалось изменить данные студента!"
	LoginErrorString           = "Такого логина не существует!"
	PasswordErrorString        = "Пароль введен неверно!"
	RoomNotFoundErrorString    = "Комната не найден!"
	ThingNotInRoomErrorString  = "Вещь не находится в этой комнате!"
	DeleteRoomErrorString      = "Не удалось удалить комнату!"
	DeleteThingErrorString     = "Не удалось удалить вещь!"
	AddThingErrorString        = "Не удалось добавить вещь!"
	SettleErrorString          = "Не удалось заселить студента!"
	EvicErrorString            = "Не удалось выселить студента!"
	BadTransferErrorString     = "Не удалось осуществить действие"
	DBConnectErrorString       = "Не получилось подключиться к Базе данных!"
	IntErrorString             = "Неправильно введено целое число!"
	InternalServerErrorString  = "Проблемы на стороне сервера"
)

type TransferDirection int
