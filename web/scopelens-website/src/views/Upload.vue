<template>
    <v-container id="upload-form" class="fill-height">
        <v-row align="center" justify="center" no-gutters>
            <v-col>
                <v-card class="elevation-6 card">
                    <v-card-text>
                        <h1 class="text-start display-1 mb-10" :class="`${bgColor}--text`"> Upload your team </h1>
                        <ValidationObserver ref="observer" v-slot="{ validate }">
                            <v-form class="upload-form-form" @submit.prevent="loginRequest">
                                <ValidationProvider v-slot="{ errors }" name="Title" rules="required|max:50">
                                    <v-text-field
                                            id="title"
                                            v-model="form.title"
                                            label="Title"
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
                                <ValidationProvider v-slot="{ errors }" name="Author" rules="max:50">
                                    <v-switch
                                            v-model="sameAsUploader"
                                            class="shrink mr-2 mt-0"
                                            label="I am not the team author. "
                                    ></v-switch>
                                    <v-text-field
                                            v-if="sameAsUploader"
                                            id="author"
                                            v-model="form.author"
                                            label="Author"
                                            :hint="hint.author"
                                            name="Author"
                                            append-icon="person"
                                            type="text"
                                            outlined
                                            :clearable="true"
                                            :color="bgColor"
                                            :counter="50"
                                            :error-messages="errors"
                                    />
                                </ValidationProvider>
                                <ValidationProvider v-slot="{ errors }" name="format" rules="required">
                                    <v-autocomplete
                                            id="format"
                                            v-model="form.format"
                                            label="Format"
                                            :items="formats"
                                            persistent-hint
                                            :hint="hint.format"
                                            menu-props="auto"
                                            outlined
                                            :error-messages="errors"
                                    ></v-autocomplete>
                                </ValidationProvider>
                                <ValidationProvider v-slot="{ errors }" name="Showdown" rules="required|max:1600">
                                    <v-textarea
                                            id="showdown"
                                            v-model="form.showdown"
                                            label="Showdown"
                                            outlined
                                            persistent-hint
                                            :hint="hint.showdown"
                                            :auto-grow="true"
                                            :clearable="true"
                                            :counter="1600"
                                    ></v-textarea>
                                </ValidationProvider>
                                <ValidationProvider v-slot="{ errors }" name="Description" rules="required|max:2800">
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
                            </v-form>
                        </ValidationObserver>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import {formats} from "../assets/formats"
    import {pmNames4Select} from "../assets/pokemonNames"
    import {required, max} from 'vee-validate/dist/rules'
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

    export default {
        name: "Upload",
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
        data() {
            return {
                sameAsUploader: true,
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
                // hint
                hint: {
                    author: 'Please fill the name of the team author here. ',
                    format: 'You can type words here to search for desired format. ',
                    pokemon: 'Please select up to 6 Pokemon. You can type words here to filter Pokemon',
                    showdown: 'Please paste the Showdown team here (if applicable). ',
                },
            }
        },
        computed: {
            formats() {
                return formats
            },
            pokemonNames() {
                return pmNames4Select
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