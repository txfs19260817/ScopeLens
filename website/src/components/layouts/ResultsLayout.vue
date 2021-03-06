<template>
  <v-container class="fill-height">
    <v-row align="start" justify="space-around">
      <v-col v-for="(t, i) in teams" :key="i" xs="12" cols="6" md="4" xl="3">
        <v-hover v-slot:default="{ hover }">
          <v-skeleton-loader
              v-if="loading"
              class="mx-auto"
              height="360"
              max-width="640"
              type="card, list-item-two-line"
              elevation="2"
          >
          </v-skeleton-loader>
          <v-card v-else class="mx-auto" max-width="640" :elevation="hover ? 10 : 2" :class="{ 'on-hover': hover }" link :to="'/team/' + t.id">
            <v-img class="align-end" :src="RedirectToCloudFront(t.image)" :aspect-ratio="16/9" :lazy-src="require('@/assets/teamlazyload.jpg')"></v-img>
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-card-subtitle class="pb-0 font-weight-bold card-text" v-bind="attrs" v-on="on">
                  {{ "[" + t.format + "] " }}{{ t.title }}
                </v-card-subtitle>
              </template>
              <span>{{ "[" + t.format + "] " }}{{ t.title }}</span>
            </v-tooltip>

            <v-card-text class="text--primary">
              <div class="card-text">by {{ t.author }}</div>
              <div>{{ DateConversion(t.created_at) }}</div>
              <v-row align="center" justify="end">
                <v-tooltip bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-icon v-if="t.has_showdown" class="mr-1 light-blue--text accent-2" v-bind="attrs" v-on="on">
                      mdi-alpha-s-circle
                    </v-icon>
                  </template>
                  <span>{{ $t('results.hasShowdown') }}</span>
                </v-tooltip>

                <v-tooltip bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-icon v-if="t.has_rental" class="mr-1 teal--text accent-2" v-bind="attrs" v-on="on">
                      mdi-alpha-r-circle
                    </v-icon>
                  </template>
                  <span>{{ $t('results.hasRental') }}</span>
                </v-tooltip>

                <v-icon class="mr-1 red--text">mdi-heart</v-icon>
                <span class="subheading mr-2">{{ t.likes }}</span>
              </v-row>
            </v-card-text>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  import {DateConversion, RedirectToCloudFront} from "../../assets/utils";

  export default {
    name: "ResultsLayout",
    props: {
      teams: {
        type: Array,
        default: () => [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},],// dummy obj * pageSize
      }
    },
    methods: {
      DateConversion: DateConversion,
      RedirectToCloudFront: RedirectToCloudFront,
    },
    computed: {
      loading() {
        return this.$store.state.loading.loading
      },
    }
  }
</script>

<style scoped>
  .card-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  @media screen and (max-width: 1000px) {
    .card-text {
      white-space: nowrap;
      overflow: auto;
      text-overflow: unset;
    }
  }

  .v-card {
    transition: opacity .3s ease-in-out;
  }

  .v-card:not(.on-hover) {
    opacity: 0.86;
  }
</style>