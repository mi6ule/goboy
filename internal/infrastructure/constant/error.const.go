package constants

const (
	ERROR_CODE_100001 = "100001"
	ERROR_CODE_100002 = "100002"
	ERROR_CODE_100003 = "100003"
	ERROR_CODE_100004 = "100004"
	ERROR_CODE_100005 = "100005"
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
	ERROR_CODE_100023 = "100023"
	ERROR_CODE_100024 = "100024"
	ERROR_CODE_100025 = "100025"
	ERROR_CODE_100026 = "100026"
	ERROR_CODE_100027 = "100027"
	ERROR_CODE_100028 = "100028"
	ERROR_CODE_100029 = "100029"
	ERROR_CODE_100030 = "100030"
	ERROR_CODE_100031 = "100031"
	ERROR_CODE_100032 = "100032"
	ERROR_CODE_100033 = "100033"
	ERROR_CODE_100034 = "100034"
	ERROR_CODE_100035 = "100035"
	ERROR_CODE_100036 = "100036"
)

var ErrorMessage = map[string]string{
	ERROR_CODE_100001: "failed to import env variables",
	ERROR_CODE_100002: "failed to run migrations",
	ERROR_CODE_100003: "failed to connect to mongoDb",
	ERROR_CODE_100004: "failed to connect to redis",
	ERROR_CODE_100005: "",
	ERROR_CODE_100006: "",
	ERROR_CODE_100007: "",
	ERROR_CODE_100008: "invalid connection string",
	ERROR_CODE_100009: "failed to connect to mongodb",
	ERROR_CODE_100010: "failed to connect to mongodb",
	ERROR_CODE_100011: "",
	ERROR_CODE_100012: "",
	ERROR_CODE_100013: "failed to run asynq scheduler",
	ERROR_CODE_100014: "failed to init mux server",
	ERROR_CODE_100015: "failed to enqueue email",
	ERROR_CODE_100016: "failed to enqueue image resize",
	ERROR_CODE_100017: "failed to connect to postgresql db",
	ERROR_CODE_100018: "failed to listen gin server",
	ERROR_CODE_100019: "failed to create elastic client",
	ERROR_CODE_100020: "",
	ERROR_CODE_100021: "failed to create admin client",
	ERROR_CODE_100022: "failed to create topics",
	ERROR_CODE_100023: "KAFKA: Error in creating producer",
	ERROR_CODE_100024: "KAFKA: Error in creating producer",
	ERROR_CODE_100025: "Error in person producer",
	ERROR_CODE_100026: "delivery failed",
	ERROR_CODE_100027: "KAFKA: Error in subscribing to the topic",
	ERROR_CODE_100028: "KAFKA: Error in fetch messages from the topic",
	ERROR_CODE_100029: "",
	ERROR_CODE_100030: "",
	ERROR_CODE_100031: "failed to enqueue task",
	ERROR_CODE_100032: "failed to delete tasks from the source queue",
	ERROR_CODE_100033: "failed to parse task payload into json",
	ERROR_CODE_100034: "failed to get pending tasks from source queue",
	ERROR_CODE_100035: "failed to push task to destination queue",
	ERROR_CODE_100036: "failed to get queue info",
}
