<template>
  <div class="tw-mt-1 tw-max-w-full sm:tw-mt-2 sm:tw-max-w-[calc(100%-236px)]">
    <!-- Collapsed view - Show summary button when there's content -->
    <v-btn
      v-if="hasContent && !expanded && !isEditingDescription && !isEditingLocation"
      text
      class="-tw-ml-2 tw-mt-0 tw-px-2 tw-text-dark-gray"
      @click="expanded = true"
    >
      <v-icon small class="tw-mr-1">mdi-information-outline</v-icon>
      <span class="tw-flex tw-items-center tw-gap-1">
        View event details
        <span class="tw-ml-1 tw-flex tw-gap-1">
          <v-icon v-if="showDescription" x-small class="tw-text-gray-500">mdi-text</v-icon>
          <v-icon v-if="showLocation" x-small class="tw-text-gray-500">mdi-map-marker</v-icon>
        </span>
      </span>
    </v-btn>

    <!-- Add details button when no content exists -->
    <v-btn
      v-else-if="!hasContent && canEdit && !expanded && !isEditingDescription && !isEditingLocation"
      text
      class="-tw-ml-2 tw-mt-0 tw-w-min tw-px-2 tw-text-dark-gray"
      @click="expanded = true"
    >
      + Add event details
    </v-btn>

    <!-- Expanded view -->
    <div
      v-if="expanded || isEditingDescription || isEditingLocation"
      class="tw-flex tw-w-full tw-flex-col tw-gap-2 tw-rounded-md tw-border tw-border-light-gray-stroke tw-bg-light-gray tw-p-2"
    >
      <!-- Header with collapse button -->
      <div class="tw-flex tw-items-center tw-justify-between">
        <span class="tw-text-xs tw-font-semibold tw-text-very-dark-gray sm:tw-text-sm">
          Event Details
        </span>
        <v-btn
          v-if="!isEditingDescription && !isEditingLocation"
          icon
          small
          @click="expanded = false"
          title="Hide details"
        >
          <v-icon small>mdi-chevron-up</v-icon>
        </v-btn>
      </div>

      <!-- Description Section -->
      <div class="tw-flex tw-flex-col tw-gap-1">
        <div
          v-if="showDescription && !isEditingDescription"
          class="tw-flex tw-items-start tw-gap-2"
        >
          <v-icon small class="tw-mt-0.5">mdi-text</v-icon>
          <div class="tw-grow tw-space-y-1 tw-text-xs tw-font-normal tw-text-very-dark-gray sm:tw-text-sm">
            <div
              class="tw-min-h-6 tw-leading-6"
              v-for="(line, i) in event.description.split('\n')"
              :key="i"
            >
              {{ line }}
            </div>
          </div>
          <v-btn
            v-if="canEdit"
            icon
            small
            @click="startEditingDescription"
            title="Edit description"
          >
            <v-icon small>mdi-pencil</v-icon>
          </v-btn>
        </div>

        <!-- Add description button -->
        <v-btn
          v-else-if="canEdit && !showDescription && !isEditingDescription"
          text
          x-small
          class="-tw-ml-2 tw-w-min tw-px-2 tw-text-dark-gray"
          @click="startEditingDescription"
        >
          + Add description
        </v-btn>

        <!-- Edit description form -->
        <div v-if="isEditingDescription" class="tw-flex tw-items-start tw-gap-2">
          <v-icon small class="tw-mt-2">mdi-text</v-icon>
          <div class="tw-flex tw-flex-grow tw-flex-col tw-gap-1">
            <v-textarea
              v-model="newDescription"
              placeholder="Enter a description..."
              class="tw-text-xs sm:tw-text-sm"
              autofocus
              :rows="1"
              auto-grow
              hide-details
            ></v-textarea>
            <div class="tw-flex tw-gap-1">
              <v-btn
                small
                @click="cancelEditingDescription"
              >
                Cancel
              </v-btn>
              <v-btn
                small
                color="primary"
                @click="saveDescription"
              >
                Save
              </v-btn>
            </div>
          </div>
        </div>
      </div>

      <!-- Location Section -->
      <div class="tw-flex tw-flex-col tw-gap-1">
        <div
          v-if="showLocation && !isEditingLocation"
          class="tw-flex tw-flex-col tw-gap-2"
        >
          <div class="tw-flex tw-items-start tw-gap-2">
            <v-icon small class="tw-mt-0.5">mdi-map-marker</v-icon>
            <div class="tw-grow tw-cursor-pointer tw-text-xs tw-font-normal tw-text-very-dark-gray sm:tw-text-sm" @click="toggleMapExpanded">
              <div class="tw-min-h-6 tw-leading-6">
                {{ event.location }}
              </div>
            </div>
            <v-btn
              v-if="showMap && !mapExpanded"
              icon
              small
              @click="toggleMapExpanded"
              title="Show map"
            >
              <v-icon small>mdi-chevron-down</v-icon>
            </v-btn>
            <v-btn
              v-else-if="showMap && mapExpanded"
              icon
              small
              @click="toggleMapExpanded"
              title="Hide map"
            >
              <v-icon small>mdi-chevron-up</v-icon>
            </v-btn>
            <v-btn
              v-if="canEdit"
              icon
              small
              @click="startEditingLocation"
              title="Edit location"
            >
              <v-icon small>mdi-pencil</v-icon>
            </v-btn>
          </div>

          <!-- Map view - Collapsible -->
          <v-expand-transition>
            <div
              v-show="showMap && mapExpanded"
              class="tw-h-48 tw-w-full tw-overflow-hidden tw-rounded"
              :id="mapId"
            ></div>
          </v-expand-transition>
        </div>

        <!-- Add location button -->
        <v-btn
          v-else-if="canEdit && !showLocation && !isEditingLocation"
          text
          x-small
          class="-tw-ml-2 tw-w-min tw-px-2 tw-text-dark-gray"
          @click="startEditingLocation"
        >
          + Add location
        </v-btn>

        <!-- Edit location form -->
        <div v-if="isEditingLocation" class="tw-flex tw-items-start tw-gap-2">
          <v-icon small class="tw-mt-2">mdi-map-marker</v-icon>
          <div class="tw-flex tw-flex-grow tw-flex-col tw-gap-1">
            <v-text-field
              v-model="newLocation"
              placeholder="Enter a location..."
              class="tw-text-xs sm:tw-text-sm"
              autofocus
              hide-details="auto"
              counter="500"
              maxlength="500"
            ></v-text-field>
            <div class="tw-flex tw-gap-1">
              <v-btn
                small
                @click="cancelEditingLocation"
              >
                Cancel
              </v-btn>
              <v-btn
                small
                color="primary"
                @click="saveLocation"
              >
                Save
              </v-btn>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from "vuex"
