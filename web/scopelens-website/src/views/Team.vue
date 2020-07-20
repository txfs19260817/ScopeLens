<template>
    <v-container>
        <v-container v-if="team">
            <div class="row">
                <div class="col-md-5 col-sm-5 col-xs-12">
                    <v-row align="center">
                        <img style="width: 100%" :src="team.image" alt="preview"/>
                    </v-row>
                    <v-row justify="space-around">
                        <v-subheader>{{ $t("team.likes") }} {{ team.likes }}</v-subheader>
                    </v-row>
                    <v-row justify="space-around">

                        <v-tooltip bottom>
                            <template v-slot:activator="{ on, attrs }">
                                <v-btn class="mx-2" fab dark small color="pink" @click="like" :loading="loading"
                                       v-bind="attrs" v-on="on">
                                    <v-icon dark>mdi-heart</v-icon>
                                </v-btn>
                            </template>
                            <span>{{ $t("team.addLike.tooltip") }}</span>
                        </v-tooltip>

                        <v-tooltip bottom>
                            <template v-slot:activator="{ on, attrs }">
                                <v-btn class="mx-2" fab dark small color="blue" v-bind="attrs" v-on="on"
                                       v-if="team.showdown.length>0"
                                       v-clipboard:copy="team.showdown"
                                       v-clipboard:success="onCopy"
                                       v-clipboard:error="onError">
                                    <v-icon dark>mdi-content-copy</v-icon>
                                </v-btn>
                            </template>
                            <span>{{ $t("team.copy.tooltip") }}</span>
                        </v-tooltip>

                        <v-btn class="mx-2" fab dark small color="blue"
                               :href="`http://twitter.com/share?text=Team%20Sharing:%20` + team.title +  `&url=` + url +`&hashtags=ScopeLens,PokemonSwordShield,NintendoSwitch`">
                            <v-icon dark>mdi-twitter</v-icon>
                        </v-btn>

                        <v-btn class="mx-2" fab dark small color="red"
                               :href="`http://service.weibo.com/share/share.php?title=队伍分享：` + team.title + `&url=`+ url + `&pic=` + team.image">
                            <v-icon dark>mdi-sina-weibo</v-icon>
                        </v-btn>
                    </v-row>
                </div>
                <div class="col-md-7 col-sm-7 col-xs-12 pl-6">
                    <p class="display-1 mb-0">{{ "[" + team.format + "] " }}{{ team.title }}</p>
                    <p class="subtitle-2 font-weight-light">{{ DateConversion(team.created_at) }}</p>
                    <v-card-actions class="pa-0">
                        <p class="headline font-weight-light pt-3">{{ $t("team.author") }} {{team.author}}
                            <span class="subtitle-2 font-weight-light">{{ $t("team.uploader") }} {{team.uploader}}</span>
                        </p>
                    </v-card-actions>
                    <p class="title">{{ $t("team.format") }}</p>
                    <v-chip class="ma-2" color="primary"> {{ team.format }}</v-chip>
                    <p class="title">{{ $t("team.pokemon") }}</p>
                    <v-breadcrumbs :items="team.pokemon">
                        <template v-slot:item="{ item }">
                            <v-breadcrumbs-item>
                                <v-img max-height="50" max-width="50" :alt="item"
                                       :src="iconUrl + ProcessStr(item) + `.png`"></v-img>
                            </v-breadcrumbs-item>
                        </template>
                    </v-breadcrumbs>
                </div>
            </div>
            <div class="row">
                <div class="col-sm-12 col-xs-12 col-md-12">
                    <v-tabs>
                        <v-tab>Showdown</v-tab>
                        <v-tab-item>
                            <pre v-if="team.showdown.length>0">
                                <p class="body-2"> {{ team.showdown }} </p>
                            </pre>
                            <p v-else class="pt-10 body-1"> {{ $t("team.noShowdown") }} </p>
                        </v-tab-item>
                        <v-tab>{{ $t("team.description") }}</v-tab>
                        <v-tab-item>
                            <pre class="pt-10 body-1"> {{ team.description }} </pre>
                        </v-tab-item>
                    </v-tabs>
                </div>
            </div>
        </v-container>
        <v-container>
            <v-row align="start" justify="center">
                <v-progress-circular v-if="loading" :size="80" color="primary" indeterminate></v-progress-circular>
            </v-row>
        </v-container>
    </v-container>
</template>
<script>
    import {ProcessStr, DateConversion} from "../assets/utils"
    import {ERROR, logErrors, SUCCESS} from "../api";
    import {getTeamByID, insertLikeByUsername} from "../api/team";

    export default {
        name: "Team",
        data: () => ({
            team: null,
            // Sprites paths
            iconUrl: process.env.VUE_APP_STATIC_ASSET_URL + '2d/',
        }),
        methods: {
            getTeamDetail(id) {
                // loading
                this.$store.commit('LOADING_ON')
                getTeamByID(id).then(res => {
                    if (res.data.code === SUCCESS) {
                        this.team = res.data.data
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
            async like() {
                // loading
                this.$store.commit('LOADING_ON')
                const res = await insertLikeByUsername(this.$store.state.user.username, this.$route.params.id, this.$store.state.user.token)
                if (res.data.code === ERROR || res.status === 401) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t("team.addLike.failed") + res.data.msg,
                        "color": "error"
                    });
                } else {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t("team.addLike.success"),
                        "color": "success"
                    });
                }
                this.$store.commit('LOADING_OFF')
            },
            onCopy(e) {
                this.$store.dispatch('snackbar/openSnackbar', {
                    "msg": this.$t("team.copy.success"),
                    "color": "success"
                });
            },
            onError(e) {
                this.$store.dispatch('snackbar/openSnackbar', {
                    "msg": this.$t("team.copy.failed"),
                    "color": "error"
                });
            },
            ProcessStr: ProcessStr,
            DateConversion: DateConversion
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
            url() {
                return encodeURIComponent(`https://scopelens.team/#` + this.$route.path)
            }
        },
        created() {
            this.getTeamDetail(this.$route.params.id)
        }
    }
</script>
<style scoped>
    .textarea {
        text-align: justify;
        text-justify: newspaper;
        word-break: break-word;
    }

    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
    }
</style>