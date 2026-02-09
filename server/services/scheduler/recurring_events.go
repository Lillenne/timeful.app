package scheduler

import (
"context"
"fmt"
"time"

"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
"schej.it/server/db"
"schej.it/server/logger"
"schej.it/server/models"
"schej.it/server/services/gcloud"
"schej.it/server/utils"
)

// calculateNextDate calculates the next date based on the recurrence interval and unit
func calculateNextDate(baseTime time.Time, interval int, unit string) (time.Time, error) {
switch unit {
case "days":
return baseTime.AddDate(0, 0, interval), nil
case "weeks":
return baseTime.AddDate(0, 0, interval*7), nil
case "months":
return baseTime.AddDate(0, interval, 0), nil
default:
return time.Time{}, fmt.Errorf("invalid recurrence unit: %s", unit)
}
}

// StartRecurringEventScheduler starts a scheduler that checks every hour for recurring events that need new instances created
func StartRecurringEventScheduler() {
logger.StdOut.Println("Starting recurring event scheduler...")

// Run immediately on startup
go checkAndCreateRecurringEvents()

// Then run every hour
ticker := time.NewTicker(1 * time.Hour)
go func() {
for range ticker.C {
checkAndCreateRecurringEvents()
}
}()
}

// checkAndCreateRecurringEvents checks for recurring events that need new instances created
func checkAndCreateRecurringEvents() {
logger.StdOut.Println("Checking for recurring events to create...")

// Find all recurring events where:
// 1. IsRecurring is true
// 2. RecurrenceEnabled is true
// 3. NextOccurrenceDate is less than or equal to now
now := primitive.NewDateTimeFromTime(time.Now())

cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
"isRecurring": true,
"recurrenceEnabled": true,
"nextOccurrenceDate": bson.M{
"$lte": now,
},
"$or": bson.A{
bson.M{"isDeleted": bson.M{"$exists": false}},
bson.M{"isDeleted": false},
},
})

if err != nil {
logger.StdErr.Printf("Error finding recurring events: %v\n", err)
return
}
defer cursor.Close(context.Background())

var events []models.Event
if err := cursor.All(context.Background(), &events); err != nil {
logger.StdErr.Printf("Error decoding recurring events: %v\n", err)
return
}

logger.StdOut.Printf("Found %d recurring events to process\n", len(events))

for _, event := range events {
if err := createNextRecurringEvent(&event); err != nil {
logger.StdErr.Printf("Error creating recurring event for %s: %v\n", event.Id.Hex(), err)
}
}
}

// createNextRecurringEvent creates the next instance of a recurring event
func createNextRecurringEvent(parentEvent *models.Event) error {
logger.StdOut.Printf("Creating next occurrence for event %s\n", parentEvent.Id.Hex())

// Calculate the new dates based on recurrence interval
interval := utils.Coalesce(parentEvent.RecurrenceInterval)
unit := utils.Coalesce(parentEvent.RecurrenceUnit)

if interval <= 0 || unit == "" {
return fmt.Errorf("invalid recurrence configuration")
}

newDates := make([]primitive.DateTime, len(parentEvent.Dates))
for i, date := range parentEvent.Dates {
newTime, err := calculateNextDate(date.Time(), interval, unit)
if err != nil {
return err
}
newDates[i] = primitive.NewDateTimeFromTime(newTime)
}

// Calculate new times if they exist
var newTimes []primitive.DateTime
if parentEvent.Times != nil && len(parentEvent.Times) > 0 {
newTimes = make([]primitive.DateTime, len(parentEvent.Times))
for i, timeVal := range parentEvent.Times {
newTime, err := calculateNextDate(timeVal.Time(), interval, unit)
if err != nil {
return err
}
newTimes[i] = primitive.NewDateTimeFromTime(newTime)
}
}

// Create the new event as a copy of the parent
numResponses := 0
newEvent := models.Event{
Id:                       primitive.NewObjectID(),
OwnerId:                  parentEvent.OwnerId,
Name:                     parentEvent.Name,
Description:              parentEvent.Description,
Location:                 parentEvent.Location,
Duration:                 parentEvent.Duration,
Dates:                    newDates,
HasSpecificTimes:         parentEvent.HasSpecificTimes,
Times:                    newTimes,
IsSignUpForm:             parentEvent.IsSignUpForm,
SignUpBlocks:             parentEvent.SignUpBlocks,
StartOnMonday:            parentEvent.StartOnMonday,
NotificationsEnabled:     parentEvent.NotificationsEnabled,
BlindAvailabilityEnabled: parentEvent.BlindAvailabilityEnabled,
DaysOnly:                 parentEvent.DaysOnly,
SendEmailAfterXResponses: parentEvent.SendEmailAfterXResponses,
CollectEmails:            parentEvent.CollectEmails,
TimeIncrement:            parentEvent.TimeIncrement,
Type:                     parentEvent.Type,
SignUpResponses:          make(map[string]*models.SignUpResponse),
NumResponses:             &numResponses,
ParentEventId:            parentEvent.Id,
}

// Generate short id
shortId := db.GenerateShortEventId(newEvent.Id)
newEvent.ShortId = &shortId

// Send reminder emails if parent event had remindees
if parentEvent.Remindees != nil && len(*parentEvent.Remindees) > 0 {
// Get owner name
var ownerName string
if parentEvent.OwnerId != primitive.NilObjectID {
owner := db.GetUserById(parentEvent.OwnerId.Hex())
if owner != nil {
ownerName = owner.FirstName
} else {
ownerName = "Somebody"
}
} else {
ownerName = "Somebody"
}

// Schedule email reminders for each remindee
remindees := make([]models.Remindee, 0)
for _, remindee := range *parentEvent.Remindees {
taskIds := gcloud.CreateEmailTask(remindee.Email, ownerName, newEvent.Name, newEvent.GetId())
remindees = append(remindees, models.Remindee{
Email:     remindee.Email,
TaskIds:   taskIds,
Responded: utils.FalsePtr(),
})
}

newEvent.Remindees = &remindees
}

// Insert the new event
_, err := db.EventsCollection.InsertOne(context.Background(), newEvent)
if err != nil {
return fmt.Errorf("error inserting new event: %v", err)
}

logger.StdOut.Printf("Created new event %s from parent %s\n", newEvent.Id.Hex(), parentEvent.Id.Hex())

// Update parent event with next occurrence date
// Calculate when the next event after this one should be created
latestDate := newDates[0].Time()
for _, date := range newDates {
if date.Time().After(latestDate) {
latestDate = date.Time()
}
}

nextEventDate, err := calculateNextDate(latestDate, interval, unit)
if err != nil {
return fmt.Errorf("error calculating next occurrence date: %v", err)
}

// Subtract advance days to get when to create the event
advanceDays := utils.Coalesce(parentEvent.RecurrenceAdvanceDays)
createDate := nextEventDate.AddDate(0, 0, -advanceDays)
nextOccurrenceDateTime := primitive.NewDateTimeFromTime(createDate)

_, err = db.EventsCollection.UpdateOne(
context.Background(),
bson.M{"_id": parentEvent.Id},
bson.M{
"$set": bson.M{
"nextOccurrenceDate": nextOccurrenceDateTime,
},
},
)

if err != nil {
return fmt.Errorf("error updating parent event: %v", err)
}

return nil
}
