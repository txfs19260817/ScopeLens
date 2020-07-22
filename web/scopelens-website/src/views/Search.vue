<template>
    <v-container>
        <v-container v-if="!gotResult" id="searchbar" class="fill-height">
            <v-row align="center" justify="center" no-gutters>
                <v-col>
                    <v-card class="elevation-2 card">
                        <v-card-text>
                            <h1 class="text-start display-1 mb-10 fg-text"> {{$t('search.title')}} </h1>
                            <v-form class="searchbar-form" @submit.prevent="goSearch">
                                <FormatSelector :value.sync="criteria.format"
                                                :hint="$t('upload.hint.format')"></FormatSelector>
                                <PokemonSelector :value.sync="criteria.pokemon"
                                                 :hint="$t('upload.hint.pokemon')"></PokemonSelector>
                                <v-checkbox v-model="criteria.has_showdown" :label="$t('search.hasShowdown')"></v-checkbox>
                                <v-checkbox v-model="criteria.has_rental" :label="$t('search.hasRental')"></v-checkbox>
                                <v-radio-group v-model="criteria.order_by" row>
                                    <template v-slot:label>
                                        <div>{{ $t('search.orderBy.title') }}</div>
                                    </template>
                                    <v-radio :label="$t('search.orderBy.time')" value="time"></v-radio>
                                    <v-radio :label="$t('search.orderBy.likes')" value="likes"></v-radio>
                                </v-radio-group>
                                <div class="text-center mt-6">
                                    <v-btn color="primary" type="submit" large dark :loading="loading">
                                        <v-icon left dark>search</v-icon>
                                        {{ $t('search.btn') }}
                                    </v-btn>
                                </div>
                            </v-form>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>
        </v-container>
        <v-container v-else>
            <v-btn color="primary" fab large dark @click="reset">
                <v-icon>mdi-arrow-left</v-icon>
            </v-btn>
            <ResultsLayout :teams="teams"></ResultsLayout>
            <v-col>
                <v-pagination v-if="total!==0" v-model="curPage" :length="pageLen" total-visible="8"
                              @input="getTeamsSearch"></v-pagination>
                <v-subheader class="justify-center" v-else>{{ $t('results.noResult') }}</v-subheader>
            </v-col>
        </v-container>
    </v-container>
</template>

<script>
    import ResultsLayout from "../components/layouts/ResultsLayout";
    import {getTeamsBySearchCriteria} from "../api/team";
    import {ERROR} from "../api";
    import FormatSelector from "../components/selectors/FormatSelector";
    import PokemonSelector from "../components/selectors/PokemonSelector";

    export default {
        name: "Search",
        components: {
            FormatSelector,
            PokemonSelector,
            ResultsLayout
        },
        data() {
            return {
                teams: [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},],
                gotResult: false,
                // search criteria form to be uploaded
                criteria: {
                    format: '',
                    pokemon: [],
                    has_showdown: false,
                    has_rental: false,
                    order_by: "time",
                },
                // page
                total: 1, // total data size
                pageSize: 12, // data size per page
                curPage: 1,
            }
        },
        methods: {
            async getTeamsSearch(page) {
                // format and pokemon should not be empty at the same time
                if (this.criteria.format.length === 0 && this.criteria.pokemon.length === 0) return

                // loading
                this.$store.commit('LOADING_ON')

                // process Pokemon names ('A/B/C' --> 'A')
                this.criteria.pokemon.forEach((item, idx) => this.criteria.pokemon[idx] = item.toString().split('/', 1)[0])

                const res = await getTeamsBySearchCriteria(page, this.criteria)
                if (res.data.code === ERROR) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('api.thenError') + res.data.msg,
                        "color": "error"
                    });
                } else {
                    this.$emit("results", res.data.data)
                    this.teams = res.data.data.teams
                    this.total = res.data.data.total
                    console.log(res.data.data)
                }
                this.gotResult = true
                this.$store.commit('LOADING_OFF')
            },
            goSearch() {
                this.getTeamsSearch(1)
            },
            reset() {
                this.gotResult = false;
                this.criteria.pokemon = [];
                this.criteria.format = '';
                this.curPage = this.total = 1;
            }
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
            pageLen() {
                return Math.ceil(this.total / this.pageSize)
            }
        }
    }
</script>

<style scoped lang="scss">
    a.no-text-decoration {
        text-decoration: none;
    }

    #searchbar {
        max-width: 70rem;
    }

    .searchbar-form {
        max-width: 40rem;
        margin: 0 auto;
    }

    .card {
        overflow: hidden;
    }

    .fg-text {
        color: #4768A1;
    }
</style>