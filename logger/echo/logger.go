package echo

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logDump "github.com/sirupsen/logrus"
)

var (
	userId = ""
)

func InitBodyDumpLog() (err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	logf, err := rotatelogs.New(
		dir+"/logs/RequestResponseDump.log.%Y%m%d",
		rotatelogs.WithLinkName(dir+"/logs/RequestResponseDump.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(365),
	)

	logDump.SetFormatter(&logDump.JSONFormatter{DisableHTMLEscape: true})
	logDump.SetOutput(io.MultiWriter(os.Stdout, logf))
	logDump.SetLevel(logDump.InfoLevel)
	logDump.SetReportCaller(true)

	return
}

func Info(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Info(message)

	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func Error(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Error(message)
	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func SentryLog(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message string) {
	if c != nil {
		if c.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", fmt.Sprintf("%v", c.Get("UserId")))
		} else {
			userId = fmt.Sprintf("%v", c.Get("RequestID"))
		}

		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub := sentryecho.GetHubFromContext(c)

			dataBreadcumb := breadcumb
			log.Println(data)
			dataBreadcumb.Data = data
			dataBreadcumb.Message = message

			hub.CaptureMessage(message)
			sentry.AddBreadcrumb(&breadcumb)
		}
	}

}

// func Error(data map[string]interface{}, message interface{}) {
// 	logDump.WithFields(logDump.Fields{
// 		data,
// 	}).Error(message)
// }

// func Fatal(data map[string]interface{}, message interface{}) {
// 	logDump.WithFields(logDump.Fields{
// 		data,
// 	}).Fatal(message)
// }

// func Debug(data map[string]interface{}, message interface{}) {
// 	logDump.WithFields(logDump.Fields{
// 		data,
// 	}).Debug(message)
// }
