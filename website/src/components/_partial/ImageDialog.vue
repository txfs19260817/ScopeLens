<template>
  <div class="text-center">
    <v-dialog v-model="valueModel" max-width="1280">
      <v-card>
        <v-img :src="url" alt="preview" :aspect-ratio="16/9" :lazy-src="require('@/assets/teamlazyload.jpg')">
          <template v-slot:placeholder>
            <v-row class="fill-height ma-0" align="center" justify="center">
              <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </v-row>
          </template>
        </v-img>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" dark fab @click="valueModel = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  name: "ImageDialog",
  props: {
    show: {
      type: Boolean,
      default: false
    },
    url: {
      type: String,
      required: true
    }
  },
  computed: {
    valueModel: {
      get() {
        return this.show
      },
      set(v) {
        this.$emit('update:show', v)
      },
    },
  }
}
</script>

<style scoped>

</style>