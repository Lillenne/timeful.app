package listmonk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/logger"
)

// GetTemplateID retrieves a template ID from environment variables with error handling
func GetTemplateID(envVar string, defaultValue int) int {
	if value := os.Getenv(envVar); value != "" {
		if id, err := strconv.Atoi(value); err == nil {
			return id
		}
		logger.StdErr.Printf("Warning: Invalid %s value '%s', using default %d\n", envVar, value, defaultValue)
	}
	if defaultValue == 0 {
		logger.StdErr.Printf("Warning: %s not set and no default provided\n", envVar)
	}
	return defaultValue
}

// Template ID environment variable names
const (
	EnvEveryoneRespondedTemplateID        = "LISTMONK_EVERYONE_RESPONDED_TEMPLATE_ID"
	EnvAvailabilityGroupInviteTemplateID  = "LISTMONK_AVAILABILITY_GROUP_INVITE_TEMPLATE_ID"
	EnvSomeoneRespondedEventTemplateID    = "LISTMONK_SOMEONE_RESPONDED_EVENT_TEMPLATE_ID"
	EnvAddedAttendeeTemplateID            = "LISTMONK_ADDED_ATTENDEE_TEMPLATE_ID"
	EnvSomeoneRespondedGroupTemplateID    = "LISTMONK_SOMEONE_RESPONDED_GROUP_TEMPLATE_ID"
	EnvXResponsesTemplateID               = "LISTMONK_X_RESPONSES_TEMPLATE_ID"
	EnvInitialEmailReminderID             = "LISTMONK_INITIAL_EMAIL_REMINDER_ID"
	EnvSecondEmailReminderID              = "LISTMONK_SECOND_EMAIL_REMINDER_ID"
	EnvFinalEmailReminderID               = "LISTMONK_FINAL_EMAIL_REMINDER_ID"
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
		logger.StdOut.Printf("Listmonk disabled - skipping email to %s (template %d)\n", email, templateId)
		return
	}

	// Get listmonk url env vars
	listmonkUrl := os.Getenv("LISTMONK_URL")
	listmonkUsername := os.Getenv("LISTMONK_USERNAME")
	listmonkPassword := os.Getenv("LISTMONK_PASSWORD")

	// Validate configuration
	if listmonkUrl == "" {
		logger.StdErr.Printf("ERROR: LISTMONK_URL not configured - cannot send email to %s\n", email)
		return
	}
	if listmonkUsername == "" || listmonkPassword == "" {
		logger.StdErr.Printf("ERROR: LISTMONK_USERNAME or LISTMONK_PASSWORD not configured - cannot send email to %s\n", email)
		return
	}

	logger.StdOut.Printf("Sending email to %s using template %d via Listmonk at %s\n", email, templateId, listmonkUrl)

	// Construct body using external mode to send to arbitrary email addresses
	// This doesn't require the email to be a subscriber in Listmonk
	body, err := json.Marshal(bson.M{
		"subscriber_mode":   "external",
		"subscriber_emails": []string{email},
		"template_id":       templateId,
		"data":              data,
		"content_type":      "html",
	})
	if err != nil {
		logger.StdErr.Printf("ERROR: Failed to marshal email request for %s: %v\n", email, err)
		return
	}

	logger.StdOut.Printf("Email request payload: %s\n", string(body))

	// Construct request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/tx", listmonkUrl), bytes.NewBuffer(body))
	if err != nil {
		logger.StdErr.Printf("ERROR: Failed to create HTTP request for %s: %v\n", email, err)
		return
	}
	req.SetBasicAuth(listmonkUsername, listmonkPassword)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Printf("ERROR: Failed to send email to %s via Listmonk: %v\n", email, err)
		return
	}
	defer response.Body.Close()

	// Read and log response
	bodyBytes, _ := io.ReadAll(response.Body)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		logger.StdOut.Printf("SUCCESS: Email sent to %s (status %d). Response: %s\n", email, response.StatusCode, string(bodyBytes))
	} else {
		logger.StdErr.Printf("ERROR: Listmonk returned status %d for email to %s. Response: %s\n", response.StatusCode, email, string(bodyBytes))
	}
}

// SendEmailAddSubscriberIfNotExist sends a transactional email using external mode
// Optionally adds the user as a subscriber for marketing emails if requested
// Uses external mode - does not require subscribers to exist in Listmonk for transactional emails
func SendEmailAddSubscriberIfNotExist(email string, templateId int, data bson.M, sendMarketingEmails bool) {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		logger.StdOut.Printf("Listmonk disabled - skipping email to %s (template %d)\n", email, templateId)
		return
	}

	logger.StdOut.Printf("Sending transactional email to %s (template %d, marketing: %v)\n", email, templateId, sendMarketingEmails)

	// Optionally add to marketing list if requested
	if sendMarketingEmails {
		if exists, _ := DoesUserExist(email); !exists {
			logger.StdOut.Printf("Adding %s to marketing list\n", email)
			AddUserToListmonk(email, "", "", "", nil, true)
		}
	}

	// Send transactional email using external mode
	SendEmail(email, templateId, data)
}

