import { calendarTypes } from "@/constants"
import store from "@/store"

/** Redirects user to the correct google sign in page */
export const signInGoogle = ({
  state = {},
  selectAccount = false,
  requestCalendarPermission = false,
  requestContactsPermission = false,
  loginHint = "",
}) => {
  const clientId = window.__TIMEFUL_CONFIG__?.googleClientId || "523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com"
  const redirectUri = `${window.location.origin}/auth`

  let scope = "openid email profile "
  if (requestCalendarPermission) {
    scope +=
      "https://www.googleapis.com/auth/calendar.calendarlist.readonly https://www.googleapis.com/auth/calendar.events.readonly "
  }
  if (requestContactsPermission) {
    scope +=
      "https://www.googleapis.com/auth/contacts.readonly https://www.googleapis.com/auth/directory.readonly "
  }
  scope = encodeURIComponent(scope)

  let stateString = ""
  if (!state) state = {}
  state.calendarType = calendarTypes.GOOGLE
  state = encodeURIComponent(JSON.stringify(state))
  stateString = `&state=${state}`

  let promptString = ""
  if (selectAccount) {
    promptString = "&prompt=select_account+consent"
  } else {
    promptString = "&prompt=consent"
    if (loginHint.length > 0) {
      promptString += `&login_hint=${loginHint}`
    } else if (store.state.authUser) {
      promptString += `&login_hint=${store.state.authUser.email}`
    }
  }

  const url = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&access_type=offline${promptString}${stateString}&include_granted_scopes=true`
  window.location.href = url
}

export const signInOutlook = ({
  state = {},
  requestCalendarPermission = false,
}) => {
  const clientId = window.__TIMEFUL_CONFIG__?.microsoftClientId
  
  // Check if Microsoft client ID is configured
  if (!clientId) {
    console.error("Microsoft client ID not configured. Please set microsoftClientId in config.js")
    alert("Outlook calendar integration is not configured. Please contact your administrator.")
    return
  }
  
  const tenant = "common"
  const redirectUri = encodeURIComponent(`${window.location.origin}/auth`)

  let scope = "offline_access User.Read"
  if (requestCalendarPermission) {
    scope += " Calendars.Read"
  }
  scope = encodeURIComponent(scope)

  let stateString = ""
  if (!state) state = {}
  state.calendarType = calendarTypes.OUTLOOK
  state.scope = scope
  state = encodeURIComponent(JSON.stringify(state))
  stateString = `&state=${state}`

  const url = `https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=${clientId}&response_type=code&redirect_uri=${redirectUri}&response_mode=query&scope=${scope}${stateString}`
  window.location.href = url
}
