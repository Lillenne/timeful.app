<template>
  <div class="location-autocomplete">
    <v-text-field
      :value="value"
      @input="handleInput"
      :placeholder="placeholder"
      :hide-details="hideDetails"
      :solo="solo"
      :counter="counter"
      :maxlength="maxlength"
      ref="textField"
      @focus="onFocus"
      @blur="onBlur"
    >
      <template v-slot:append>
        <v-icon v-if="loading" small>mdi-loading mdi-spin</v-icon>
      </template>
    </v-text-field>
    
    <!-- Autocomplete suggestions -->
    <v-list
      v-if="showSuggestions && suggestions.length > 0"
      class="tw-absolute tw-z-50 tw-mt-1 tw-max-h-64 tw-w-full tw-overflow-auto tw-rounded tw-border tw-border-gray tw-bg-white tw-shadow-lg"
    >
      <v-list-item
        v-for="(suggestion, index) in suggestions"
        :key="index"
        @click="selectSuggestion(suggestion)"
        class="tw-cursor-pointer hover:tw-bg-gray-100"
      >
        <v-list-item-content>
          <v-list-item-title>{{ suggestion.place_name || suggestion.full_address }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </div>
</template>

<script>
export default {
  name: "LocationAutocomplete",

  props: {
    value: {
      type: String,
      default: "",
    },
    placeholder: {
      type: String,
      default: "Add a location (optional)...",
    },
    hideDetails: {
      type: String,
      default: "auto",
    },
    solo: {
      type: Boolean,
      default: true,
    },
    counter: {
      type: String,
      default: "500",
    },
    maxlength: {
      type: String,
      default: "500",
    },
  },

  data() {
    return {
      suggestions: [],
      loading: false,
      showSuggestions: false,
      debounceTimer: null,
      mapboxApiKey: null,
    }
  },

  mounted() {
    // Get Mapbox API key from global config
    if (window.__TIMEFUL_CONFIG__ && window.__TIMEFUL_CONFIG__.mapboxApiKey) {
      this.mapboxApiKey = window.__TIMEFUL_CONFIG__.mapboxApiKey
    }
  },

  methods: {
    handleInput(value) {
      this.$emit("input", value)
      
      // Only fetch suggestions if Mapbox API key is available
      if (!this.mapboxApiKey || !value || value.length < 3) {
        this.suggestions = []
        this.showSuggestions = false
        return
      }

      // Debounce the API call
      clearTimeout(this.debounceTimer)
      this.debounceTimer = setTimeout(() => {
        this.fetchSuggestions(value)
      }, 300)
    },

    async fetchSuggestions(query) {
      if (!this.mapboxApiKey) return

      this.loading = true
      try {
        // Use Mapbox Search Box API
        const response = await fetch(
          `https://api.mapbox.com/search/searchbox/v1/suggest?q=${encodeURIComponent(
            query
          )}&access_token=${this.mapboxApiKey}&session_token=${this.getSessionToken()}&language=en&limit=5`
        )

        if (response.ok) {
          const data = await response.json()
          this.suggestions = data.suggestions || []
          this.showSuggestions = this.suggestions.length > 0
        }
      } catch (error) {
        console.error("Error fetching location suggestions:", error)
        this.suggestions = []
        this.showSuggestions = false
      } finally {
        this.loading = false
      }
    },

    selectSuggestion(suggestion) {
      const locationText = suggestion.place_name || suggestion.full_address || suggestion.name
      this.$emit("input", locationText)
      this.suggestions = []
      this.showSuggestions = false
      this.$refs.textField.blur()
    },

    onFocus() {
      // Show suggestions if we have any
      if (this.suggestions.length > 0) {
        this.showSuggestions = true
      }
    },

    onBlur() {
      // Delay hiding to allow click on suggestion
      setTimeout(() => {
        this.showSuggestions = false
      }, 200)
    },

    getSessionToken() {
      // Generate or retrieve a session token for the autocomplete session
      // This helps Mapbox track usage and optimize suggestions
      if (!this._sessionToken) {
        this._sessionToken = `timeful-${Date.now()}-${Math.random().toString(36).substring(2, 11)}`
      }
      return this._sessionToken
    },
  },

  beforeDestroy() {
    clearTimeout(this.debounceTimer)
  },
}
</script>

<style scoped>
.location-autocomplete {
  position: relative;
  width: 100%;
}
</style>
