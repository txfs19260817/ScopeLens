<template>
    <v-container class="fill-height">
        <v-row align="start" justify="space-around">
            <v-col v-for="(t, i) in teams" :key="i" xs="12" cols="6" md="4" xl="3">
                <v-hover v-slot:default="{ hover }">
                    <v-skeleton-loader
                            v-if="loading"
                            class="mx-auto"
                            height="340"
                            max-width="640"
                            type="card, list-item-two-line"
                            elevation="2"
                    >
                    </v-skeleton-loader>
                    <v-card v-else link class="mx-auto" max-width="640" :elevation="hover ? 10 : 2"
                            :class="{ 'on-hover': hover }">
                        <v-img class="white--text align-end" :src="t.image"></v-img>
                        <v-card-subtitle class="pb-0 font-weight-bold">
                            {{ "[" + t.format + "] " }}{{ t.title }}
                        </v-card-subtitle>
                        <v-card-text class="text--primary">
                            <div>by {{ t.author }}</div>
                            <div>{{ DateConversion(t.created_at) }}</div>
                        </v-card-text>
                    </v-card>
                </v-hover>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import {DateConversion} from "../../assets/utils";

    export default {
        name: "ResultsLayout",
        props: {
            teams: {
                type: Array,
                default: () => [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},],// dummy obj * pageSize
            }
        },
        methods: {
            DateConversion: DateConversion
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
        }
    }
</script>

<style scoped>

</style>