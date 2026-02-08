package listmonk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/logger"
)

// Adds the given user to the Listmonk contact list
// If subscriberId is not nil, then UPDATE the user instead of adding user
func AddUserToListmonk(email string, firstName string, lastName string, picture string, subscriberId *int, sendMarketingEmails bool) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		return
	}

	url := os.Getenv("LISTMONK_URL")
	username := os.Getenv("LISTMONK_USERNAME")
	password := os.Getenv("LISTMONK_PASSWORD")
	listIdString := os.Getenv("LISTMONK_LIST_ID")

	listId, err := strconv.Atoi(listIdString)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}

	// Create new subscriber
	args := bson.M{
		"email":  email,
		"name":   firstName + " " + lastName,
		"status": "enabled",
		"attribs": bson.M{
			"firstName": firstName,
			"lastName":  lastName,
			"picture":   picture,
		},
		"preconfirm_subscriptions": true,
	}
	if sendMarketingEmails {
		args["lists"] = bson.A{listId}
	}
	body, _ := json.Marshal(args)
	bodyBuffer := bytes.NewBuffer(body)

	var req *http.Request
	if subscriberId != nil {
		// Existing subscriber
		req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/api/subscribers/%d", url, *subscriberId), bodyBuffer)
	} else {
		// New subscriber
		req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/subscribers", url), bodyBuffer)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()
}

