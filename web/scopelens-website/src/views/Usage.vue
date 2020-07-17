<template>
    <v-container>
        <v-container class="searchbar fill-height">
            <v-row align="center" justify="center" no-gutters>
                <v-col>
                    <v-card class="elevation-2 card">
                        <v-card-text>
                            <h1 class="text-start display-1 mb-10 fg-text"> Pokemon Usage </h1>
                            <v-form class="searchbar-form" @submit.prevent="getUsage">
                                <FormatSelector :value.sync="format"></FormatSelector>
                                <div class="text-center mt-6">
                                    <v-btn color="primary" type="submit" large dark :loading="loading">
                                        <v-icon left dark>search</v-icon>
                                        Search
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
                                "msg": "No usage retrieved because no Pokemon used more than once among all teams under this format. ",
                                "color": "warning"
                            });
                        }
                        console.log(this.usage)
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": "Failed to retrieve usages from server! " + res.data.msg,
                            "color": "error"
                        });
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": "Failed to connect to server! ",
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

    .fg-text {
        color: #4768A1;
    }

    #pie {
        width: 800px;
        height: 800px;
    }

    @media screen and (max-width: 800px) {
        #pie {
            width: 600px;
            height: 600px;
        }
    }
    @media screen and (max-width: 400px) {
        #pie {
            width: 280px;
            height: 280px;
        }
    }

</style>