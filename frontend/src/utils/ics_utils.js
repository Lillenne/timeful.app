/**
 * Utility functions for parsing and creating ICS (iCalendar) files
 */

import ICAL from "ical.js"

/**
 * Parse ICS file content and extract calendar events
 * @param {string} icsContent - The raw ICS file content
 * @returns {Array} Array of calendar events with startDate and endDate
 */
export function parseICSFile(icsContent) {
  try {
    const jcalData = ICAL.parse(icsContent)
    const comp = new ICAL.Component(jcalData)
    const vevents = comp.getAllSubcomponents("vevent")

    const events = []
    for (const vevent of vevents) {
      const event = new ICAL.Event(vevent)
      
      // Skip all-day events as they don't block availability
      if (event.startDate.isDate || event.endDate.isDate) {
        continue
      }

      // Check if the event shows as free/transparent
      const transp = vevent.getFirstPropertyValue("transp")
      if (transp === "TRANSPARENT") {
        continue
      }

      events.push({
        summary: event.summary || "Busy",
        startDate: event.startDate.toJSDate(),
        endDate: event.endDate.toJSDate(),
        isAllDay: event.startDate.isDate,
      })
    }

    return events
  } catch (error) {
    console.error("Error parsing ICS file:", error)
    throw new Error("Failed to parse ICS file. Please ensure it's a valid calendar file.")
  }
}

/**
 * Fetch ICS content from a URL
 * @param {string} url - The URL to fetch the ICS file from
 * @returns {Promise<string>} The ICS file content
 */
export async function fetchICSFromURL(url) {
  try {
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`Failed to fetch calendar: ${response.statusText}`)
    }
    const content = await response.text()
    return content
  } catch (error) {
    console.error("Error fetching ICS from URL:", error)
    throw new Error("Failed to fetch calendar from URL. Please check the URL and try again.")
  }
}

/**
 * Create ICS file content from event details
 * @param {Object} eventDetails - Event details
 * @param {string} eventDetails.title - Event title
 * @param {Date} eventDetails.startDate - Event start date
 * @param {Date} eventDetails.endDate - Event end date
 * @param {string} eventDetails.description - Event description
 * @param {string} eventDetails.location - Event location
 * @param {Array<string>} eventDetails.attendees - Array of attendee emails
 * @returns {string} ICS file content
 */
export function createICSFile(eventDetails) {
  const { title, startDate, endDate, description, location, attendees = [] } = eventDetails

  const comp = new ICAL.Component(["vcalendar", [], []])
  comp.updatePropertyWithValue("prodid", "-//Timeful//Timeful App//EN")
  comp.updatePropertyWithValue("version", "2.0")
  comp.updatePropertyWithValue("calscale", "GREGORIAN")
  comp.updatePropertyWithValue("method", "REQUEST")

  const vevent = new ICAL.Component("vevent")
  const event = new ICAL.Event(vevent)

  // Set event properties
  event.summary = title
  event.description = description || ""
  event.location = location || ""
  
  // Set start and end times
  event.startDate = ICAL.Time.fromJSDate(startDate, false)
  event.endDate = ICAL.Time.fromJSDate(endDate, false)

  // Generate UID
  const uid = `timeful-${Date.now()}-${Math.random().toString(36).substr(2, 9)}@timeful.app`
  vevent.updatePropertyWithValue("uid", uid)

  // Set timestamp
  const now = ICAL.Time.now()
  vevent.updatePropertyWithValue("dtstamp", now)
  vevent.updatePropertyWithValue("created", now)

  // Add organizer and attendees
  if (attendees.length > 0) {
    for (const attendee of attendees) {
      if (attendee) {
        const attendeeProp = vevent.addPropertyWithValue("attendee", `mailto:${attendee}`)
        attendeeProp.setParameter("role", "REQ-PARTICIPANT")
        attendeeProp.setParameter("partstat", "NEEDS-ACTION")
        attendeeProp.setParameter("rsvp", "TRUE")
      }
    }
  }

  comp.addSubcomponent(vevent)

  return comp.toString()
}

/**
 * Download ICS file
 * @param {string} icsContent - The ICS file content
 * @param {string} filename - The filename for the download
 */
export function downloadICSFile(icsContent, filename = "event.ics") {
  const blob = new Blob([icsContent], { type: "text/calendar;charset=utf-8" })
  const link = document.createElement("a")
  link.href = URL.createObjectURL(blob)
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(link.href)
}

/**
 * Create a mailto link with ICS attachment
 * @param {string} icsContent - The ICS file content
 * @param {string} subject - Email subject
 * @param {string} body - Email body
 * @returns {string} Mailto URL
 */
export function createMailtoWithICS(icsContent, subject = "", body = "") {
  // Note: Most modern email clients don't support attachments via mailto:
  // This will open the email client with subject and body
  // The ICS file will need to be attached manually or downloaded separately
  const encodedSubject = encodeURIComponent(subject)
  const encodedBody = encodeURIComponent(body + "\n\n(Please attach the downloaded .ics file)")
  return `mailto:?subject=${encodedSubject}&body=${encodedBody}`
}
