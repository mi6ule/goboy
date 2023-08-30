package constants

const (
	ERROR_CODE_100001 = "100001"
	ERROR_CODE_100004 = "100004"
	ERROR_CODE_100006 = "100006"
	ERROR_CODE_100007 = "100007"
	ERROR_CODE_100008 = "100008"
	ERROR_CODE_100009 = "100009"
	ERROR_CODE_100010 = "100010"
	ERROR_CODE_100011 = "100011"
	ERROR_CODE_100012 = "100012"
	ERROR_CODE_100013 = "100013"
	ERROR_CODE_100014 = "100014"
	ERROR_CODE_100015 = "100015"
	ERROR_CODE_100016 = "100016"
	ERROR_CODE_100017 = "100017"
	ERROR_CODE_100018 = "100018"
	ERROR_CODE_100019 = "100019"

	ERROR_CODE_100020 = "100020"
	ERROR_CODE_100021 = "100021"
	ERROR_CODE_100022 = "100022"
	ERROR_CODE_100023 = "100035"
	ERROR_CODE_100024 = "100024"
	ERROR_CODE_100025 = "100025"
	ERROR_CODE_100026 = "100026"
	ERROR_CODE_100027 = "100027"
	ERROR_CODE_100028 = "100028"
	ERROR_CODE_100029 = "100029"
	ERROR_CODE_100030 = "100030"
)

var ErrorMessage = map[string]string{
	"100001": "failed to import env variables",
	"100002": "failed to run migrations",
	"100003": "failed to connect to mongoDb",
	"100004": "failed to connect to redis",
	"100005": "",
	"100006": "",
	"100007": "",
	"100008": "invalid connection string",
	"100009": "failed to connect to mongodb",
	"100010": "failed to connect to mongodb",
	"100011": "",
	"100012": "",
	"100013": "failed to run asynq scheduler",
	"100014": "failed to init mux server",
	"100015": "failed to enqueue email",
	"100016": "failed to enqueue image resize",
	"100017": "failed to connect to postgresql db",
	"100018": "",
	"100019": "failed to create elastic client",
}
