package handler_test

import (
	handler "main/internal/router/handlers"
	"testing"
	"time"
)

func TestGetDateString(t *testing.T) {
    date:=time.Date(2001,9,11,12,0,0,0,&time.Location{}) 
	want:="11.09.2001 12:00"

    t.Run("GetDateString",func(t *testing.T) {
		result:=handler.GetDateString(date)
		if (result != want){
			t.Errorf("Error date parse;wanted:%s;result:%s;",want,result)
		}
	})
}