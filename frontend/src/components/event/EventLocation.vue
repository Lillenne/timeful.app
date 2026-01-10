<template>
  <div class="tw-mt-1 tw-max-w-full sm:tw-mt-2 sm:tw-max-w-[calc(100%-236px)]">
    <div
      v-if="showLocation"
      class="tw-flex tw-w-full tw-flex-col tw-gap-2 tw-rounded-md tw-border tw-border-light-gray-stroke tw-bg-light-gray tw-p-2 tw-text-xs tw-font-normal tw-text-very-dark-gray sm:tw-text-sm"
    >
      <div class="tw-flex tw-items-start tw-gap-2">
        <v-icon small class="tw-mt-0.5">mdi-map-marker</v-icon>
        <div class="tw-grow tw-cursor-pointer" @click="toggleMapExpanded">
          <div class="tw-min-h-6 tw-leading-6">
            {{ event.location }}
          </div>
        </div>
        <v-btn
          v-if="showMap && !mapExpanded"
          key="expand-map-btn"
          class="-tw-my-1"
          icon
          small
          @click="toggleMapExpanded"
          title="Show map"
        >
          <v-icon small>mdi-chevron-down</v-icon>
        </v-btn>
        <v-btn
          v-else-if="showMap && mapExpanded"
          key="collapse-map-btn"
          class="-tw-my-1"
          icon
          small
          @click="toggleMapExpanded"
          title="Hide map"
        >
          <v-icon small>mdi-chevron-up</v-icon>
        </v-btn>
        <v-btn
          v-if="canEdit"
          key="edit-location-btn"
          class="-tw-my-1"
          icon
          small
          @click="isEditing = true"
          title="Edit location"
        >
          <v-icon small>mdi-pencil</v-icon>
        </v-btn>
      </div>

      <!-- Map view - Collapsible, starts minimized -->
      <v-expand-transition>
        <div
          v-show="showMap && mapExpanded"
          class="tw-h-48 tw-w-full tw-overflow-hidden tw-rounded"
          :id="mapId"
        ></div>
      </v-expand-transition>
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
          hide-details="auto"
          counter="500"
          maxlength="500"
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
      mapExpanded: false, // Map starts collapsed
      mapId: `map-${Math.random().toString(36).substring(7)}`,
      mapInitialized: false, // Track if map has been rendered
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
            this.geocodeLocation()
          })
        }
      },
    },
    'event.location': {
      handler(newVal, oldVal) {
        // Re-initialize map if location changes
        if (newVal && newVal !== oldVal) {
          this.showMap = false
          this.mapInitialized = false
          this.$nextTick(() => {
            this.geocodeLocation()
          })
        }
      },
    },
    mapExpanded: {
      handler(val) {
        // When map is expanded, render it if not already done
        if (val && this.showMap && !this.mapInitialized) {
          this.$nextTick(() => {
            this.renderMap()
          })
        }
      },
    },
  },

  methods: {
    ...mapActions(["showError"]),
    toggleMapExpanded() {
      this.mapExpanded = !this.mapExpanded
    },
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
        description: this.event.description,
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
    async geocodeLocation() {
      if (!this.event.location) {
        return
      }

      const location = this.event.location

      try {
        // Use Nominatim to geocode the address
        // Following Nominatim usage policy: https://operations.osmfoundation.org/policies/nominatim/
        // User-Agent header is required by Nominatim to identify the application
        const response = await fetch(
          `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(location)}`,
          {
            headers: {
              'User-Agent': 'Timeful/1.0 (https://timeful.app)'
            }
          }
        )

        if (!response.ok) {
          throw new Error(`Geocoding failed: ${response.status}`)
        }

        let data = await response.json()

        // If no results, try with simplified address (remove suite/apartment numbers)
        if (!data || data.length === 0) {
          const simplifiedAddress = location
            .replace(/Suite \d+,?\s*/i, '')
            .replace(/Apt\.? \d+,?\s*/i, '')
            .replace(/#\d+,?\s*/i, '')
          
          if (simplifiedAddress !== location) {
            const retryResponse = await fetch(
              `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(simplifiedAddress)}`,
              {
                headers: {
                  'User-Agent': 'Timeful/1.0 (https://timeful.app)'
                }
              }
            )
            
            if (retryResponse.ok) {
              data = await retryResponse.json()
            }
          }
        }
        
        if (data && data.length > 0) {
          const lat = parseFloat(data[0].lat)
          const lon = parseFloat(data[0].lon)
          
          // Store coordinates for later rendering
          this.mapLat = lat
          this.mapLon = lon
          this.showMap = true
          
          // If map is already expanded, render it now
          if (this.mapExpanded) {
            this.$nextTick(() => {
              this.renderMap()
            })
          }
        }
      } catch (err) {
        console.error('[EventLocation] Error geocoding location:', err)
        // Silently fail - user can still see the location text
      }
    },
    calculateBboxFromZoom(lat, lon, zoom) {
      // Calculate the bounding box based on zoom level
      // At zoom 18, we want to show roughly 1-2 city blocks
      // At zoom 17, we want to show roughly 2-3 city blocks
      // The formula: latDelta = 180 / (2^(zoom+1)), lonDelta = latDelta / cos(lat)
      
      const latDelta = 180 / Math.pow(2, zoom + 1)
      const lonDelta = latDelta / Math.cos(lat * Math.PI / 180)
      
      const south = lat - latDelta
      const north = lat + latDelta
      const west = lon - lonDelta
      const east = lon + lonDelta
      
      return { south, north, west, east }
    },
    renderMap() {
      if (this.mapInitialized) {
        return
      }

      if (!this.mapLat || !this.mapLon) {
        console.error('[EventLocation] Missing coordinates for map rendering')
        return
      }

      // Create a simple map container with OpenStreetMap
      const mapContainer = document.getElementById(this.mapId)
      if (!mapContainer) {
        console.error('[EventLocation] Map container not found:', this.mapId)
        return
      }

      const lat = this.mapLat
      const lon = this.mapLon

      // Determine zoom level based on screen size
      // Smaller screens get more zoom (higher zoom number = more zoomed in)
      const isSmallScreen = window.innerWidth < 640
      const zoom = isSmallScreen ? 18 : 17

      // Calculate bounding box from zoom level
      const bbox = this.calculateBboxFromZoom(lat, lon, zoom)
      
      const bboxString = `${bbox.west},${bbox.south},${bbox.east},${bbox.north}`
      const mapUrl = `https://www.openstreetmap.org/export/embed.html?bbox=${bboxString}&layer=mapnik&marker=${lat},${lon}`

      // Create an OpenStreetMap embed
      // sandbox with allow-same-origin and allow-scripts is needed for the map to function
      // This is safe because we're loading from the trusted openstreetmap.org domain
      const iframe = document.createElement("iframe")
      iframe.width = "100%"
      iframe.height = "100%"
      iframe.frameBorder = "0"
      iframe.scrolling = "no"
      iframe.marginHeight = "0"
      iframe.marginWidth = "0"
      iframe.sandbox = "allow-scripts allow-same-origin"
      iframe.referrerPolicy = "no-referrer"
      iframe.src = mapUrl

      mapContainer.innerHTML = ""
      mapContainer.appendChild(iframe)
      this.mapInitialized = true
    },
  },
}
</script>
