package logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const SessionIDLabel = "API-FRESH-PRODUCTS-SESSION"
const EventName = "API-FRESH-PRODUCTS-EVENT-NAME"

var log Logger = NoLogger{}

type Logger interface {
	Info(ctx context.Context, message string, details any)
	Error(ctx context.Context, message string, details any)
	Start(r *http.Request)
}

func SetLogger(l Logger) {
	log = l
}

func GetContext(r *http.Request, eventName string) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), SessionIDLabel, uuid.New().String()))
	return r.WithContext(context.WithValue(r.Context(), EventName, eventName))
}

func Info(ctx context.Context, message string, details any) {
	log.Info(ctx, message, details)
}
func Start(r *http.Request) {
	log.Start(r)
}
func Error(ctx context.Context, message string, details any) {
	log.Error(ctx, message, details)
}

type NoLogger struct{}

func (l NoLogger) Info(ctx context.Context, message string, details any) {
}
func (l NoLogger) Error(ctx context.Context, message string, details any) {
}
func (l NoLogger) Start(r *http.Request) {
}

type DBLogger struct {
	db *sql.DB
}

func NewDBLogger(db *sql.DB) Logger {
	return &DBLogger{db}
}
func (l DBLogger) Start(r *http.Request) {
	message := fmt.Sprintf("START: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
	log.Info(r.Context(), message, nil)
}
func (l DBLogger) Info(ctx context.Context, message string, details any) {
	sessionID := ctx.Value(SessionIDLabel).(string)
	eventName := ctx.Value(EventName).(string)
	now := time.Now().UnixMilli()

	var bContent []byte
	if details != nil {
		bContent, _ = json.Marshal(details)
	}

	_, err := l.db.Exec("INSERT INTO logs(session_id, event_name, message, details, date) VALUES (?, ?, ?, ?, ?)", sessionID, eventName, message, string(bContent), now)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (l DBLogger) Error(ctx context.Context, message string, details any) {
	sessionID := ctx.Value(SessionIDLabel).(string)
	eventName := ctx.Value(EventName).(string)
	now := time.Now().UnixMilli()

	var bContent []byte
	if details != nil {
		bContent, _ = json.Marshal(details)
	}

	_, err := l.db.Exec("INSERT INTO logs(session_id, event_name, message, details, date, is_error) VALUES (?, ?, ?, ?, ?, 1)", sessionID, eventName, message, string(bContent), now)
	if err != nil {
		fmt.Println(err.Error())
	}
}
