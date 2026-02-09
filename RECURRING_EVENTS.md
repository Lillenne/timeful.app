# Recurring Events Feature

This document describes the recurring events feature and how to test it.

## Overview

The recurring events feature allows authenticated users to automatically create new event instances at regular intervals. This is ideal for weekly meetings, monthly check-ins, or any regularly occurring event.

## Features

- **Authentication Required**: Only signed-in users can create recurring events
- **Flexible Intervals**: Repeat every N days, weeks, or months
- **Advance Creation**: Create new instances X days before the event date
- **Stop Anytime**: Event owners can stop future occurrences at any time
- **Email Reminders**: Each new instance sends reminders to the same recipients

## User Workflow

### Creating a Recurring Event

1. **Sign in** to your account (required for recurring events)

2. **Create a new event** with your desired dates and times

3. **Expand "Recurring event" section** below the email reminders

4. **Enable recurring** by checking "Make this a recurring event"

5. **Configure recurrence**:
   - **Repeat every**: Enter a number and select unit (days/weeks/months)
     - Example: "1 weeks" for weekly events
     - Example: "2 weeks" for bi-weekly events
     - Example: "1 months" for monthly events
   
   - **Create event**: Enter how many days in advance to create the next occurrence
     - Example: "2 days" means the next event is created 2 days before it should happen
     - This gives participants time to respond before the event

6. **Create the event** - The first instance is created immediately

7. **System behavior**:
   - The scheduler runs every hour
   - When it's time (based on advance days), a new event instance is automatically created
   - All details are copied: name, description, location, times, remindees
   - Email reminders are sent to remindees as usual

### Example: Weekly Team Meeting

**Scenario**: Create a weekly team meeting that repeats every Monday at 10 AM, created 2 days in advance.

**Steps**:
1. Sign in to Timeful
2. Create new event:
   - Name: "Weekly Team Sync"
   - Date: Next Monday
   - Time: 10:00 AM - 11:00 AM
3. Add email reminders for team members
4. Enable recurring:
   - Repeat every: 1 weeks
   - Create event: 2 days in advance
5. Create event

**Result**:
- First meeting appears for next Monday
- On Saturday (2 days before), the following Monday's meeting is automatically created
- This continues indefinitely until you stop it

### Viewing Recurring Event Details

When viewing a recurring event, you'll see:
- **Recurring indicator** with a repeat icon
- **Configuration details**: "Repeats every N days/weeks/months"
- **Advance creation info**: "(created X days in advance)"
- **Status**: Shows if recurring has been stopped

### Stopping Recurring Events

1. **Open the recurring event** (any instance will show the recurring status)

2. **View event details** section

3. **Click "Stop recurring"** button (only visible to event owner)

4. **Confirm** the action in the dialog

5. **Result**:
   - Future instances will no longer be created
   - Existing instances remain unchanged
   - Status updates to "Recurring event creation has been stopped"

## Technical Details

### Database Fields

New fields added to the Event model:
- `isRecurring`: Boolean indicating if event is recurring
- `recurrenceInterval`: Number of units between occurrences (e.g., 1, 2, 3)
- `recurrenceUnit`: Unit of time ("days", "weeks", or "months")
- `recurrenceAdvanceDays`: Days in advance to create next occurrence
- `parentEventId`: Reference to parent event (for instances)
- `recurrenceEnabled`: Whether to continue creating occurrences
- `nextOccurrenceDate`: When to create the next occurrence

### Scheduler

- Runs every hour
- Queries for events where:
  - `isRecurring = true`
  - `recurrenceEnabled = true`
  - `nextOccurrenceDate <= now`
  - Event is not deleted
- Creates new instance with dates shifted by interval
- Updates parent event's `nextOccurrenceDate`

### API Endpoints

**Create Event** (`POST /api/events`):
```json
{
  "name": "Weekly Meeting",
  "duration": 1,
  "dates": ["2024-02-12T10:00:00Z"],
  "type": "specific_dates",
  "isRecurring": true,
  "recurrenceInterval": 1,
  "recurrenceUnit": "weeks",
  "recurrenceAdvanceDays": 2
}
```

**Edit Event** (`PUT /api/events/:eventId`):
```json
{
  "name": "Weekly Meeting",
  "duration": 1,
  "dates": ["2024-02-12T10:00:00Z"],
  "type": "specific_dates",
  "recurrenceEnabled": false
}
```

## Testing Checklist

### Basic Functionality
- [ ] Cannot create recurring event as guest user
- [ ] Can create recurring event as authenticated user
- [ ] Daily recurring event works correctly
- [ ] Weekly recurring event works correctly
- [ ] Monthly recurring event works correctly
- [ ] Advance creation days works correctly

### Event Creation
- [ ] First event is created immediately
- [ ] All details are copied (name, description, location)
- [ ] Times are shifted correctly
- [ ] Email reminders are sent for new instances
- [ ] Parent event ID is set on instances

### Scheduler
- [ ] Scheduler runs on server startup
- [ ] Scheduler creates events at correct time
- [ ] Multiple recurring events are processed
- [ ] Failed creation doesn't stop other events

### Stopping Recurring
- [ ] Stop button only visible to event owner
- [ ] Confirmation dialog appears
- [ ] Future instances stop being created
- [ ] Existing instances remain unchanged
- [ ] Status updates correctly in UI

### Validation
- [ ] Interval must be > 0
- [ ] Unit must be days/weeks/months
- [ ] Advance days must be >= 0
- [ ] Authentication required

### Edge Cases
- [ ] Event with no dates (should not break)
- [ ] Event with specific times (times should shift)
- [ ] Event deleted (should not create instances)
- [ ] Very large interval (e.g., 100 weeks)
- [ ] Zero advance days (creates on same day)

## Troubleshooting

### Recurring event not being created

1. **Check scheduler is running**:
   - Look for "Starting recurring event scheduler..." in server logs
   - Look for "Checking for recurring events to create..." every hour

2. **Check nextOccurrenceDate**:
   - Query the database for your event
   - Verify `nextOccurrenceDate` is in the past

3. **Check recurrenceEnabled**:
   - Ensure it's set to `true`
   - If `false`, recurring has been stopped

4. **Check server logs**:
   - Look for error messages from scheduler
   - Check for "Creating next occurrence for event..."

### UI not showing recurring options

1. **Check authentication**:
   - User must be signed in
   - Recurring section only shown for authenticated users

2. **Check event type**:
   - Recurring only available when creating new events
   - Not available when editing existing events

### Dates not shifting correctly

1. **Check recurrence unit**:
   - Ensure unit is spelled correctly in database
   - Valid values: "days", "weeks", "months"

2. **Check interval**:
   - Must be a positive integer

3. **Check for DST issues**:
   - Times should maintain local time across DST changes

## Future Enhancements

Possible improvements for future versions:

- Allow editing recurring configuration on existing events
- Add end date for recurring events (stop after date)
- Add occurrence count limit (stop after N occurrences)
- Support more complex patterns (every weekday, monthly on 15th, etc.)
- Bulk update all future occurrences
- Link between parent and child instances in UI
- Recurring event analytics

## Support

For issues or questions about recurring events:
1. Check this documentation
2. Review server logs for errors
3. Check database for event configuration
4. Open an issue on GitHub with details
