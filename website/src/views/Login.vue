<!-- https://gist.github.com/awunder/5ea26d7b7c6b1e3cf29ab265070478cb -->
<template>
    <v-container id="signinup-form" class="fill-height">
        <v-row align="center" justify="center" no-gutters>
            <v-col cols="12" sm="8" md="8" class="">
                <v-card class="elevation-12 card">
                    <v-window v-model="step">
                        <!--SignIn-->
                        <v-window-item :value="1">
                            <v-row class="">
                                <v-col cols="12" md="8" class="pt-6 pb-6">
                                    <v-card-text>
                                        <v-form class="signup-form-form" @submit.prevent="loginRequest">
                                            <h1 class="text-center display-1 mb-10 primary--text">{{ $t('login.signin') }}</h1>
                                            <v-text-field
                                                    id="username"
                                                    v-model="login.username"
                                                    :label="$t('login.username')"
                                                    name="Username"
                                                    append-icon="person"
                                                    type="text"
                                                    required
                                            />
                                            <v-text-field
                                                    id="password"
                                                    v-model="login.password"
                                                    :label="$t('login.password')"
                                                    name="Password"
                                                    append-icon="lock"
                                                    type="password"
                                                    required
                                            />
                                            <div class="text-center grey--text">
                                                {{ $t('login.forget') }}
                                            </div>
                                            <div class="text-center mt-6">
                                                <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}" type="submit" large dark color="primary" :loading="loading">
                                                    {{ $t('login.signin') }}
                                                </v-btn>
                                            </div>
                                        </v-form>
                                    </v-card-text>
                                </v-col>
                                <v-col cols="12" md="4" class="darken-1 vcenter primary">
                                    <div>
                                        <v-card-text :class="{'black--text': $vuetify.theme.dark, 'white--text': !$vuetify.theme.dark}">
                                            <h1 class="text-center headline mb-3">{{ $t('login.noUser') }}</h1>
                                            <h5 class="text-center overline mb-3">{{ $t('login.goSignup') }}</h5>
                                        </v-card-text>
                                        <div class="text-center mb-6">
                                            <v-btn :class="{'black--text': $vuetify.theme.dark, 'white--text': !$vuetify.theme.dark}" dark outlined @click="step = 2">
                                                {{ $t('login.signup') }}
                                            </v-btn>
                                        </div>
                                    </div>
                                </v-col>
                            </v-row>
                        </v-window-item>

                        <!--SignUp-->
                        <v-window-item :value="2">
                            <v-row class="fill-height">
                                <v-col cols="12" md="4" class="darken-1 vcenter primary">
                                    <div>
                                        <v-card-text :class="{'black--text': $vuetify.theme.dark, 'white--text': !$vuetify.theme.dark}">
                                            <h1 class="text-center headline mb-3">{{ $t('login.already') }}</h1>
                                            <h5 class="text-center overline mb-3">{{ $t('login.goSignin') }}</h5>
                                        </v-card-text>
                                        <div class="text-center mb-6">
                                            <v-btn :class="{'black--text': $vuetify.theme.dark, 'white--text': !$vuetify.theme.dark}" dark outlined @click="step = 1">
                                                {{ $t('login.signin') }}
                                            </v-btn>
                                        </div>
                                    </div>
                                </v-col>
                                <v-col cols="12" md="8" class=" pt-6 pb-6">
                                    <v-card-text>
                                        <h1 class="text-center display-1 mb-10 primary--text">{{ $t('login.signup') }}</h1>
                                        <ValidationObserver ref="observer" v-slot="{ validate }">
                                            <v-form class="signup-form-form" @submit.prevent="registerRequest">
                                                <ValidationProvider v-slot="{ errors }" name="username"
                                                                    rules="required|max:20">
                                                    <v-text-field
                                                            id="username"
                                                            v-model="register.username"
                                                            :label="$t('login.username')"
                                                            name="username"
                                                            append-icon="person"
                                                            type="text"
                                                            required
                                                            :counter="20"
                                                            :error-messages="errors"
                                                    />
                                                </ValidationProvider>
                                                <ValidationProvider v-slot="{ errors }" name="email"
                                                                    rules="required|email">
                                                    <v-text-field
                                                            id="email"
                                                            v-model="register.email"
                                                            :label="$t('login.email')"
                                                            name="email"
                                                            append-icon="email"
                                                            type="email"
                                                            required
                                                            :error-messages="errors"
                                                    />
                                                </ValidationProvider>
                                                <ValidationProvider v-slot="{ errors }" name="password"
                                                                    rules="required|max:20">
                                                    <v-text-field
                                                            id="password"
                                                            v-model="register.password"
                                                            :label="$t('login.password')"
                                                            name="password"
                                                            append-icon="lock"
                                                            type="password"
                                                            required
                                                            :counter="20"
                                                            :error-messages="errors"
                                                    />
                                                </ValidationProvider>
                                                <vue-recaptcha :sitekey="siteKey" @verify="onVerify"></vue-recaptcha>
                                                <div class="text-center mt-6">
                                                    <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}" type="submit" large dark color="primary" :loading="loading">
                                                        {{ $t('login.signup') }}
                                                    </v-btn>
                                                </div>
                                            </v-form>
                                        </ValidationObserver>
                                    </v-card-text>
                                </v-col>
                            </v-row>
                        </v-window-item>
                    </v-window>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import {required, email, max} from 'vee-validate/dist/rules'
    import {extend, ValidationObserver, ValidationProvider, setInteractionMode} from 'vee-validate'
    import VueRecaptcha from 'vue-recaptcha'

    setInteractionMode('eager')
    extend('required', {
        ...required,
        message: '{_field_} can not be empty',
    })
    extend('max', {
        ...max,
        message: '{_field_} may not be greater than {length} characters',
    })
    extend('email', {
        ...email,
        message: 'Email must be valid',
    })

    export default {
        name: 'Login',
        components: {
            ValidationProvider,
            ValidationObserver,
            VueRecaptcha
        },
        data: () => ({
            // Active window
            step: 1,
            // Forms
            login: {
                username: "",
                password: ""
            },
            register: {
                username: "",
                email: "",
                password: ""
            },
            // reCaptcha
            siteKey: process.env.VUE_APP_RECAPTCHA,
            reCaptchaResponse: ""
        }),
        methods: {
            onVerify: function (response) {
                this.reCaptchaResponse = response
            },
            async registerRequest() {
                const v = await this.$refs.observer.validate()
                if (!v) return
                const success = await this.$store.dispatch('user/register', {data: this.register, recaptcha: this.reCaptchaResponse});
                if (success) {
                    // back to login window
                    this.step = 1;
                    this.login.username = this.register.username;
                    this.login.password = this.register.password;
                }
            },
            async loginRequest() {
                const success = await this.$store.dispatch('user/login', {data: this.login});
                if (success) {
                    // back to last page
                    await this.$router.go(-1)
                }
            }
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            }
        }
    }
</script>

<style scoped lang="scss">
    .v-input__icon--double .v-input__icon {
        margin-left: -4.25rem !important;
    }

    a.no-text-decoration {
        text-decoration: none;
    }

    #signinup-form {
        max-width: 75rem;
    }

    .signup-form-form {
        max-width: 23rem;
        margin: 0 auto;
    }

    .card {
        overflow: hidden;
    }

    .vcenter {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .bg-text {
        color: black;
    }
</style>
