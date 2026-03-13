package api

import (
	"final-project/pkg/logic"
	"fmt"
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time

	if nowStr == "" {
		now = time.Now()
	} else {
		var err error

		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	}

	next, err := logic.NextDate(now, dateStr, repeat)
	if err != nil {
		if err == logic.ErrNoRepeat {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, next)
}
