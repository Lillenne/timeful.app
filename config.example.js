window.__TIMEFUL_CONFIG__ = {
  // Google OAuth Client ID (required for login/calendar features)
  // Must match the CLIENT_ID from your .env file
  // Get this from: https://console.cloud.google.com/apis/credentials
  googleClientId: '',

  // PostHog analytics API key (optional)
  // Leave empty to disable PostHog analytics
  // Get this from: https://posthog.com/
  posthogApiKey: '',

  // Disable all analytics (Google Tag Manager, PostHog, etc.)
  // Set to true to completely disable analytics - no tracking at all
  // Useful for privacy-focused deployments
  disableAnalytics: false,
}
