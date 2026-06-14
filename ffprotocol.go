package main

//структуры запросов

// BaseRequest содержит базовые поля для всех POST-запросов
type BaseRequest struct {
	SerialNumber string `json:"serialNumber"` // Уникальный серийный номер устройства
	CheckCode    string `json:"checkCode"`    // Код проверки доступа
}

// CheckCodeRequest используется для валидации доступа по LAN
type CheckCodeRequest struct {
	BaseRequest
}

// DetailRequest используется для получения информации о принтере
type DetailRequest struct {
	BaseRequest
}

// ProductRequest используется для получения информации о возможностях устройства
type ProductRequest struct {
	BaseRequest
}

// ControlRequest представляет запрос на управление принтером
type ControlRequest struct {
	BaseRequest
	Payload interface{} `json:"payload"` // Полезная нагрузка команды
}

// Возможные типы Payload для команды управления
/*type Payload interface {
    LightControlCmd
    CirculateCtlCmd
    TemperatureCtlCmd
    PrinterCtlCmd
    JobCtlCmd
    DelayCloseCmd
    ReNameCmd
    CalibrationCmd
    StreamCtrlCmd
    UserProfileCmd
    StateCtrlCmd
    DeviceUpdateDetailCmd
    NewLocalJobCmd
}*/

// Структура для запроса на печать G-code файла
type PrintGcodeRequest struct {
	BaseRequest
	FileName            string `json:"fileName"`            // Имя файла для печати
	LevelingBeforePrint bool   `json:"levelingBeforePrint"` // Флаг выравнивания перед печатью
}

// LightControlCmd управление освещением
type LightControlCmd struct {
	Cmd  string    `json:"cmd"`  // Команда: lightControl_cmd
	Args LightArgs `json:"args"` // Параметры команды
}

type LightArgs struct {
	Status string `json:"status"` // Статус освещения (open/close)
}

// CirculateCtlCmd представляет команду управления вентиляторами
type CirculateCtlCmd struct {
	Cmd  string        `json:"cmd"`
	Args CirculateArgs `json:"args"`
}

type CirculateArgs struct {
	Internal string `json:"internal"` // Состояние внутреннего вентилятора
	External string `json:"external"` // Состояние внешнего вентилятора
}

// TemperatureCtlCmd управление температурой
type TemperatureCtlCmd struct {
	Cmd  string          `json:"cmd"`  // Команда: temperatureCtl_cmd
	Args TemperatureArgs `json:"args"` // Параметры команды
}

type TemperatureArgs struct {
	Platform    int `json:"platform"`    // Температура платформы
	RightNozzle int `json:"rightNozzle"` // Температура правого сопла
	LeftNozzle  int `json:"leftNozzle"`  // Температура левого сопла
	Chamber     int `json:"chamber"`     // Температура камеры
}

// JobCtlCmd управление заданиями печати
type JobCtlCmd struct {
	Cmd  string  `json:"cmd"`  // Команда: jobCtl_cmd
	Args JobArgs `json:"args"` // Параметры команды
}

type JobArgs struct {
	JobID  string `json:"jobID"`  // ID задания
	Action string `json:"action"` // Действие (pause, continue, resume, cancel, stop)
}

// GcodeListRequest используется для получения списка файлов для печати
type GcodeListRequest struct {
	BaseRequest
}

// GcodeThumbRequest используется для получения превью файла
type GcodeThumbRequest struct {
	BaseRequest
	FileName string `json:"fileName"` // Имя файла для получения превью
}

//Структуры ответов

// Базовый ответ с кодом и сообщением
type CodeMessageResponse struct {
	Code    int    `json:"code"`    // Код результата операции
	Message string `json:"message"` // Описание результата
}

// Ответ на запрос детальной информации
type DetailResponse struct {
	Code    int     `json:"code"`    // Код результата
	Message string  `json:"message"` // Описание результата
	Detail  Details `json:"detail"`  // Подробная информация о принтере
}

// Ответ с информацией о возможностях устройства
type ProductResponse struct {
	Code    int                 `json:"code"`    // Код результата
	Message string              `json:"message"` // Описание результата
	Product ProductCapabilities `json:"product"` // Возможности устройства
}

type ProductCapabilities struct {
	NozzleTempCtrlState   int `json:"nozzleTempCtrlState"`   // Состояние управления температурой сопла
	ChamberTempCtrlState  int `json:"chamberTempCtrlState"`  // Состояние управления температурой камеры
	PlatformTempCtrlState int `json:"platformTempCtrlState"` // Состояние управления температурой платформы
	LightCtrlState        int `json:"lightCtrlState"`        // Состояние управления освещением
	InternalFanCtrlState  int `json:"internalFanCtrlState"`  // Состояние управления внутренним вентилятором
	ExternalFanCtrlState  int `json:"externalFanCtrlState"`  // Состояние управления внешним вентилятором
}

// Ответ со списком файлов для печати
type GcodeListResponse struct {
	Code      int      `json:"code"`      // Код результата
	Message   string   `json:"message"`   // Описание результата
	GcodeList []string `json:"gcodeList"` // Список файлов
}

