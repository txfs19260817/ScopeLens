<template>
    <v-container>
        <v-tabs class="fill-height" grow slider-color="primary" background-color="bg_primary">
            <v-tab class="font-weight-bold title">{{ $t('usage.tabs.official') }}</v-tab>
            <v-tab-item class="text-center">
                <div class="hidden-xs-only tab-item-wrapper">
                    <iframe :src="officialSrc"></iframe>
                </div>
                <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark, 'hidden-sm-and-up': true}"
                       style="margin: 20px" :href="officialSrc" target="_blank"
                       color="primary" type="submit" large dark :loading="loading">
                    <v-icon left dark>mdi-open-in-new</v-icon>
                    {{ $t('usage.newWindow') }}
                </v-btn>
            </v-tab-item>
            <v-tab class="font-weight-bold title">{{ $t('usage.tabs.website') }}</v-tab>
            <v-tab-item>
                <v-container class="searchbar fill-height">
                    <v-row align="center" justify="center" no-gutters>
                        <v-col>
                            <v-card class="elevation-2 card">
                                <v-card-text>
                                    <h1 class="text-start display-1 mb-10 primary--text"> {{$t('usage.title')}} </h1>
                                    <v-form class="searchbar-form" @submit.prevent="getUsage">
                                        <FormatSelector :value.sync="format"
                                                        :hint="$t('upload.hint.format')"></FormatSelector>
                                        <div class="text-center mt-6">
                                            <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}"
                                                   color="primary" type="submit" large dark :loading="loading">
                                                <v-icon left dark>search</v-icon>
                                                {{$t('usage.btn')}}
                                            </v-btn>
                                        </div>
                                    </v-form>
                                </v-card-text>
                            </v-card>
                        </v-col>
                    </v-row>
                </v-container>
                <v-container v-show="usage.length > 0" class="searchbar fill-height">
                    <v-row align="center" justify="center" no-gutters>
                        <v-card class="elevation-2 card">
                            <v-card-text>
                                <div id="pie"></div>
                            </v-card-text>
                        </v-card>
                    </v-row>
                </v-container>
            </v-tab-item>
        </v-tabs>

    </v-container>
</template>

<script>
    import FormatSelector from "../components/selectors/FormatSelector";
    import {getPokemonUsageByFormat} from "../api/team";
    import {logErrors, SUCCESS} from "../api";
    import {ProcessStr} from "../assets/utils";

    let echarts = require('echarts');

    export default {
        name: "Usage",
        components: {
            FormatSelector,
        },
        data() {
            return {
                format: "",
                usage: [],
                // Sprites paths
                iconUrl: process.env.VUE_APP_STATIC_ASSET_URL + '2d/',
                officialSrc: "https://resource.pokemon-home.com/battledata/rankmatch_detail.html?l=9&c=302&v=443199&lcc=1595518253"
            }
        },
        methods: {
            getUsage() {
                // loading
                if (this.format === "") return
                this.$store.commit('LOADING_ON')
                getPokemonUsageByFormat(this.format).then(res => {
                    if (res.data.code === SUCCESS) {
                        if (res.data.data !== null) {
                            this.usage = res.data.data
                            this.drawPieChart(this.usage);
                        } else {
                            // no usage retrieved
                            this.usage = [];
                            this.$store.dispatch('snackbar/openSnackbar', {
                                "msg": this.$t('usage.noUsage'),
                                "color": "error"
                            });
                        }
                        console.log(this.usage)
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('api.thenError') + res.data.msg,
                            "color": "error"
                        });
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('api.catchError'),
                        "color": "error"
                    });
                }).finally(() => {
                    this.$store.commit('LOADING_OFF')
                })
            },
            drawPieChart(d) {
                const iconUrl = process.env.VUE_APP_STATIC_ASSET_URL + '2d/'
                this.echarts.init(document.getElementById('pie'), 'light').setOption({
                    title: {
                        text: this.format,
                        left: 'left'
                    },
                    tooltip: {
                        trigger: 'item',
                        formatter: function (params, ticket, callback) {
                            let img = `<img src="` + iconUrl + ProcessStr(params.data.name) + `.png` + `" />`;
                            let res = params.seriesName + img + "<br />" +
                                params.data.name + ': ' + params.data.value +
                                ' (' + params.percent + '%)';
                            setTimeout(function () {
                                callback(ticket, res);
                            }, 100);
                            return 'loading...';
                        }
                    },
                    series: [
                        {
                            name: this.format + ' Usage',
                            type: 'pie',
                            radius: '50%',
                            roseType: 'angle',
                            data: d
                        }
                    ],
                });
            },
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
            tabs() {
                return [
                    this.$i18n.t('usage.tabs.official'),
                    this.$i18n.t('usage.tabs.website'),
                ]
            }
        },
        created() {
            this.echarts = echarts;
        }
    }
</script>

<style scoped>
    .searchbar {
        max-width: 70rem;
    }

    .searchbar-form {
        max-width: 40rem;
        margin: 0 auto;
    }

    .card {
        overflow: hidden;
    }

    #pie {
        width: 800px;
        height: 800px;
    }

    @media screen and (max-width: 800px) {
        #pie {
            width: 400px;
            height: 400px;
        }
    }

    @media screen and (max-width: 420px) {
        #pie {
            width: 280px;
            height: 280px;
        }
    }

    .tab-item-wrapper {
        /* vuetify sets the v-tabs__container height to 48px */
        height: calc(100vh - 160px);
        overflow: auto;
        -webkit-overflow-scrolling:touch;
        width:100%;
    }

    iframe {
        border:none;
        width: 1px;
        min-width: 100%;
        *width: 100%;
        height: 100%;
    }

</style>