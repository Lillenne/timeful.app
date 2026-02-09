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

  // Enable advertising (Google AdSense)
  // Set to true to enable ads - disabled by default for privacy
  // When false, advertising scripts will not load even if user consents
  enableAdvertising: false,

  // Enable Google Fonts
  // Set to true to load DM Sans font from Google Fonts CDN
  // When false, uses system font stack for better privacy and performance
  enableGoogleFonts: false,

  // Mapbox API key (optional - for address autocomplete)
  // Leave empty to disable address autocomplete feature
  // Get this from: https://www.mapbox.com/
  // Free tier includes 100,000 requests per month
  mapboxApiKey: '',

  // Blog configuration (optional)
  // Set blogUrl to the URL of your blog
  // Set blogButtonText to customize the button text (default: 'Blog')
  // Set blogEnabled to false to hide the blog button
  blogUrl: 'https://schej-blog.vercel.app/blog/',
  blogButtonText: 'Blog',
  blogEnabled: true,
}
