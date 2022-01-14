package echo

import (
	"fmt"
	"io"
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

	if c != nil {
		if c.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", fmt.Sprintf("%v", c.Get("UserId")))
		} else {
			userId = fmt.Sprintf("%v", c.Get("RequestID"))
		}

		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub := sentryecho.GetHubFromContext(c)

			dataBreadcumb := breadcumb

			dataBreadcumb.Data = data
			dataBreadcumb.Message = fmt.Sprintf("%v", message)

			hub.CaptureMessage(fmt.Sprintf("%v", message))
			hub.AddBreadcrumb(&breadcumb, nil)
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
