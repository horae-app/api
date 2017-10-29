package Cron

import (
	"time"
	"log"

	calendar "github.com/horae-app/api/Calendar"
	device "github.com/horae-app/api/Device"
	settings "github.com/horae-app/api/Settings"
	util "github.com/horae-app/api/Util"
)

const DAY_INTERVAL time.Duration = 24 * time.Hour

func SentReminderToContacts() {
	ticker := updateTicker(settings.CONTACT_REMEMBER_HOUR, settings.CONTACT_REMEMBER_MINUTE, 0, DAY_INTERVAL)
	for {
		<-ticker.C
		log.Println("Sending notifications...")
		for _, cal := range calendar.GetAllTomorrow(time.Now()) {
			token := device.GetTokenByEmail(cal.Contact.Email)
			if token == "" {
				log.Println("[Notify Auto] Error Token not found to", cal.Contact.Email)
				continue
			}
			err := util.SentContactRemember(token, cal.ID.String())
			if err != nil {
				log.Println("[Notify Auto] Error", token, err)
				continue
			}
			log.Println("[Notify Auto] Success", cal.ID)
		}

		ticker = updateTicker(settings.CONTACT_REMEMBER_HOUR, settings.CONTACT_REMEMBER_MINUTE, 0, DAY_INTERVAL)
	}
}

func updateTicker(hour int, minute int, second int, period time.Duration) *time.Ticker {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, second, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(period)
	}
	log.Println("Notify Auto] Next ", nextTick)
	diff := nextTick.Sub(time.Now())
	return time.NewTicker(diff)
}


func Start() {
	go SentReminderToContacts()
}
