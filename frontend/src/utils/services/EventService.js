import { post } from "../fetch_utils"

export const archiveEvent = (eventId, archive) => {
  return post(`/events/${eventId}/archive`, {
    archive: archive,
  })
}

export const scheduleEvent = (eventId, scheduledEvent) => {
  return post(`/events/${eventId}/schedule-event`, {
    scheduledEvent: scheduledEvent,
  })
}
