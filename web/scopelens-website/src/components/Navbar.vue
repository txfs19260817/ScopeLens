<template>
    <nav>
        <v-app-bar v-if="appbar.tabs" app fixed flat color="bg_primary">
            <v-app-bar-nav-icon @click="drawer.display = !drawer.display" class="d-lg-none"></v-app-bar-nav-icon>
        </v-app-bar>

        <v-navigation-drawer v-model="drawer.display" color="bg_secondary" app fixed left flat>
            <v-sheet color="bg_secondary">
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
                            <v-list-item-title class="title">{{ 'Login/Register' }}</v-list-item-title>
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
                                active-class="amber--text"
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
    export default {
        data: () => ({
            appbar: {
                tabs: [
                    'Posts',
                    'Trending'
                ]
            },
            drawer: {
                display: null,
                list: [
                    {heading: 'Menu'},
                    {
                        icon: 'home',
                        text: 'Teams',
                        link: '/'
                    },
                    {
                        icon: 'mdi-cloud-upload',
                        text: 'Upload',
                        link: '/upload'
                    },
                    {
                        icon: 'mdi-card-search',
                        text: 'Search',
                        link: '/search'
                    },
                    {
                        icon: 'mdi-chart-pie',
                        text: 'Usage',
                        link: '/usage'
                    },
                    {divider: true},
                    {heading: 'Extra'},
                    {
                        icon: 'forum',
                        text: 'Forum',
                        link: '/Forum'
                    },
                    {
                        icon: 'mdi-information',
                        text: 'About',
                        link: '/About'
                    },
                    {divider: true},
                    {heading: 'Navigation'},
                    {
                        icon: 'mdi-github',
                        text: 'GitHub',
                        link: 'https://github.com/txfs19260817/ScopeLens',
                        target: true
                    },
                    // {
                    //     icon: 'mdi-twitter',
                    //     text: 'Twitter',
                    //     link: 'https://twitter.com/stackFlam1ngo',
                    //     target: true
                    // }
                ]
            },
            user: {
                list: [
                    {
                        icon: 'mdi-text-box',
                        text: 'My Teams',
                        link: '/myteams',
                    },
                    {
                        icon: 'logout',
                        text: 'Logout',
                        link: '/logout',
                    },
                ],
            },
        }),
        computed: {
            isLogin() {
                return this.$store.state.user.isLogin
            },
            username() {
                return this.$store.state.user.username
            }
        }
    }
</script>
