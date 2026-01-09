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

  mounted() {
    console.log('[EventLocation] Component mounted', {
      hasLocation: !!this.event.location,
      location: this.event.location,
      canEdit: this.canEdit,
      mapId: this.mapId
    })
  },

  data() {
    return {
      isEditing: false,
      newLocation: this.event.location ?? "",
      showMap: false,
      mapExpanded: false, // Map starts collapsed
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
    'event.location': {
      handler(newVal, oldVal) {
        // Re-initialize map if location changes
        if (newVal && newVal !== oldVal) {
          this.showMap = false
          this.$nextTick(() => {
            this.initMap()
          })
        }
      },
    },
  },

  methods: {
    ...mapActions(["showError"]),
    toggleMapExpanded() {
      this.mapExpanded = !this.mapExpanded
      console.log('[EventLocation] Map toggled:', this.mapExpanded ? 'expanded' : 'collapsed')
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
    async initMap() {
      console.log('[EventLocation] initMap called', {
        hasLocation: !!this.event.location,
        location: this.event.location,
        showMap: this.showMap,
        mapId: this.mapId
      })

      if (!this.event.location) {
        console.log('[EventLocation] No location provided, skipping map initialization')
        return
      }

      if (this.showMap) {
        console.log('[EventLocation] Map already showing, skipping re-initialization')
        return
      }

      // Simple OpenStreetMap embed using Leaflet
      // For now, we'll use a simpler approach with an iframe to OSM
      const location = this.event.location

      // Create a simple map container with OpenStreetMap
      const mapContainer = document.getElementById(this.mapId)
      if (!mapContainer) {
        console.error('[EventLocation] Map container not found:', this.mapId)
        return
      }

      console.log('[EventLocation] Geocoding address:', this.event.location)

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
        console.log('[EventLocation] Geocoding response:', data)

        // If no results, try with simplified address (remove suite/apartment numbers)
        if (!data || data.length === 0) {
          const simplifiedAddress = location
            .replace(/Suite \d+,?\s*/i, '')
            .replace(/Apt\.? \d+,?\s*/i, '')
            .replace(/#\d+,?\s*/i, '')
          
          if (simplifiedAddress !== location) {
            console.log('[EventLocation] Retrying with simplified address:', simplifiedAddress)
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
              console.log('[EventLocation] Simplified geocoding response:', data)
            }
          }
        }
        
        if (data && data.length > 0) {
          const lat = parseFloat(data[0].lat)
          const lon = parseFloat(data[0].lon)
          console.log('[EventLocation] Creating map at coordinates:', { lat, lon })

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
          iframe.src = `https://www.openstreetmap.org/export/embed.html?bbox=${
            lon - 0.01
          },${lat - 0.01},${lon + 0.01},${
            lat + 0.01
          }&layer=mapnik&marker=${lat},${lon}`

          mapContainer.innerHTML = ""
          mapContainer.appendChild(iframe)
          this.showMap = true
          console.log('[EventLocation] Map initialized successfully')
        } else {
          console.warn('[EventLocation] No geocoding results found for:', this.event.location)
        }
      } catch (err) {
        console.error('[EventLocation] Error geocoding location:', err)
        // Silently fail - user can still see the location text
      }
    },
  },
}
</script>
