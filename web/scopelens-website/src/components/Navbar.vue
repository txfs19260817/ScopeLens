<template>
    <nav>
        <v-app-bar app fixed flat color="secondary" :class="{'lighten-4':!$vuetify.theme.dark, 'darken-4':$vuetify.theme.dark}">
            <v-row justify="space-between">
                <v-app-bar-nav-icon @click="display = !display" class="d-lg-none"></v-app-bar-nav-icon>
                <v-btn v-if="$route.path.includes(`team`)" icon @click="$router.back()">
                    <v-icon>mdi-arrow-left</v-icon>
                </v-btn>
                <v-spacer></v-spacer>
                <DarkModeSwitch></DarkModeSwitch>
                <LanguageSelector></LanguageSelector>
            </v-row>
        </v-app-bar>

        <v-navigation-drawer v-model="display" color="secondary" :class="{'lighten-5':!$vuetify.theme.dark, 'darken-3':$vuetify.theme.dark}" app fixed left flat>
            <v-sheet color="secondary" :class="{'lighten-5':!$vuetify.theme.dark, 'darken-3':$vuetify.theme.dark}">
                <v-list>
                    <v-list-item two-line>
                        <v-list-item-avatar>
                            <v-img :src="require('@/assets/logo.png')"></v-img>
                        </v-list-item-avatar>
                        <v-list-item-title>
                            <span class="font-weight-medium title">Scope</span>
                            <span class="font-weight-light title">Lens</span>
                        </v-list-item-title>
                    </v-list-item>

                    <v-list-group v-if="isLogin" prepend-icon="account_circle">
                        <template v-slot:activator>
                            <v-list-item-title class="font-weight-bold">{{ username }}</v-list-item-title>
                        </template>
                        <template v-for="(item, i) in user.list">
                            <v-list-item :key="i" link :to="item.link">
                                <v-list-item-icon>
                                    <v-icon v-text="item.icon"></v-icon>
                                </v-list-item-icon>
                                <v-list-item-content class="text-body-2">{{ item.text }}</v-list-item-content>
                            </v-list-item>
                        </template>
                    </v-list-group>

                    <v-list-item v-else link to="/login">
                        <v-list-item-content>
                            <v-list-item-title class="title">{{ $t("nav.login") }}</v-list-item-title>
                        </v-list-item-content>
                        <v-list-item-action>
                            <v-icon>mdi-arrow-right</v-icon>
                        </v-list-item-action>
                    </v-list-item>
                </v-list>

                <v-divider dark></v-divider>

                <v-list nav dense flat>
                    <template v-for="(item, i) in drawer.list">
                        <v-layout v-if="item.heading" :key="i">
                            <v-flex xs6>
                                <v-subheader v-if="item.heading">{{ item.heading }}</v-subheader>
                            </v-flex>
                        </v-layout>

                        <v-divider v-else-if="item.divider" :key="i" dark class="my-4"></v-divider>

                        <v-list-item
                                v-else
                                :key="i"
                                :to="item.target ? '' : item.link"
                                :href="item.target ? item.link : ''"
                                :target="item.target ? '_black' : ''"
                                :disabled="!item.link"
                                link
                                active-class="primary--text"
                        >
                            <v-list-item-action>
                                <v-icon v-text="item.icon"></v-icon>
                            </v-list-item-action>
                            <v-list-item-content>
                                <v-list-item-title>{{ item.text }}</v-list-item-title>
                            </v-list-item-content>
                        </v-list-item>
                    </template>
                </v-list>
            </v-sheet>
        </v-navigation-drawer>
    </nav>
</template>

<script>
    import LanguageSelector from "./selectors/LanguageSelector";
    import DarkModeSwitch from "./_partial/DarkModeSwitch";

    export default {
        components: {
            LanguageSelector,
            DarkModeSwitch
        },
        data:()=>({
            display: null,
        }),
        computed: {
            isLogin() {
                return this.$store.state.user.isLogin
            },
            username() {
                return this.$store.state.user.username
            },
            // menu
            drawer() {
                return {
                    list: [
                        {heading: this.$i18n.t('nav.menu.menu')},
                        {
                            icon: 'home',
                            text: this.$i18n.t('nav.menu.teams'),
                            link: '/'
                        },
                        {
                            icon: 'mdi-cloud-upload',
                            text: this.$i18n.t('nav.menu.upload'),
                            link: '/upload'
                        },
                        {
                            icon: 'mdi-card-search',
                            text: this.$i18n.t('nav.menu.search'),
                            link: '/search'
                        },
                        {
                            icon: 'mdi-chart-pie',
                            text: this.$i18n.t('nav.menu.usage'),
                            link: '/usage'
                        },
                        {divider: true},
                        {heading: this.$i18n.t('nav.extra.extra')},
                        {
                            icon: 'forum',
                            text: this.$i18n.t('nav.extra.forum'),
                            link: '/forum'
                        },
                        {
                            icon: 'mdi-information',
                            text: this.$i18n.t('nav.extra.about'),
                            link: '/about'
                        },
                        {divider: true},
                        {heading: this.$i18n.t('nav.external.external')},
                        {
                            icon: 'mdi-github',
                            text: 'GitHub',
                            link: 'https://github.com/txfs19260817/ScopeLens',
                            target: true
                        },
                    ]
                }
            },
            user() {
                return {
                    list: [
                        {
                            icon: 'mdi-text-box',
                            text: this.$i18n.t('nav.user.myteams'),
                            link: '/myteams',
                        },
                        {
                            icon: 'logout',
                            text: this.$i18n.t('nav.user.logout'),
                            link: '/logout',
                        },
                    ],
                }
            }
        }
    }
</script>
