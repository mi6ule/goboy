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
)

var ClientSideErrors = map[any]any{
	"BAD_REQUEST":                     400,
	400:                               "BAD_REQUEST",
	"UNAUTHORIZED":                    401,
	401:                               "UNAUTHORIZED",
	"PAYMENT_REQUIRED":                402,
	402:                               "PAYMENT_REQUIRED",
	"FORBIDDEN":                       403,
	403:                               "FORBIDDEN",
	"NOT_FOUND":                       404,
	404:                               "NOT_FOUND",
	"METHOD_NOT_ALLOWED":              405,
	405:                               "METHOD_NOT_ALLOWED",
	"NOT_ACCEPTABLE":                  406,
	406:                               "NOT_ACCEPTABLE",
	"PROXY_AUTHENTICATION":            407,
	407:                               "PROXY_AUTHENTICATION",
	"REQUEST_TIMEOUT":                 408,
	408:                               "REQUEST_TIMEOUT",
	"CONFLICT":                        409,
	409:                               "CONFLICT",
	"GONE":                            410,
	410:                               "GONE",
	"LENGTH_REQUIRED":                 411,
	411:                               "LENGTH_REQUIRED",
	"PRECONDITION_FAILED":             412,
	412:                               "PRECONDITION_FAILED",
	"PAYLOAD_TOO_LARGE":               413,
	413:                               "PAYLOAD_TOO_LARGE",
	"URI_TOO_LONG":                    414,
	414:                               "URI_TOO_LONG",
	"UNSUPPORTED_MEDIA_TYPE":          415,
	415:                               "UNSUPPORTED_MEDIA_TYPE",
	"RANGE_NOT_SATISFIABLE":           416,
	416:                               "RANGE_NOT_SATISFIABLE",
	"EXPECTATION_FAILED":              417,
	417:                               "EXPECTATION_FAILED",
	"IM_A_TEAPOT":                     418,
	418:                               "IM_A_TEAPOT",
	"UNPROCESSABLE_ENTITY":            422,
	422:                               "UNPROCESSABLE_ENTITY",
	"LOCKED":                          423,
	423:                               "LOCKED",
	"FAILED_DEPENDENCY":               424,
	424:                               "FAILED_DEPENDENCY",
	"UPGRADE_REQUIRED":                426,
	426:                               "UPGRADE_REQUIRED",
	"PRECONDITION_REQUIRED":           428,
	428:                               "PRECONDITION_REQUIRED",
	"TOO_MANY_REQUESTS":               429,
	429:                               "TOO_MANY_REQUESTS",
	"REQUEST_HEADER_FIELDS_TOO_LARGE": 431,
	431:                               "REQUEST_HEADER_FIELDS_TOO_LARGE",
	"UNAVAILABLE_FOR_LEGAL_REASONS":   451,
	451:                               "UNAVAILABLE_FOR_LEGAL_REASONS",
}

var ServerSideErrors = map[any]any{
	"INTERNAL_SERVER_ERROR":           500,
	500:                               "INTERNAL_SERVER_ERROR",
	"NOT_IMPLEMENTED":                 501,
	501:                               "NOT_IMPLEMENTED",
	"BAD_GATEWAY":                     502,
	502:                               "BAD_GATEWAY",
	"SERVICE_UNAVAILABLE":             503,
	503:                               "SERVICE_UNAVAILABLE",
	"GATEWAY_TIMEOUT":                 504,
	504:                               "GATEWAY_TIMEOUT",
	"HTTP_VERSION_NOT_SUPPORTED":      505,
	505:                               "HTTP_VERSION_NOT_SUPPORTED",
	"VARIANT_ALSO_NEGOTIATES":         506,
	506:                               "VARIANT_ALSO_NEGOTIATES",
	"INSUFFICIENT_STORAGE":            507,
	507:                               "INSUFFICIENT_STORAGE",
	"LOOP_DETECTED":                   508,
	508:                               "LOOP_DETECTED",
	"NOT_EXTENDED":                    510,
	510:                               "NOT_EXTENDED",
	"NETWORK_AUTHENTICATION_REQUIRED": 511,
	511:                               "NETWORK_AUTHENTICATION_REQUIRED",
}
