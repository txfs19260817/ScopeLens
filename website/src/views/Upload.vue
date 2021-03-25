<template>
    <v-container id="upload-form" class="fill-height">
        <v-row align="center" justify="center" no-gutters>
            <v-col>
                <v-card class="elevation-2 card">
                    <v-card-text>
                        <h1 class="text-start display-1 mb-10 primary--text"> {{$t('upload.title')}} </h1>
                        <ValidationObserver ref="observer" v-slot="{ validate }">
                            <v-form class="upload-form-form" @submit.prevent="submit">
<!--                                Title-->
                                <ValidationProvider v-slot="{ errors }" name="Title" rules="required|max:50">
                                    <v-text-field
                                            id="title"
                                            v-model="form.title"
                                            :label="$t('upload.form.title')"
                                            name="Title"
                                            append-icon="mdi-pencil"
                                            type="text"
                                            outlined
                                            :clearable="true"
                                            :color="bgColor"
                                            required
                                            :counter="50"
                                            :error-messages="errors"
                                    />
                                </ValidationProvider>
<!--                                NotAuthor switch-->
                                <v-switch
                                        v-model="notAuthor"
                                        class="shrink mr-2 mt-0"
                                        :label="$t('upload.form.authorSwitch')"
                                        :hint="$t('upload.hint.authorSwitch')"
                                        persistent-hint
                                ></v-switch>
<!--                                Author-->
                                <ValidationProvider v-slot="{ errors }" name="Author" :rules="`${notAuthor ? 'required|max:20' : ''}`">
                                    <v-text-field
                                            v-if="notAuthor"
                                            id="author"
                                            v-model="form.author"
                                            :label="$t('upload.form.author')"
                                            persistent-hint
                                            :hint="$t('upload.hint.author')"
                                            name="Author"
                                            append-icon="person"
                                            type="text"
                                            outlined
                                            :clearable="true"
                                            :color="bgColor"
                                            :counter="20"
                                            :error-messages="errors"
                                    />
                                </ValidationProvider>
<!--                                Format-->
                                <FormatSelector :value.sync="form.format" :hint="$t('upload.hint.format')" :required="true">
                                </FormatSelector>
                                <v-file-input
                                        ref="image"
                                        v-model="imageFile"
                                        show-size
                                        outlined
                                        persistent-hint
                                        :hint="$t('upload.hint.image')"
                                        :rules="fileRules"
                                        accept="image/png, image/jpeg, image/jpg"
                                        :placeholder="$t('upload.form.image')"
                                        prepend-icon=""
                                        append-icon="mdi-camera"
                                        label="Image"
                                ></v-file-input>
<!--                                haveShowdown switch-->
                                <v-switch
                                        v-model="haveShowdown"
                                        class="shrink mr-2 mt-0"
                                        :label="$t('upload.form.showdownSwitch')"
                                        :hint="$t('upload.hint.showdownSwitch')"
                                        persistent-hint
                                ></v-switch>
<!--                                Showdown-->
                                <ValidationProvider v-if="haveShowdown" v-slot="{ errors }" name="Showdown" :rules="`${haveShowdown ? 'required|max:1600' : ''}`">
                                    <v-textarea
                                            id="showdown"
                                            v-model="form.showdown"
                                            label="*Showdown"
                                            outlined
                                            :hint="$t('upload.hint.showdown')"
                                            :clearable="true"
                                            :counter="1600"
                                            :error-messages="errors"
                                    ></v-textarea>
                                </ValidationProvider>
<!--                                Pokemon-->
                                <v-row v-if="!haveShowdown">
                                    <v-col cols="10">
                                        <PokemonSelector :value.sync="form.pokemon" :hint="$t('upload.hint.pokemon')" :required="!haveShowdown">
                                        </PokemonSelector>
                                    </v-col>
                                    <v-col cols="1">
                                        <v-tooltip top close-delay="3000">
                                            <template v-slot:activator="{ on, attrs }">
                                                <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}"
                                                       @click="recognize"
                                                       color="primary"
                                                       v-bind="attrs"
                                                       v-on="on"
                                                       fab
                                                       dark
                                                       :loading="loading">
                                                    <v-icon dark>mdi-face-recognition</v-icon>
                                                </v-btn>
                                            </template>
                                            <span>{{ $t("upload.btn.recognition") }}</span>
                                        </v-tooltip>
                                    </v-col>
                                </v-row>
<!--                                Description-->
                                <ValidationProvider v-slot="{ errors }" name="Description" rules="max:2800">
                                    <v-textarea
                                            id="description"
                                            v-model="form.description"
                                            :label="$t('upload.form.description')"
                                            outlined
                                            :auto-grow="true"
                                            :clearable="true"
                                            :counter="2800"
                                    ></v-textarea>
                                </ValidationProvider>
