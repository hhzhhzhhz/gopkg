package xgin

// ErrorParams params 100-150
var (
	ErrorParams      = NewResponse(100, "invalid parameter.")
	ErrorBodyRead    = NewResponse(102, "http body read failed.")
	ErrorBodyDecode  = NewResponse(103, "http body decode failed.")
	ErrorBodyMarshal = NewResponse(104, "json marshal failed.")
)

// mysql 150 - 250
var (
	SqlCreateError    = NewResponse(150, "db insert failed.")
	SqlDeleteError    = NewResponse(151, "db delete failed.")
	SqlUpdateError    = NewResponse(152, "db update failed.")
	SqlQueryError     = NewResponse(153, "db query failed.")
	SqlUnChangedError = NewResponse(154, "db unchanged.")
)

// PublishError business 250 - 399
var (
	PublishError = NewResponse(251, "publish message failed.")
)

// other >=400
var (
	Error404            = NewResponse(404, "404 not found")
	Error405            = NewResponse(405, "method does not allow")
	ErrorServerClosed   = NewResponse(401, "server is closing.")
	ErrorForwarderUrl   = NewResponse(500, "forwarder url parse failed.")
	ErrorNotSchedule    = NewResponse(501, "service is not schedule")
	ErrorNotEdge        = NewResponse(502, "service is not edge")
	ErrorUnknown        = NewResponse(503, "unknown error")
	ErrorServerNotReady = NewResponse(504, "server not ready.")
)
