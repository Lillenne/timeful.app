/**
 * Utility functions for parsing and creating ICS (iCalendar) files
 */

import ICAL from "ical.js"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import timezone from "dayjs/plugin/timezone"

dayjs.extend(utc)
dayjs.extend(timezone)

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
      })
    }

    return events
  } catch (error) {
    console.error("Error parsing ICS file:", error)
    throw new Error("Failed to parse ICS file. Please ensure it's a valid calendar file.")
  }
}

/**
 * Validate and sanitize calendar URL
 * @param {string} url - The URL to validate
 * @returns {string} Validated URL
 * @throws {Error} If URL is invalid or uses an unsupported protocol
 */
function validateCalendarURL(url) {
  try {
    // Convert webcal:// to https://
    let sanitizedUrl = url.trim().replace(/^webcal:\/\//i, "https://")
    
    // Parse the URL
    const urlObj = new URL(sanitizedUrl)
    
    // Only allow https and http protocols
    if (!['https:', 'http:'].includes(urlObj.protocol)) {
      throw new Error("Only https:// and http:// URLs are supported")
    }
    
    // Don't allow localhost or private IP ranges for security
    const hostname = urlObj.hostname.toLowerCase()
    if (hostname === 'localhost' || 
        hostname === '127.0.0.1' ||
        hostname.startsWith('192.168.') ||
        hostname.startsWith('10.') ||
        hostname.match(/^172\.(1[6-9]|2[0-9]|3[01])\./)) {
      throw new Error("Private or local URLs are not allowed for security reasons")
    }
    
    return sanitizedUrl
  } catch (error) {
    throw new Error(`Invalid calendar URL: ${error.message}`)
  }
}

/**
 * Fetch ICS content from a URL
 * @param {string} url - The URL to fetch the ICS file from
 * @returns {Promise<string>} The ICS file content
 */
export async function fetchICSFromURL(url) {
  try {
    // Validate URL first
    const validatedUrl = validateCalendarURL(url)
    
    const response = await fetch(validatedUrl)
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
 * @param {Date} eventDetails.startDate - Event start date (in UTC)
 * @param {Date} eventDetails.endDate - Event end date (in UTC)
 * @param {string} eventDetails.description - Event description
 * @param {string} eventDetails.location - Event location
 * @param {Array<string>} eventDetails.attendees - Array of attendee emails
 * @param {string} eventDetails.timezone - IANA timezone identifier (e.g., "America/New_York")
 * @returns {string} ICS file content
 */
export function createICSFile(eventDetails) {
  const { title, startDate, endDate, description, location, attendees = [], timezone } = eventDetails

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
  
  // Set start and end times with timezone
  if (timezone) {
    // When timezone is provided, we need to create local time values (not UTC)
    // The incoming dates are UTC times that represent the correct moment
    // We convert them to local time in the specified timezone for the ICS file
    
    const localStart = dayjs(startDate).tz(timezone)
    const localEnd = dayjs(endDate).tz(timezone)
    
    // Create ICAL.Time as local time (not UTC)
    const startTime = new ICAL.Time({
      year: localStart.year(),
      month: localStart.month() + 1, // dayjs months are 0-indexed, ICAL months are 1-indexed
      day: localStart.date(),
      hour: localStart.hour(),
      minute: localStart.minute(),
      second: localStart.second(),
      isDate: false
    })
    
    const endTime = new ICAL.Time({
      year: localEnd.year(),
      month: localEnd.month() + 1,
      day: localEnd.date(),
      hour: localEnd.hour(),
      minute: localEnd.minute(),
      second: localEnd.second(),
      isDate: false
    })
    
    event.startDate = startTime
    event.endDate = endTime
    
    // Set the TZID parameter on the date-time properties
    const dtstart = vevent.getFirstProperty("dtstart")
    const dtend = vevent.getFirstProperty("dtend")
    
    if (dtstart) {
      dtstart.setParameter("tzid", timezone)
    }
    if (dtend) {
      dtend.setParameter("tzid", timezone)
    }
  } else {
    // No timezone provided, use UTC times
    event.startDate = ICAL.Time.fromJSDate(startDate, false)
    event.endDate = ICAL.Time.fromJSDate(endDate, false)
  }

  // Generate UID using crypto API if available, fallback to timestamp + random
  let uid
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    uid = `timeful-${crypto.randomUUID()}@timeful.app`
  } else {
    uid = `timeful-${Date.now()}-${Math.random().toString(36).substring(2, 11)}@timeful.app`
  }
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
