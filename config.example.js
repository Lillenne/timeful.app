window.__TIMEFUL_CONFIG__ = {
  // Google OAuth Client ID (required for login/calendar features)
  // Must match the CLIENT_ID from your .env file
  // Get this from: https://console.cloud.google.com/apis/credentials
  googleClientId: '',

  // Microsoft OAuth Client ID (optional - required for Outlook calendar integration)
  // Must match the MICROSOFT_CLIENT_ID from your .env file
  // Get this from: https://portal.azure.com/ -> Azure Active Directory -> App registrations
  // Leave empty to disable Outlook calendar integration
  microsoftClientId: '',

  // PostHog analytics API key (optional)
  // Leave empty to disable PostHog analytics
  // Get this from: https://posthog.com/
  posthogApiKey: '',

  // Disable all analytics (Google Tag Manager, PostHog, etc.)
  // Set to true to completely disable analytics - no tracking at all
  // Useful for privacy-focused deployments
  disableAnalytics: false,

  // Mapbox API key (optional - for address autocomplete)
  // Leave empty to disable address autocomplete feature
  // Get this from: https://www.mapbox.com/
  // Free tier includes 100,000 requests per month
  mapboxApiKey: '',
}
