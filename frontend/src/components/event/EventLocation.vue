<template>
  <div class="tw-mt-1 tw-max-w-full sm:tw-mt-2 sm:tw-max-w-[calc(100%-236px)]">
    <div
      v-if="showLocation"
      class="tw-flex tw-w-full tw-cursor-pointer tw-flex-col tw-gap-2 tw-rounded-md tw-border tw-border-light-gray-stroke tw-bg-light-gray tw-p-2 tw-text-xs tw-font-normal tw-text-very-dark-gray hover:tw-bg-[#eeeeee] sm:tw-text-sm"
    >
      <div class="tw-flex tw-items-start tw-gap-2">
        <v-icon small class="tw-mt-0.5">mdi-map-marker</v-icon>
        <div class="tw-grow">
          <div class="tw-min-h-6 tw-leading-6">
            {{ event.location }}
          </div>
        </div>
        <v-btn
          v-if="canEdit"
          key="edit-location-btn"
          class="-tw-my-1"
          icon
          small
          @click="isEditing = true"
        >
          <v-icon small>mdi-pencil</v-icon>
        </v-btn>
      </div>

      <!-- Map view -->
      <div
        v-if="showMap"
        class="tw-mt-2 tw-h-48 tw-w-full tw-overflow-hidden tw-rounded"
        :id="mapId"
      ></div>
    </div>

    <v-btn
      v-else-if="canEdit && !isEditing"
      text
      class="-tw-ml-2 tw-mt-0 tw-w-min tw-px-2 tw-text-dark-gray"
      @click="isEditing = true"
    >
      + Add location
    </v-btn>
    <div
      :class="
        canEdit && !showLocation && isEditing
          ? ''
          : 'tw-absolute tw-opacity-0'
      "
    >
      <div
        class="-tw-mt-[6px] tw-flex tw-w-full tw-flex-grow tw-items-center tw-gap-2"
      >
        <v-text-field
          v-model="newLocation"
          placeholder="Enter a location..."
          class="tw-flex-grow tw-p-2 tw-text-xs sm:tw-text-sm"
          autofocus
          hide-details
        ></v-text-field>
        <v-btn
          icon
          :small="isPhone"
          @click="
            newLocation = event.location
            isEditing = false
          "
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <v-btn icon :small="isPhone" color="primary" @click="saveLocation"
          ><v-icon>mdi-check</v-icon></v-btn
        >
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from "vuex"
import { isPhone, put } from "@/utils"

export default {
  name: "EventLocation",

  props: {
    event: {
      type: Object,
      required: true,
    },
    canEdit: {
      type: Boolean,
      required: true,
    },
  },

  data() {
    return {
      isEditing: false,
      newLocation: this.event.location ?? "",
      showMap: false,
      mapId: `map-${Math.random().toString(36).substring(7)}`,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showLocation() {
      return this.event.location && !this.isEditing
    },
  },

  watch: {
    showLocation: {
      immediate: true,
      handler(val) {
        if (val) {
          this.$nextTick(() => {
            this.initMap()
          })
        }
      },
    },
  },

  methods: {
    ...mapActions(["showError"]),
    saveLocation() {
      const oldEvent = { ...this.event }

      const newEvent = {
        ...this.event,
        location: this.newLocation,
      }

      const eventPayload = {
        name: this.event.name,
        duration: this.event.duration,
        dates: this.event.dates,
        type: this.event.type,
        location: this.newLocation,
      }

      this.$emit("update:event", newEvent)
      this.isEditing = false
      put(`/events/${this.event._id}`, eventPayload).catch((err) => {
        console.error(err)
        this.showError("Failed to save location! Please try again later.")
        this.$emit("update:event", {
          ...oldEvent,
        })
      })
    },
    initMap() {
      if (!this.event.location || this.showMap) return

      // Simple OpenStreetMap embed using Leaflet
      // For now, we'll use a simpler approach with an iframe to OSM
      const location = encodeURIComponent(this.event.location)

      // Create a simple map container with OpenStreetMap
      const mapContainer = document.getElementById(this.mapId)
      if (!mapContainer) return

      // Use Nominatim to geocode the address
      fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${location}`
      )
        .then((response) => response.json())
        .then((data) => {
          if (data && data.length > 0) {
            const { lat, lon } = data[0]

            // Create an OpenStreetMap embed
            const iframe = document.createElement("iframe")
            iframe.width = "100%"
            iframe.height = "100%"
            iframe.frameBorder = "0"
            iframe.scrolling = "no"
            iframe.marginHeight = "0"
            iframe.marginWidth = "0"
            iframe.src = `https://www.openstreetmap.org/export/embed.html?bbox=${
              lon - 0.01
            },${lat - 0.01},${lon + 0.01},${
              lat + 0.01
            }&layer=mapnik&marker=${lat},${lon}`

            mapContainer.innerHTML = ""
            mapContainer.appendChild(iframe)
            this.showMap = true
          }
        })
        .catch((err) => {
          console.error("Error geocoding location:", err)
        })
    },
  },
}
</script>
