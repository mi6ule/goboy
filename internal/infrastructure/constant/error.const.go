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
)

var ErrorMessage = map[string]string{
	"100001": "could not import env variables",
	"100002": "failed to run migrations",
	"100003": "failed to connect to mongoDb",
	"100004": "failed to connect to redis",
	"100005": "",
	"100006": "",
	"100007": "",
	"100008": "invalid connection string",
	"100009": "Error while trying to connect to mongodb",
	"100010": "Error while trying to connect to mongodb",
	"100011": "",
	"100012": "",
	"100013": "Could not run asynq scheduler",
	"100014": "Could not init mux server",
	"100015": "Could not enqueue email",
	"100016": "Could not enqueue image resize",
	"100017": "could not connect to postgresql db",
	"100018": "",
	"100019": "error creating elastic client",
}
