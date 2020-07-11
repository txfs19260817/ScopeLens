<!-- https://gist.github.com/awunder/5ea26d7b7c6b1e3cf29ab265070478cb -->
<template>
    <v-container id="signinup-form" class="fill-height">
        <v-row align="center" justify="center" no-gutters>
            <v-col cols="12" sm="8" md="8" class="">
                <v-card class="evelation-12 card">
                    <v-window v-model="step">
                        <!--SignIn-->
                        <v-window-item :value="1">
                            <v-row class="">
                                <v-col cols="12" md="8" class="pt-6 pb-6">
                                    <v-card-text>
                                        <v-form class="signup-form-form" @submit.prevent="loginRequest">
                                            <h1 class="text-center display-1 mb-10" :class="`${bgColor}--text`">
                                                Sign in
                                            </h1>
                                            <v-text-field
                                                    id="username"
                                                    v-model="login.username"
                                                    label="Username"
                                                    name="Username"
                                                    append-icon="person"
                                                    type="text"
                                                    :color="bgColor"
                                                    required
                                            />
                                            <v-text-field
                                                    id="password"
                                                    v-model="login.password"
                                                    label="Password"
                                                    name="Password"
                                                    append-icon="lock"
                                                    type="password"
                                                    :color="bgColor"
                                                    required
                                            />
                                            <div class="text-center grey--text">
                                                Please contact administrator if you forgot password.
                                            </div>
                                            <div class="text-center mt-6">
                                                <v-btn type="submit" large dark :color="bgColor" :loading="loading">
                                                    Sign In
                                                </v-btn>
                                            </div>
                                        </v-form>
                                    </v-card-text>
                                </v-col>
                                <v-col cols="12" md="4" class="darken-2 vcenter" :class="`${bgColor}`">
                                    <div>
                                        <v-card-text :class="`${fgColor}--text`">
                                            <h1 class="text-center headline mb-3">No User?</h1>
                                            <h5 class="text-center overline mb-3">Please Sign Up to continue</h5>
                                        </v-card-text>
                                        <div class="text-center mb-6">
                                            <v-btn dark outlined @click="step = 2">Sign Up</v-btn>
                                        </div>
                                    </div>
                                </v-col>
                            </v-row>
                        </v-window-item>

                        <!--SignUp-->
                        <v-window-item :value="2">
                            <v-row class="fill-height">
                                <v-col cols="12" md="4" class="darken-2 vcenter" :class="`${bgColor}`">
                                    <div>
                                        <v-card-text :class="`${fgColor}--text`">
                                            <h1 class="text-center headline mb-3">Already a user?</h1>
                                            <h5 class="text-center overline mb-3">Please Sign In</h5>
                                        </v-card-text>
                                        <div class="text-center mb-6">
                                            <v-btn dark outlined @click="step = 1">Sign In</v-btn>
                                        </div>
                                    </div>
                                </v-col>
                                <v-col cols="12" md="8" class=" pt-6 pb-6">
                                    <v-card-text>
                                        <h1 class="text-center display-1 mb-10" :class="`${bgColor}--text`">Sign Up</h1>
                                        <ValidationObserver ref="observer" v-slot="{ validate }">
                                            <v-form class="signup-form-form" @submit.prevent="registerRequest">
                                                <ValidationProvider v-slot="{ errors }" name="username"
                                                                    rules="required|max:20">
                                                    <v-text-field
                                                            id="username"
                                                            v-model="register.username"
                                                            label="Username"
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
                                                            label="E-mail"
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
                                                            label="Password"
                                                            name="password"
                                                            append-icon="lock"
                                                            type="password"
                                                            required
                                                            :counter="20"
                                                            :error-messages="errors"
                                                    />
                                                </ValidationProvider>
                                                <div class="text-center mt-6">
                                                    <v-btn type="submit" large dark :color="bgColor" :loading="loading">
                                                        Sign Up
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
        },
        props: {
            bgColor: {
                type: String,
                default: 'blue'
            },
            fgColor: {
                type: String,
                default: 'white'
            }
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
        }),
        methods: {
            async registerRequest() {
                const v = await this.$refs.observer.validate()
                if (!v) return
                const success = await this.$store.dispatch('user/register', {data: this.register});
                if (success) {
                    this.step = 1;
                    this.login.username = this.register.username;
                    this.login.password = this.register.password;
                }
            },
            async loginRequest() {
                const success = await this.$store.dispatch('user/login', {data: this.login});
                if (success) {
                    await this.$router.push("/")
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
</style>