<!--                                Submit button-->
                                <div class="text-center mt-6">
                                    <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}" type="submit" color="primary" large dark :loading="loading">
                                        <v-icon left dark>mdi-upload</v-icon>
                                        {{ $t('upload.btn.submit') }}
                                    </v-btn>
                                </div>
                            </v-form>
                        </ValidationObserver>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import FormatSelector from "../components/selectors/FormatSelector";
    import PokemonSelector from "../components/selectors/PokemonSelector";
    import {Koffing} from "koffing"
    import {toBase64} from "../assets/utils"
    import {required, max} from 'vee-validate/dist/rules'
    import {extend, ValidationObserver, ValidationProvider, setInteractionMode} from 'vee-validate'
    import {insertTeam} from "../api/team";
    import {getSpeciesByImage} from "../api/predict";
    import {ERROR} from "../api";
    import {forme, en2pmNames} from "../assets/pokemonNames"

    setInteractionMode('eager')
    extend('required', {
        ...required,
        message: '{_field_} can not be empty',
    });
    extend('max', {
        ...max,
        message: '{_field_} may not be greater than {length} characters',
    });


    export default {
        name: "Upload",
        components: {
            ValidationProvider,
            ValidationObserver,
            FormatSelector,
            PokemonSelector
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
        data() {
            return {
                // author field will be autofilled as username when switch on.
                notAuthor: true,
                haveShowdown: true,
                // image upload rules
                imageTypes: ['image/png', 'image/jpeg', 'image/jpg'],
                imageMaxSize: 2000000,
                fileRules: [
                    value => !value || value.size <= this.imageMaxSize || 'Photo size should be less than 2 MB!',
                    type => !type || this.imageTypes.includes(type.type) || 'Only accept .png, .jpg or .jpeg image file!',
                ],
                imageFile: undefined,
                // form
                form: {
                    title: '',
                    author: '',
                    format: '',
                    pokemon: [],
                    showdown: '',
                    image: '',
                    description: '',
                    uploader: '',
                    state: 1,
                    has_showdown: false,
                    has_rental: false
                },
            }
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
        },
        methods: {
            async recognize() {
                // loading on
                this.$store.commit('LOADING_ON')
                // validate
                let v = false
                if (this.imageFile === undefined) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('upload.msg.recognize.empty'),
                        "color": "error"
                    });
                }
                else if (!(this.imageFile instanceof File) || !(this.imageTypes.includes(this.imageFile.type))) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('upload.msg.recognize.wrongType'),
                        "color": "error"
                    });
                }
                else if (this.imageFile.size > this.imageMaxSize) {
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('upload.msg.recognize.wrongType'),
                        "color": "error"
                    });
                } else {
                    v = true
                }
                // call api
                if (v) {
                    const res = await getSpeciesByImage(this.imageFile)
                    if (res.status !== 200) {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('upload.msg.recognize.failed'),
                            "color": "error"
                        });
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('upload.msg.recognize.success'),
                            "color": "success"
                        });
                    }
                    this.form.pokemon = res.data.names.map(x => en2pmNames[x])
                }
                // loading off
                this.$store.commit('LOADING_OFF')
            },
            async submit() {
                // loading on
                this.$store.commit('LOADING_ON')
                // validate
                const v = await this.$refs.observer.validate()
                if (v) {
                    // form preparation
                    // 1. Convert image to base64 string
                    if (this.imageFile !== undefined) {
                        this.form.image = await toBase64(this.imageFile)
                        this.form.has_rental = true
                    } else {
                        this.form.has_rental = false
                    }
                    // 2. Assign uploader
                    this.form.uploader = this.$store.state.user.username
                    if (!this.notAuthor) this.form.author = this.form.uploader
                    // 3. Auto push Pokemon names to form if Showdown text is provided
                    if (this.haveShowdown) {
                        this.form.pokemon = [];
                        this.form.has_showdown = true;
                        for (const p of Koffing.parse(this.form.showdown).teams[0].pokemon) {
                            if (forme[p.name] === undefined){
                                // no alter forme, push it to form directly
                                this.form.pokemon.push(p.name)
                            } else {
                                // push the origin species name to form
                                this.form.pokemon.push(forme[p.name])
                            }
                        }
                    } else {
                        this.form.showdown = "";
                        this.form.has_showdown = false;
                        // process Pokemon names ('A/B/C' --> 'A')
                        this.form.pokemon.forEach((item, idx) => this.form.pokemon[idx] = item.toString().split('/', 1)[0])
                    }
                    // call api
                    const res = await insertTeam(this.form, this.$store.state.user.token);
                    if (res.data.code === ERROR || res.status === 401) {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('upload.msg.failed') + res.data.msg,
                            "color": "error"
                        });
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('upload.msg.success'),
                            "color": "success"
                        });
                        await this.$router.push("/")
                    }
                }
                // loading off
                this.$store.commit('LOADING_OFF')
            },
        }
    }
</script>

<style scoped lang="scss">
    a.no-text-decoration {
        text-decoration: none;
    }

    #upload-form {
        max-width: 70rem;
    }

    .upload-form-form {
        max-width: 40rem;
        margin: 0 auto;
    }

    .card {
        overflow: hidden;
    }
</style>