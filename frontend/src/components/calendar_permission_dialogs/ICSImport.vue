<template>
  <div class="tw-flex tw-flex-col tw-gap-6">
    <div class="tw-flex tw-flex-col tw-gap-3">
      <div class="tw-text-md tw-font-medium">Import from ICS file or URL</div>
      <div class="tw-text-sm tw-text-very-dark-gray">
        Import your availability from an ICS calendar file or a calendar URL
        (such as a Google Calendar public link or webcal URL).
      </div>
    </div>

    <v-tabs v-model="activeTab" class="tw-mb-2">
      <v-tab>File Upload</v-tab>
      <v-tab>Calendar URL</v-tab>
    </v-tabs>

    <v-tabs-items v-model="activeTab">
      <!-- File Upload Tab -->
      <v-tab-item>
        <div class="tw-flex tw-flex-col tw-gap-3">
          <div class="tw-text-sm tw-text-very-dark-gray">
            Upload an .ics file from your computer. Events marked as busy will
            block your availability.
          </div>
          <v-file-input
            v-model="file"
            accept=".ics"
            label="Select ICS file"
            placeholder="Choose a file"
            prepend-icon="mdi-calendar-import"
            outlined
            dense
            hide-details
            @change="handleFileSelect"
          />
        </div>
      </v-tab-item>

      <!-- URL Tab -->
      <v-tab-item>
        <div class="tw-flex tw-flex-col tw-gap-3">
          <div class="tw-text-sm tw-text-very-dark-gray">
            Enter a calendar URL (webcal://, https://, or http://). This is
            useful for public calendar feeds.
          </div>
          <v-text-field
            v-model="url"
            label="Calendar URL"
            placeholder="https://calendar.google.com/calendar/ical/..."
            outlined
            dense
            hide-details
          />
        </div>
      </v-tab-item>
    </v-tabs-items>

    <!-- Error message -->
    <v-alert v-if="errorMessage" type="error" dense text>
      {{ errorMessage }}
    </v-alert>

    <!-- Success message -->
    <v-alert v-if="successMessage" type="success" dense text>
      {{ successMessage }}
    </v-alert>

    <!-- Action buttons -->
    <div class="tw-flex tw-items-center tw-gap-2">
      <v-btn text class="tw-grow" @click="$emit('back')">Back</v-btn>
      <v-btn
        :disabled="!canImport"
        color="primary"
        class="tw-grow"
        :loading="loading"
        @click="importCalendar"
      >
        Import
      </v-btn>
    </div>
  </div>
</template>

<script>
import { parseICSFile, fetchICSFromURL } from "@/utils/ics_utils"

export default {
  name: "ICSImport",

  data() {
    return {
      activeTab: 0,
      file: null,
      url: "",
      loading: false,
      errorMessage: "",
      successMessage: "",
    }
  },

  computed: {
    canImport() {
      return (this.activeTab === 0 && this.file) || (this.activeTab === 1 && this.url)
    },
  },

  methods: {
    handleFileSelect() {
      this.errorMessage = ""
      this.successMessage = ""
    },

    async importCalendar() {
      this.loading = true
      this.errorMessage = ""
      this.successMessage = ""

      try {
        let icsContent = ""

        if (this.activeTab === 0) {
          // File upload
          icsContent = await this.readFile(this.file)
        } else {
          // URL
          // Convert webcal:// to https://
          let calendarUrl = this.url.replace(/^webcal:\/\//i, "https://")
          
          // For security, we need to proxy the request through our backend
          // to avoid CORS issues
          try {
            icsContent = await fetchICSFromURL(calendarUrl)
          } catch (err) {
            // If direct fetch fails, show a helpful message
            this.errorMessage =
              "Unable to fetch calendar from URL. Please ensure the URL is public and accessible, or try downloading the ICS file and uploading it instead."
            return
          }
        }

        // Parse the ICS content
        const events = parseICSFile(icsContent)

        if (events.length === 0) {
          this.errorMessage =
            "No busy events found in the calendar. Please check the file or URL."
          return
        }

        this.successMessage = `Successfully imported ${events.length} event(s)!`

        // Emit the events to the parent component
        this.$emit("eventsImported", events)

        this.$posthog.capture("ICS Calendar Imported", {
          source: this.activeTab === 0 ? "file" : "url",
          eventCount: events.length,
        })
      } catch (error) {
        console.error("Error importing calendar:", error)
        this.errorMessage = error.message || "Failed to import calendar. Please try again."
      } finally {
        this.loading = false
      }
    },

    readFile(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = (e) => resolve(e.target.result)
        reader.onerror = (e) => reject(new Error("Failed to read file"))
        reader.readAsText(file)
      })
    },
  },

  watch: {
    activeTab() {
      this.errorMessage = ""
      this.successMessage = ""
    },
  },
}
</script>