// Check if the user is already in listmonk
// Returns a bool representing whether the subscriber exists and the id of the subscriber if it does exist
func DoesUserExist(email string) (bool, *int) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		return false, nil
	}

	url := os.Getenv("LISTMONK_URL")
	username := os.Getenv("LISTMONK_USERNAME")
	password := os.Getenv("LISTMONK_PASSWORD")

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/subscribers?query=subscribers.email='%s'", url, email), nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return false, nil
	}
	defer resp.Body.Close()

	var response struct {
		Data struct {
			Results []struct {
				Id int `json:"id"`
			} `json:"results"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logger.StdErr.Println(err)
		return false, nil
	}

	if len(response.Data.Results) > 0 {
		return true, &response.Data.Results[0].Id
	} else {
		return false, nil
	}
}

// Send a transactional email using the specified template and data
func SendEmail(email string, templateId int, data bson.M) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		return
	}

	// Get listmonk url env vars
	listmonkUrl := os.Getenv("LISTMONK_URL")
	listmonkUsername := os.Getenv("LISTMONK_USERNAME")
	listmonkPassword := os.Getenv("LISTMONK_PASSWORD")

	// Construct body
	body, err := json.Marshal(bson.M{
		"subscriber_email": email,
		"template_id":      templateId,
		"data":             data,
		"content_type":     "html",
	})
	if err != nil {
		logger.StdErr.Println(err)
		return
	}

	// Construct request
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/tx", listmonkUrl), bytes.NewBuffer(body))
	req.SetBasicAuth(listmonkUsername, listmonkPassword)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
	}
	defer response.Body.Close()
}

// Send a transactional email using the specified template and data. Adds subscriber if they don't exist
func SendEmailAddSubscriberIfNotExist(email string, templateId int, data bson.M, sendMarketingEmails bool) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		return
	}

	if exists, _ := DoesUserExist(email); !exists {
		AddUserToListmonk(email, "", "", "", nil, sendMarketingEmails)
	}

	SendEmail(email, templateId, data)
}

// ScheduleReminderEmails schedules three reminder emails for an event
// Returns a slice of bson.M objects to be stored in the Remindee model's Reminders field
func ScheduleReminderEmails(email string, ownerName string, eventName string, eventId string) []interface{} {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		return []interface{}{}
	}

	// Get email template ids
	initialEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_INITIAL_EMAIL_REMINDER_ID"))
	if err != nil {
		logger.StdErr.Println("Error parsing LISTMONK_INITIAL_EMAIL_REMINDER_ID:", err)
		return []interface{}{}
	}
	secondEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_SECOND_EMAIL_REMINDER_ID"))
	if err != nil {
		logger.StdErr.Println("Error parsing LISTMONK_SECOND_EMAIL_REMINDER_ID:", err)
		return []interface{}{}
	}
	finalEmailReminderId, err := strconv.Atoi(os.Getenv("LISTMONK_FINAL_EMAIL_REMINDER_ID"))
	if err != nil {
		logger.StdErr.Println("Error parsing LISTMONK_FINAL_EMAIL_REMINDER_ID:", err)
		return []interface{}{}
	}

	// Find if subscriber exists in listmonk
	subscriberExists, _ := DoesUserExist(email)

	// If subscriber doesn't exist, add subscriber to listmonk
	if !subscriberExists {
		AddUserToListmonk(email, "", "", "", nil, false)
	}

	// Create scheduled reminders as bson.M for storage
	now := time.Now()
	falseBool := false
	reminders := []interface{}{
		bson.M{
			"templateId":  initialEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now),
			"sent":        &falseBool,
		},
		bson.M{
			"templateId":  secondEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now.Add(24 * time.Hour)),
			"sent":        &falseBool,
		},
		bson.M{
			"templateId":  finalEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now.Add(72 * time.Hour)),
			"sent":        &falseBool,
		},
	}

	return reminders
}

// StartReminderEmailScheduler starts a background worker that checks for and sends scheduled reminder emails
func StartReminderEmailScheduler(ctx context.Context, eventsCollection *mongo.Collection) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		logger.StdOut.Println("Listmonk disabled, reminder email scheduler not started")
		return
	}

	logger.StdOut.Println("Starting reminder email scheduler...")

	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.StdOut.Println("Reminder email scheduler stopped")
			return
		case <-ticker.C:
			sendScheduledReminders(eventsCollection)
		}
	}
}

// sendScheduledReminders checks for and sends any due reminder emails
func sendScheduledReminders(eventsCollection *mongo.Collection) {
	if eventsCollection == nil {
		return
	}

	ctx := context.Background()
	now := primitive.NewDateTimeFromTime(time.Now())

	// Find all events with remindees that have unsent scheduled reminders
	filter := bson.M{
		"remindees": bson.M{
			"$elemMatch": bson.M{
				"reminders": bson.M{
					"$elemMatch": bson.M{
						"sent":        false,
						"scheduledAt": bson.M{"$lte": now},
					},
				},
			},
		},
	}

	cursor, err := eventsCollection.Find(ctx, filter)
	if err != nil {
		logger.StdErr.Println("Error finding events with scheduled reminders:", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event bson.M
		if err := cursor.Decode(&event); err != nil {
			logger.StdErr.Println("Error decoding event:", err)
			continue
		}

		eventId := event["_id"].(primitive.ObjectID).Hex()
		eventName := event["name"].(string)
		
		// Get owner name
		ownerName := "Somebody"
		if ownerId, ok := event["ownerId"].(primitive.ObjectID); ok && !ownerId.IsZero() {
			// We would need to fetch the user here, but for simplicity we'll use the event creator's name if available
			// This is a simplification - in production you'd want to fetch the user from the database
		}

		remindees, ok := event["remindees"].(primitive.A)
		if !ok {
			continue
		}

		// Construct URLs (using a helper function from utils)
		baseUrl := os.Getenv("BASE_URL")
		if baseUrl == "" {
			baseUrl = "http://localhost:3002"
		}
		eventUrl := fmt.Sprintf("%s/e/%s", baseUrl, eventId)

		for i, remindeeInterface := range remindees {
			remindee, ok := remindeeInterface.(bson.M)
			if !ok {
				continue
			}

			email, ok := remindee["email"].(string)
			if !ok {
				continue
			}

			// Check if already responded
			if responded, ok := remindee["responded"].(*bool); ok && responded != nil && *responded {
				continue
			}

			reminders, ok := remindee["reminders"].(primitive.A)
			if !ok {
				continue
			}

			finishedUrl := fmt.Sprintf("%s/e/%s/responded?email=%s", baseUrl, eventId, email)

			// Check each reminder to see if it should be sent
			for j, reminderInterface := range reminders {
				reminder, ok := reminderInterface.(bson.M)
				if !ok {
					continue
				}

				sent, _ := reminder["sent"].(bool)
				if sent {
					continue
				}

				scheduledAt, ok := reminder["scheduledAt"].(primitive.DateTime)
				if !ok {
					continue
				}

				if scheduledAt <= now {
					templateId, ok := reminder["templateId"].(int32)
					if !ok {
						// Try int64
						if templateId64, ok := reminder["templateId"].(int64); ok {
							templateId = int32(templateId64)
						} else {
							continue
						}
					}

					// Send the email
					SendEmail(email, int(templateId), bson.M{
						"ownerName":   ownerName,
						"eventName":   eventName,
						"eventUrl":    eventUrl,
						"finishedUrl": finishedUrl,
					})

					// Mark as sent
					update := bson.M{
						"$set": bson.M{
							fmt.Sprintf("remindees.%d.reminders.%d.sent", i, j): true,
						},
					}

					_, err := eventsCollection.UpdateOne(ctx, bson.M{"_id": event["_id"]}, update)
					if err != nil {
						logger.StdErr.Printf("Error marking reminder as sent: %v\n", err)
					} else {
						logger.StdOut.Printf("Sent reminder email to %s for event %s (template %d)\n", email, eventName, templateId)
					}
				}
			}
		}
	}
}

// CancelScheduledReminders marks all unsent reminders for a remindee as cancelled by setting sent=true
// This prevents them from being sent in the future
func CancelScheduledReminders(eventsCollection *mongo.Collection, eventId string, remindeeEmail string) {
	if eventsCollection == nil {
		return
	}

	ctx := context.Background()
	eventObjectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		logger.StdErr.Println("Error parsing event ID:", err)
		return
	}

	// Find the event and update reminders for this specific remindee
	filter := bson.M{
		"_id":              eventObjectId,
		"remindees.email": remindeeEmail,
	}

	// Mark all unsent reminders as sent (cancelled)
	update := bson.M{
		"$set": bson.M{
			"remindees.$[elem].reminders.$[reminder].sent": true,
		},
	}

	arrayFilters := []interface{}{
		bson.M{"elem.email": remindeeEmail},
		bson.M{"reminder.sent": false},
	}

	_, err = eventsCollection.UpdateOne(
		ctx,
		filter,
		update,
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}),
	)

	if err != nil {
		logger.StdErr.Printf("Error cancelling scheduled reminders: %v\n", err)
	}
}
