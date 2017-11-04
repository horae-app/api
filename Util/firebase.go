package Util

import (
    "github.com/NaySoftware/go-fcm"

    settings "github.com/horae-app/api/Settings"
)


func SentContactRemember(token string, calendar_id string) (err error) {
    data := map[string]string{
        "calendar_id": calendar_id,
    }

    c := fcm.NewFcmClient(settings.FCM_SERVER_KEY)
    c.NewFcmRegIdsMsg([]string{token}, data)
    
    status, err := c.Send()
    status.PrintResults()
    return err
}
