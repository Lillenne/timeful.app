# Listmonk Email Templates Setup Guide

This guide provides instructions and examples for setting up email templates in Listmonk for use with Timeful.

## Overview

Timeful uses Listmonk for sending various transactional emails including:
- Event reminders
- Availability group invitations
- Response notifications
- Group attendee notifications

## Required Templates

You need to create 9 email templates in Listmonk. Below are the details for each template.

### Template 1: Everyone Responded Notification (ID: 8)

**Purpose**: Sent to event owner when all attendees have responded

**Template Variables**:
- `{{ .ownerName }}` - First name of the event owner
- `{{ .eventName }}` - Name of the event
- `{{ .eventUrl }}` - URL to view the event

**Subject Line**: `Everyone has responded to {{ .eventName }}! 🎉`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Everyone Responded</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">Great news, {{ .ownerName }}! 🎉</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Everyone has responded to your event <strong>{{ .eventName }}</strong>!
        </p>
        <p style="font-size: 16px; margin-bottom: 25px;">
            You can now view all the responses and find the best time for your event.
        </p>
        <div style="text-align: center;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Responses</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 2: Availability Group Invite (ID: 9)

**Purpose**: Sent when someone is invited to an availability group

**Template Variables**:
- `{{ .ownerName }}` - First name of the person who invited them
- `{{ .groupName }}` - Name of the availability group
- `{{ .groupUrl }}` - URL to view the group

