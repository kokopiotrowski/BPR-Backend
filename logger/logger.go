package logger

import (
	"log"
	"net/http"
	"os"
	"stockx-backend/conf"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		if conf.FlagConf.IsLoggingOn {
			//year, month, day := start.Date()

			//absPath, _ := filepath.Abs(strconv.FormatInt(int64(year), 10) + "/" + month.String() + "/" + strconv.FormatInt(int64(day), 10) + ".log")
			f, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()

			log.SetOutput(f)
		}

		if r.RequestURI != "/index" {
			log.Printf(
				"%s %s %s %s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start),
			)
		}
	})
}