// Ответ с превью файла
type GcodeThumbResponse struct {
	Code      int    `json:"code"`      // Код результата
	Message   string `json:"message"`   // Описание результата
	ImageData string `json:"imageData"` // Превью в формате base64
}

// Details представляет подробную информацию о состоянии и конфигурации 3D-принтера
type Details struct {
	// Настройки автоматического выключения
	AutoShutdown     string `json:"autoShutdown"`     // Возможные значения: "open", "close"
	AutoShutdownTime int    `json:"autoShutdownTime"` // Время до автоматического выключения

	// Параметры потоковой передачи камеры
	CameraStreamURL string `json:"cameraStreamUrl"` // URL для потоковой передачи

	// Параметры вентилятора камеры
	ChamberFanSpeed   int     `json:"chamberFanSpeed"`   // Скорость вентилятора камеры
	ChamberTargetTemp float64 `json:"chamberTargetTemp"` // Целевая температура камеры
	ChamberTemp       float64 `json:"chamberTemp"`       // Текущая температура камеры

	// Параметры охлаждающего вентилятора
	CoolingFanSpeed int `json:"coolingFanSpeed"`

	// Статистика использования филамента
	CumulativeFilament  float64 `json:"cumulativeFilament"`  // Общий расход филамента
	CumulativePrintTime int     `json:"cumulativePrintTime"` // Общее время печати

	// Параметры текущей печати
	CurrentPrintSpeed    int     `json:"currentPrintSpeed"`
	DoorStatus           string  `json:"doorStatus"`           // Статус дверцы: "open", "close"
	ErrorCode            string  `json:"errorCode"`            // Код ошибки
	EstimatedLeftLen     float64 `json:"estimatedLeftLen"`     // Ожидаемая длина левого филамента
	EstimatedLeftWeight  float64 `json:"estimatedLeftWeight"`  // Ожидаемый вес левого филамента
	EstimatedRightLen    float64 `json:"estimatedRightLen"`    // Ожидаемая длина правого филамента
	EstimatedRightWeight float64 `json:"estimatedRightWeight"` // Ожидаемый вес правого филамента
	EstimatedTime        float64 `json:"estimatedTime"`        // Ожидаемое время печати

	// Статус внешних вентиляторов
	ExternalFanStatus string `json:"externalFanStatus"` // "open", "close"

	// Дополнительные параметры
	FillAmount         int     `json:"fillAmount"`      // Уровень заполнения
	FirmwareVersion    string  `json:"firmwareVersion"` // Версия прошивки
	FlashRegisterCode  string  `json:"flashRegisterCode"`
	InternalFanStatus  string  `json:"internalFanStatus"` // "open", "close"
	IPAddr             string  `json:"ipAddr"`            // IP-адрес принтера
	LeftFilamentType   string  `json:"leftFilamentType"`  // Тип левого филамента
	LeftTargetTemp     float64 `json:"leftTargetTemp"`    // Целевая температура левого сопла
	LeftTemp           float64 `json:"leftTemp"`          // Текущая температура левого сопла
	LightStatus        string  `json:"lightStatus"`       // Статус подсветки: "open", "close"
	Location           string  `json:"location"`          // Местоположение
	MACAddr            string  `json:"macAddr"`           // MAC-адрес
	Measure            string  `json:"measure"`           // Единицы измерения
	Name               string  `json:"name"`              // Имя принтера
	NozzleCnt          int     `json:"nozzleCnt"`         // Количество сопел
	NozzleModel        string  `json:"nozzleModel"`       // Модель сопла
	NozzleStyle        int     `json:"nozzleStyle"`
	PID                int     `json:"pid"`
	PlatTargetTemp     float64 `json:"platTargetTemp"` // Целевая температура платформы
	PlatTemp           float64 `json:"platTemp"`       // Текущая температура платформы
	PolarRegisterCode  string  `json:"polarRegisterCode"`
	PrintDuration      int     `json:"printDuration"`      // Продолжительность печати
	PrintFileName      string  `json:"printFileName"`      // Имя печатаемого файла
	PrintFileThumbURL  string  `json:"printFileThumbUrl"`  // URL превью файла
	PrintLayer         int     `json:"printLayer"`         // Текущий слой
	PrintProgress      float64 `json:"printProgress"`      // Прогресс печати
	PrintSpeedAdjust   float64 `json:"printSpeedAdjust"`   // Регулировка скорости печати
	RemainingDiskSpace float64 `json:"remainingDiskSpace"` // Свободное место на диске
	RightFilamentType  string  `json:"rightFilamentType"`  // Тип правого филамента
	RightTargetTemp    float64 `json:"rightTargetTemp"`    // Целевая температура правого сопла
	RightTemp          float64 `json:"rightTemp"`          // Текущая температура правого сопла
	Status             string  `json:"status"`             // Общий статус принтера
	TargetPrintLayer   int     `json:"targetPrintLayer"`   // Целевой слой
	TVOC               int     `json:"tvoc"`               // Уровень TVOC
	ZAxisCompensation  float64 `json:"zAxisCompensation"`  // Компенсация по оси Z
}
