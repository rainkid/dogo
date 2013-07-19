package dogo

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Util_UnixTime(t time.Time) string {
	ftime := t.Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}

func Util_UCFirst(s string) string {
	str := strings.Split(strings.ToLower(s), "")
	str[0] = strings.ToUpper(str[0])
	return strings.Join(str, "")
}

func Util_FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Util_Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Util_Log(prefix, format string, args ...interface{}) {

	loghd, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE, 644)
	if err != nil {
		return
	}
	defer loghd.Close()
	logger := log.New(loghd, prefix, log.Ldate|log.Ltime)
	logger.Fatalf(format, args)
}
