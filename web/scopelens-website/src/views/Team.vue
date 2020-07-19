<template>
    <v-container>
        <v-container v-if="team">
            <div class="row">
                <div class="col-md-5 col-sm-5 col-xs-12">
                    <v-img class="align-end" :src="team.image"></v-img>
                    <v-row justify="space-around">
                        <v-subheader>Likes: {{team.likes}}</v-subheader>
                    </v-row>
                    <v-row justify="space-around">
                        <v-btn class="mx-2" fab dark small color="blue" disabled>
                            <v-icon dark>mdi-twitter</v-icon>
                        </v-btn>
                        <v-btn class="mx-2" fab dark small color="pink" @click="like" :loading="loading">
                            <v-icon dark>mdi-heart</v-icon>
                        </v-btn>
                        <v-btn class="mx-2" fab dark small color="red" disabled>
                            <v-icon dark>mdi-sina-weibo</v-icon>
                        </v-btn>
                    </v-row>
                </div>
                <div class="col-md-7 col-sm-7 col-xs-12 pl-6">
                    <p class="display-1 mb-0">{{ "[" + team.format + "] " }}{{team.title}}</p>
                    <v-card-actions class="pa-0">
                        <p class="headline font-weight-light pt-3">Author by: {{team.author}}
                            <span class="subtitle-2 font-weight-light">Uploaded by: {{team.uploader}}</span>
                        </p>
                    </v-card-actions>
                    <p class="title">Format</p>
                    <v-chip class="ma-2" color="primary"> {{ team.format }}</v-chip>
                    <p class="title">Pokemon</p>
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
                            <pre v-if="team.showdown.length>0" class="pt-10 body-2 textarea"> {{ team.showdown }} </pre>
                            <p v-else class="pt-10 body-1 textarea"> No Showdown team provided. </p>
                        </v-tab-item>
                        <v-tab>Description</v-tab>
                        <v-tab-item>
                            <pre class="pt-10 body-1 textarea"> {{ team.description }} </pre>
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
    import {ProcessStr} from "../assets/utils"
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
                        console.log(this.team.showdown)
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": "Failed to retrieve the team from server! " + res.data.msg,
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
            async like() {
                // loading
                this.$store.commit('LOADING_ON')
                const res = await insertLikeByUsername(this.$store.state.user.username, this.$route.params.id, this.$store.state.user.token)
                if (res.data.code === ERROR || res.status === 401) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": "Like this team error: " + res.data.msg,
                        "color": "error"
                    });
                } else {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": "Added to your liked teams! ",
                        "color": "success"
                    });
                }
                this.$store.commit('LOADING_OFF')
            },
            ProcessStr: ProcessStr,
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
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
        word-break: break-all;
    }
</style>