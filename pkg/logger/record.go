package logger

import (
	"api/pkg/tracer"
	"context"
	"errors"
)

// Record stores log data that will be sent to graylog and written to std.out
type Record struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
	Error   error                  `json:"error"`
	Context context.Context        `json:"-"`
	Type    string                 `json:"type"`
	Where   string                 `json:"where"`
}

func (rec *Record) validate() {
	if rec.Context == nil {
		rec.Context = context.TODO()
	}

	if rec.Type == "" {
		rec.Type = "DEBUG"
	}

	if rec.Data == nil {
		rec.Data = make(map[string]interface{})
	}

	if rec.Type == "ERROR" && rec.Error == nil {
		rec.Error = errors.New("Unknown error")
	}
}

func (rec *Record) addValuesFromContext() {
	traceID := tracer.GetTraceID(rec.Context)
	if traceID != "" {
		rec.Data["trace_id"] = traceID
	}

	sessionID := tracer.GetSessionID(rec.Context)
	if sessionID != "" {
		rec.Data["session_id"] = sessionID
	}

	spanName := tracer.GetSpanName(rec.Context)
	if spanName != "" {
		rec.Data["span_name"] = spanName
	}
}
