package api

import (
	"fmt"
	"net/http"
	"time"

	"final-project/pkg/logic"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time

	if nowStr == "" {
		now = time.Now()
	} else {
		var err error

		now, err = time.Parse(logic.DateFormat, nowStr)
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