// ScheduleReminderEmails schedules three reminder emails for an event
// All reminders are scheduled via the background cron job
// Returns a slice of bson.M objects to be stored in the Remindee model's Reminders field
func ScheduleReminderEmails(email string, ownerName string, eventName string, eventId string) []interface{} {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		logger.StdOut.Printf("Listmonk disabled - skipping reminder setup for %s\n", email)
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

	logger.StdOut.Printf("Setting up reminders for %s (event: %s)\n", email, eventId)

	// Find if subscriber exists in listmonk
	subscriberExists, _ := DoesUserExist(email)

	// If subscriber doesn't exist, add subscriber to listmonk
	if !subscriberExists {
		AddUserToListmonk(email, "", "", "", nil, false)
	}

	// Create scheduled reminders as bson.M for storage
	// Schedule first reminder for immediate pickup (set time slightly in the past)
	// This ensures the background scheduler picks it up within 1 minute
	// Create separate boolean pointers for each reminder to avoid shared state
	now := time.Now()
	falseBool1 := false
	falseBool2 := false
	falseBool3 := false
	reminders := []interface{}{
		bson.M{
			"templateId":  initialEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now.Add(-1 * time.Second)), // Slightly in the past to ensure immediate pickup
			"sent":        &falseBool1,
		},
		bson.M{
			"templateId":  secondEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now.Add(24 * time.Hour)),
			"sent":        &falseBool2,
		},
		bson.M{
			"templateId":  finalEmailReminderId,
			"scheduledAt": primitive.NewDateTimeFromTime(now.Add(72 * time.Hour)),
			"sent":        &falseBool3,
		},
	}

	logger.StdOut.Printf("Scheduled 3 reminders for %s (immediate, 24h, and 72h)\n", email)

	return reminders
}

// StartReminderEmailScheduler starts a background worker using a cron scheduler
// to check for and send scheduled reminder emails every minute
func StartReminderEmailScheduler(ctx context.Context, eventsCollection *mongo.Collection, usersCollection *mongo.Collection) *cron.Cron {
	if os.Getenv("LISTMONK_ENABLED") == "false" {
		logger.StdOut.Println("Listmonk disabled, reminder email scheduler not started")
		return nil
	}

	logger.StdOut.Println("Starting reminder email scheduler using robfig/cron...")

	// Create a new cron scheduler
	c := cron.New()

	// Schedule the job to run every minute
	_, err := c.AddFunc("* * * * *", func() {
		sendScheduledReminders(eventsCollection, usersCollection)
	})
	if err != nil {
		logger.StdErr.Printf("Error scheduling reminder email job: %v\n", err)
		return nil
	}

	// Start the scheduler
	c.Start()

	// Stop the scheduler when context is cancelled
	go func() {
		<-ctx.Done()
		logger.StdOut.Println("Reminder email scheduler stopped")
		c.Stop()
	}()

	return c
}

// sendScheduledReminders checks for and sends any due reminder emails
func sendScheduledReminders(eventsCollection *mongo.Collection, usersCollection *mongo.Collection) {
	if eventsCollection == nil || usersCollection == nil {
		logger.StdErr.Println("ERROR: Collections not initialized for sendScheduledReminders")
		return
	}

	ctx := context.Background()
	now := primitive.NewDateTimeFromTime(time.Now())

	logger.StdOut.Printf("Checking for scheduled reminders at %v\n", time.Now())

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

	reminderCount := 0
	for cursor.Next(ctx) {
		var event bson.M
		if err := cursor.Decode(&event); err != nil {
			logger.StdErr.Println("Error decoding event:", err)
			continue
		}

		eventId := event["_id"].(primitive.ObjectID).Hex()
		eventName := event["name"].(string)
		
		// Get owner name by fetching from database
		ownerName := "Somebody"
		if ownerId, ok := event["ownerId"].(primitive.ObjectID); ok && !ownerId.IsZero() {
			var user bson.M
			err := usersCollection.FindOne(ctx, bson.M{"_id": ownerId}).Decode(&user)
			if err == nil {
				if firstName, ok := user["firstName"].(string); ok {
					ownerName = firstName
				}
			} else {
				logger.StdErr.Printf("Warning: Could not fetch user %s for event %s: %v\n", ownerId.Hex(), eventId, err)
			}
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

				// Check if reminder has been sent
				// Handle both bool and *bool types for backwards compatibility
				var sent bool
				if sentBool, ok := reminder["sent"].(bool); ok {
					sent = sentBool
				} else if sentBoolPtr, ok := reminder["sent"].(*bool); ok && sentBoolPtr != nil {
					sent = *sentBoolPtr
				}
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
							logger.StdErr.Printf("Warning: Invalid templateId type for reminder in event %s\n", eventId)
							continue
						}
					}

					reminderCount++
					logger.StdOut.Printf("Processing reminder %d for email %s, event %s (template %d)\n", reminderCount, email, eventName, templateId)

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
						logger.StdOut.Printf("Marked reminder as sent for %s in event %s\n", email, eventName)
					}
				}
			}
		}
	}
	
	if reminderCount == 0 {
		logger.StdOut.Println("No pending reminders found")
	} else {
		logger.StdOut.Printf("Processed %d reminder(s)\n", reminderCount)
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