import { isPhone, put } from "@/utils"

export default {
  name: "EventDetails",

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
      expanded: false, // Collapsed by default
      isEditingDescription: false,
      isEditingLocation: false,
      newDescription: this.event.description ?? "",
      newLocation: this.event.location ?? "",
      showMap: false,
      mapExpanded: false,
      mapId: `map-${Math.random().toString(36).substring(7)}`,
      mapInitialized: false,
      mapLat: null,
      mapLon: null,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showDescription() {
      return this.event.description
    },
    showLocation() {
      return this.event.location
    },
    hasContent() {
      return this.showDescription || this.showLocation
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
    
    startEditingDescription() {
      this.isEditingDescription = true
      this.newDescription = this.event.description ?? ""
    },
    
    cancelEditingDescription() {
      this.isEditingDescription = false
      this.newDescription = this.event.description ?? ""
    },
    
    startEditingLocation() {
      this.isEditingLocation = true
      this.newLocation = this.event.location ?? ""
    },
    
    cancelEditingLocation() {
      this.isEditingLocation = false
      this.newLocation = this.event.location ?? ""
    },

    toggleMapExpanded() {
      this.mapExpanded = !this.mapExpanded
    },

    saveDescription() {
      const oldEvent = { ...this.event }

      const newEvent = {
        ...this.event,
        description: this.newDescription,
      }

      const eventPayload = {
        name: this.event.name,
        duration: this.event.duration,
        dates: this.event.dates,
        type: this.event.type,
        description: this.newDescription,
        location: this.event.location,
      }

      this.$emit("update:event", newEvent)
      this.isEditingDescription = false
      put(`/events/${this.event._id}`, eventPayload).catch((err) => {
        console.error(err)
        this.showError("Failed to save description! Please try again later.")
        this.$emit("update:event", {
          ...oldEvent,
        })
      })
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
      this.isEditingLocation = false
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

          this.mapLat = lat
          this.mapLon = lon
          this.showMap = true

          if (this.mapExpanded) {
            this.$nextTick(() => {
              this.renderMap()
            })
          }
        }
      } catch (err) {
        console.error('[EventDetails] Error geocoding location:', err)
      }
    },

    calculateBboxFromZoom(lat, lon, zoom) {
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
        console.error('[EventDetails] Missing coordinates for map rendering')
        return
      }

      const mapContainer = document.getElementById(this.mapId)
      if (!mapContainer) {
        console.error('[EventDetails] Map container not found:', this.mapId)
        return
      }

      const lat = this.mapLat
      const lon = this.mapLon

      const isSmallScreen = window.innerWidth < 640
      const zoom = isSmallScreen ? 18 : 17

      const bbox = this.calculateBboxFromZoom(lat, lon, zoom)

      const bboxString = `${bbox.west},${bbox.south},${bbox.east},${bbox.north}`
      const mapUrl = `https://www.openstreetmap.org/export/embed.html?bbox=${bboxString}&layer=mapnik&marker=${lat},${lon}`

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
