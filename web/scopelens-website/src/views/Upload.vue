<template>
    <v-container id="upload-form" class="fill-height">
        <v-row align="center" justify="center" no-gutters>
            <v-col>
                <v-card class="elevation-2 card">
                    <v-card-text>
                        <h1 class="text-start display-1 mb-10 fg-text"> Upload your team </h1>
                        <ValidationObserver ref="observer" v-slot="{ validate }">
                            <v-form class="upload-form-form" @submit.prevent="submit">
                                <ValidationProvider v-slot="{ errors }" name="Title" rules="required|max:50">
                                    <v-text-field
                                            id="title"
                                            v-model="form.title"
                                            label="*Title"
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
                                <v-switch
                                        v-model="notAuthor"
                                        class="shrink mr-2 mt-0"
                                        label="I am not the team author. "
                                        :hint="hint.authorSwitch"
                                        persistent-hint
                                ></v-switch>
                                <ValidationProvider v-slot="{ errors }" name="Author" :rules="`${notAuthor ? 'required|max:20' : ''}`">
                                    <v-text-field
                                            v-if="notAuthor"
                                            id="author"
                                            v-model="form.author"
                                            label="*Author"
                                            persistent-hint
                                            :hint="hint.author"
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
                                <FormatSelector :value.sync="form.format" :hint="hint.format" :required="true">
                                </FormatSelector>
                                <v-switch
                                        v-model="haveShowdown"
                                        class="shrink mr-2 mt-0"
                                        label="I have the Showdown paste. "
                                        :hint="hint.showdownSwitch"
                                        persistent-hint
                                ></v-switch>
                                <ValidationProvider v-if="haveShowdown" v-slot="{ errors }" name="Showdown" :rules="`${haveShowdown ? 'required|max:1600' : ''}`">
                                    <v-textarea
                                            id="showdown"
                                            v-model="form.showdown"
                                            label="*Showdown"
                                            outlined
                                            :hint="hint.showdown"
                                            :clearable="true"
                                            :counter="1600"
                                            :error-messages="errors"
                                    ></v-textarea>
                                </ValidationProvider>
                                <PokemonSelector v-else :value.sync="form.pokemon" :hint="hint.pokemon" :required="!haveShowdown">
                                </PokemonSelector>
                                <v-file-input
                                        ref="image"
                                        v-model="imageFile"
                                        show-size
                                        outlined
                                        persistent-hint
                                        :hint="hint.image"
                                        :rules="fileRules"
                                        accept="image/png, image/jpeg, image/jpg"
                                        placeholder="Pick a rental team preview photo. Only accept .png/.jpg/.jpeg format."
                                        prepend-icon=""
                                        append-icon="mdi-camera"
                                        label="Image"
                                ></v-file-input>
                                <ValidationProvider v-slot="{ errors }" name="Description" rules="max:2800">
                                    <v-textarea
                                            id="description"
                                            v-model="form.description"
                                            label="Description"
                                            outlined
                                            :auto-grow="true"
                                            :clearable="true"
                                            :counter="2800"
                                    ></v-textarea>
                                </ValidationProvider>
                                <div class="text-center mt-6">
                                    <v-btn type="submit" color="primary" large dark :loading="loading">
                                        <v-icon left dark>mdi-upload</v-icon>
                                        Submit
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
    import {ERROR} from "../api";
    import {forme} from "../assets/pokemonNames"

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
                fileRules: [
                    value => !value || value.size < 2000000 || 'Photo size should be less than 2 MB!',
                    type => !type || ['image/png', 'image/jpeg', 'image/jpg'].includes(type.type) || 'Only accept .png, .jpg or .jpeg image file!',
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
                    state: 1
                },
                // hints
                hint: {
                    author: 'Please fill the name of team author here. ',
                    authorSwitch: 'Please turn on the switch if you are not the team author. ',
                    format: 'You can type words here to search for desired format. ',
                    showdownSwitch:'Please select all team members if no Showdown paste provided. ',
                    showdown: 'Please paste the Showdown team here. ',
                    pokemon: 'Please select up to 6 Pokemon. You can type words here to filter Pokemon.',
                    image: 'Optional. **Note that we only accept image with available rental ID.**',
                },
            }
        },
        computed: {
            loading() {
                return this.$store.state.loading.loading
            },
        },
        methods: {
            async submit() {
                // loading
                this.$store.commit('LOADING_ON')
                // validate
                const v = await this.$refs.observer.validate()
                if (v) {
                    // form preparation
                    // convert image to base64 string
                    if(this.imageFile !== undefined) this.form.image = await toBase64(this.imageFile)
                    // assign uploader
                    this.form.uploader = this.$store.state.user.username
                    if (!this.notAuthor) this.form.author = this.form.uploader
                    // Auto push Pokemon names to form if Showdown text is provided
                    if (this.haveShowdown) {
                        this.form.pokemon = []
                        for (const p of Koffing.parse(this.form.showdown).teams[0].pokemon) {
                            if (forme[p.name] === undefined){
                                // no alter forme, push to form directly
                                this.form.pokemon.push(p.name)
                            } else {
                                // push the origin species name
                                this.form.pokemon.push(forme[p.name])
                            }
                        }
                    } else {
                        this.form.showdown = ""
                        // process Pokemon names ('A/B/C' --> 'A')
                        this.form.pokemon.forEach((item, idx) => this.form.pokemon[idx] = item.toString().split('/', 1)[0])
                    }

                    const res = await insertTeam(this.form, this.$store.state.user.token);
                    if (res.data.code === ERROR || res.status === 401) {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": "Upload team error: " + res.data.msg,
                            "color": "error"
                        });
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": "Upload team success!",
                            "color": "success"
                        });
                        await this.$router.push("/")
                    }
                }
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

    .fg-text {
        color: #4768A1;
    }

</style>