**Subject Line**: `{{ .ownerName }} invited you to {{ .groupName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Group Invitation</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">You've been invited! 📅</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            <strong>{{ .ownerName }}</strong> has invited you to join the availability group <strong>{{ .groupName }}</strong>.
        </p>
        <p style="font-size: 16px; margin-bottom: 25px;">
            Share your availability to help find times that work for everyone.
        </p>
        <div style="text-align: center;">
            <a href="{{ .groupUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Group</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 3: Someone Responded (Events - ID: 10)

**Purpose**: Sent to event owner when someone responds to a non-group event

**Template Variables**:
- `{{ .ownerName }}` - First name of the event owner
- `{{ .eventName }}` - Name of the event
- `{{ .respondentName }}` - Name of the person who responded
- `{{ .eventUrl }}` - URL to view the event

**Subject Line**: `{{ .respondentName }} responded to {{ .eventName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Response</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">New response received! 📋</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Hi {{ .ownerName }},
        </p>
        <p style="font-size: 16px; margin-bottom: 20px;">
            <strong>{{ .respondentName }}</strong> has responded to your event <strong>{{ .eventName }}</strong>.
        </p>
        <div style="text-align: center; margin-top: 25px;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Response</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 4: Added as Attendee (ID: 11)

**Purpose**: Sent when someone is added as an attendee to an availability group

**Template Variables**:
- `{{ .ownerName }}` - First name of the person who added them
- `{{ .groupName }}` - Name of the availability group
- `{{ .groupUrl }}` - URL to view the group

**Subject Line**: `{{ .ownerName }} added you to {{ .groupName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Added as Attendee</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">You've been added! 👥</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            <strong>{{ .ownerName }}</strong> has added you to the availability group <strong>{{ .groupName }}</strong>.
        </p>
        <p style="font-size: 16px; margin-bottom: 25px;">
            You can now view and share availability with other members of this group.
        </p>
        <div style="text-align: center;">
            <a href="{{ .groupUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Group</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 5: Someone Responded (Groups - ID: 13)

**Purpose**: Sent to group owner when someone responds to an availability group

**Template Variables**:
- `{{ .ownerName }}` - First name of the group owner
- `{{ .groupName }}` - Name of the availability group
- `{{ .respondentName }}` - Name of the person who responded
- `{{ .groupUrl }}` - URL to view the group

**Subject Line**: `{{ .respondentName }} responded to {{ .groupName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Response</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">New availability shared! 📋</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Hi {{ .ownerName }},
        </p>
        <p style="font-size: 16px; margin-bottom: 20px;">
            <strong>{{ .respondentName }}</strong> has shared their availability in the group <strong>{{ .groupName }}</strong>.
        </p>
        <div style="text-align: center; margin-top: 25px;">
            <a href="{{ .groupUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Availability</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 6: X Responses Received (ID: 14)

**Purpose**: Sent to event owner when a threshold number of responses is reached

**Template Variables**:
- `{{ .ownerName }}` - First name of the event owner
- `{{ .eventName }}` - Name of the event
- `{{ .numResponses }}` - Number of responses received
- `{{ .eventUrl }}` - URL to view the event

**Subject Line**: `{{ .numResponses }} people have responded to {{ .eventName }}!`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Response Milestone</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">Milestone reached! 🎯</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Hi {{ .ownerName }},
        </p>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Great news! You've received <strong>{{ .numResponses }} responses</strong> to your event <strong>{{ .eventName }}</strong>.
        </p>
        <p style="font-size: 16px; margin-bottom: 25px;">
            Check out the responses to find the best time for your event.
        </p>
        <div style="text-align: center;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">View Responses</a>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 7: Initial Reminder Email (Configurable ID)

**Purpose**: Sent immediately when someone is added to the reminder list

**Template Variables**:
- `{{ .ownerName }}` - First name of the event owner
- `{{ .eventName }}` - Name of the event
- `{{ .eventUrl }}` - URL to respond to the event
- `{{ .finishedUrl }}` - URL to mark as responded

**Subject Line**: `{{ .ownerName }} wants to know when you're free for {{ .eventName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Reminder</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 10px; padding: 30px; margin-bottom: 20px;">
        <h1 style="color: #2563eb; margin-top: 0;">When are you free? 📅</h1>
        <p style="font-size: 16px; margin-bottom: 20px;">
            Hi there! 👋
        </p>
        <p style="font-size: 16px; margin-bottom: 20px;">
            <strong>{{ .ownerName }}</strong> is trying to schedule <strong>{{ .eventName }}</strong> and would like to know when you're available.
        </p>
        <p style="font-size: 16px; margin-bottom: 25px;">
            Please take a moment to share your availability - it only takes a minute!
        </p>
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #2563eb; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">Share Your Availability</a>
        </div>
        <div style="background-color: #fff; border: 1px solid #e5e7eb; border-radius: 6px; padding: 15px; margin-top: 20px;">
            <p style="font-size: 14px; color: #6b7280; margin: 0;">
                Already responded? <a href="{{ .finishedUrl }}" style="color: #2563eb; text-decoration: none;">Let us know</a> to stop receiving reminders.
            </p>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 8: Second Reminder Email (Configurable ID)

**Purpose**: Sent 24 hours after initial reminder if no response

**Template Variables**: Same as Initial Reminder

**Subject Line**: `Reminder: {{ .ownerName }} needs your availability for {{ .eventName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reminder</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 30px; margin-bottom: 20px; border-radius: 6px;">
        <h1 style="color: #92400e; margin-top: 0;">Friendly reminder! ⏰</h1>
        <p style="font-size: 16px; margin-bottom: 20px; color: #78350f;">
            Just a quick reminder that <strong>{{ .ownerName }}</strong> is still waiting for your availability for <strong>{{ .eventName }}</strong>.
        </p>
        <p style="font-size: 16px; margin-bottom: 25px; color: #78350f;">
            It only takes a minute to respond!
        </p>
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #f59e0b; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">Share Your Availability</a>
        </div>
        <div style="background-color: #fff; border: 1px solid #fbbf24; border-radius: 6px; padding: 15px; margin-top: 20px;">
            <p style="font-size: 14px; color: #92400e; margin: 0;">
                Already responded? <a href="{{ .finishedUrl }}" style="color: #d97706; text-decoration: none;">Let us know</a> to stop receiving reminders.
            </p>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
    </div>
</body>
</html>
```

### Template 9: Final Reminder Email (Configurable ID)

**Purpose**: Sent 72 hours after initial reminder if no response (last reminder)

**Template Variables**: Same as Initial Reminder

**Subject Line**: `Last reminder: Please share your availability for {{ .eventName }}`

**HTML Body**:
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Final Reminder</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #fee2e2; border-left: 4px solid #ef4444; padding: 30px; margin-bottom: 20px; border-radius: 6px;">
        <h1 style="color: #991b1b; margin-top: 0;">Final reminder! ⚠️</h1>
        <p style="font-size: 16px; margin-bottom: 20px; color: #7f1d1d;">
            This is the last reminder from <strong>{{ .ownerName }}</strong> about <strong>{{ .eventName }}</strong>.
        </p>
        <p style="font-size: 16px; margin-bottom: 20px; color: #7f1d1d;">
            Your input is important! Please take a moment to share your availability so {{ .ownerName }} can find a time that works for everyone.
        </p>
        <p style="font-size: 14px; margin-bottom: 25px; color: #991b1b; font-style: italic;">
            This is the final reminder you'll receive.
        </p>
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="{{ .eventUrl }}" style="display: inline-block; background-color: #ef4444; color: white; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-weight: 600; font-size: 16px;">Share Your Availability Now</a>
        </div>
        <div style="background-color: #fff; border: 1px solid #fca5a5; border-radius: 6px; padding: 15px; margin-top: 20px;">
            <p style="font-size: 14px; color: #991b1b; margin: 0;">
                Already responded? <a href="{{ .finishedUrl }}" style="color: #dc2626; text-decoration: none;">Let us know</a> to confirm.
            </p>
        </div>
    </div>
    <div style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
        <p>Timeful - Simple scheduling for teams</p>
        <p style="margin-top: 10px;">You won't receive any more reminders about this event.</p>
    </div>
</body>
</html>
```

## Setup Instructions

### Step 1: Access Listmonk Admin

1. Open your browser and navigate to your Listmonk instance (e.g., http://localhost:9000)
2. Log in with your admin credentials

### Step 2: Create Templates

For each template above:

1. Go to **Campaigns** → **Templates** in the Listmonk admin interface
2. Click **"Create New"**
3. Enter the template name (e.g., "Everyone Responded Notification")
4. Set the **Subject** as shown above
5. Paste the **HTML Body** into the template editor
6. Click **Save**
7. **Note the Template ID** - you'll see it in the URL (e.g., `/campaigns/templates/8`)

### Step 3: Configure Template IDs

After creating all templates, update your `.env` file with the template IDs:

```bash
# Fixed template IDs (must match the code expectations)
# Create these templates with these exact IDs or update the code
# Template 8: Everyone Responded Notification
# Template 9: Availability Group Invite
# Template 10: Someone Responded (Events)
# Template 11: Added as Attendee
# Template 13: Someone Responded (Groups)
# Template 14: X Responses Received

# Reminder email template IDs (configurable)
LISTMONK_INITIAL_EMAIL_REMINDER_ID=1
LISTMONK_SECOND_EMAIL_REMINDER_ID=2
LISTMONK_FINAL_EMAIL_REMINDER_ID=3
```

**Note**: The fixed template IDs (8, 9, 10, 11, 13, 14) are hardcoded in the application. You have two options:

1. **Create templates with these exact IDs** by creating and deleting templates until you get the right IDs
2. **Update the code** to use configurable template IDs (recommended for self-hosted deployments)

### Step 4: Configure SMTP

Listmonk needs SMTP settings to send emails:

1. Go to **Settings** → **SMTP**
2. Add your SMTP server details
3. Test the connection
4. Enable the SMTP server

### Step 5: Test Templates

Test each template by sending a test email:

1. Go to **Campaigns** → **Templates**
2. Click on a template
3. Click **"Send Test"**
4. Enter your email address
5. Provide test values for template variables
6. Send and verify the email looks correct

## Template Variable Reference

| Variable | Description | Used In |
|----------|-------------|---------|
| `ownerName` | First name of event/group owner | All templates |
| `eventName` | Name of the event | Event templates |
| `groupName` | Name of the availability group | Group templates |
| `eventUrl` | URL to view/respond to event | Most templates |
| `groupUrl` | URL to view availability group | Group templates |
| `respondentName` | Name of person who responded | Response notifications |
| `numResponses` | Number of responses received | Milestone notification |
| `finishedUrl` | URL to mark as responded | Reminder emails |

## Troubleshooting

### Emails Not Sending

- Check SMTP configuration in Listmonk settings
- Verify template IDs in `.env` match your created templates
- Check Listmonk logs: `docker compose logs -f listmonk`
- Verify `LISTMONK_ENABLED` is not set to "false"

### Template Not Found Errors

- Ensure template IDs in `.env` match the actual template IDs in Listmonk
- Check backend logs: `docker compose logs -f backend`

### Variables Not Rendering

- Ensure template uses correct variable syntax: `{{ .variableName }}`
- Variable names are case-sensitive
- Check that the application is passing the correct data to the template

## Customization Tips

- Use your brand colors in the templates
- Add your logo by including an `<img>` tag with a public URL
- Modify the text and tone to match your organization's voice
- Test on multiple email clients (Gmail, Outlook, etc.)
- Keep mobile users in mind - the templates use responsive design principles

## Advanced: Making Template IDs Configurable

For better self-hosted flexibility, you may want to make all template IDs configurable via environment variables. This requires modifying the Go code in `server/routes/events.go` to read template IDs from environment variables instead of using hardcoded values.

Example changes needed:
```go
// Instead of:
someoneRespondedEmailId := 10

// Use:
someoneRespondedEmailId, _ := strconv.Atoi(os.Getenv("LISTMONK_SOMEONE_RESPONDED_EMAIL_ID"))
```

This would require adding these to `.env`:
```bash
LISTMONK_EVERYONE_RESPONDED_EMAIL_ID=8
LISTMONK_GROUP_INVITE_EMAIL_ID=9
LISTMONK_SOMEONE_RESPONDED_EVENT_EMAIL_ID=10
LISTMONK_ADDED_ATTENDEE_EMAIL_ID=11
LISTMONK_SOMEONE_RESPONDED_GROUP_EMAIL_ID=13
LISTMONK_X_RESPONSES_EMAIL_ID=14
```
