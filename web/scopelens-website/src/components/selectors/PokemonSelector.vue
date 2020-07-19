<template>
    <ValidationProvider v-slot="{ errors }" name="Pokemon" :rules="`${required ? 'required|minmax:1,6' : ''}`">
        <v-autocomplete
                v-model="valueModel"
                :items="pokemon"
                outlined
                chips
                :label="required? `*`+$t('pokemonSelector.pokemon'):$t('pokemonSelector.pokemon')"
                persistent-hint
                :hint="hint"
                item-text="name"
                item-value="name"
                multiple
                :counter="6"
                :error-messages="errors"
        >
            <template v-slot:selection="data">
                <v-chip
                        v-bind="data.attrs"
                        :input-value="data.selected"
                        close
                        @click="data.select"
                        @click:close="removePokemon(data.item)"
                >
                    <v-avatar left>
                        <v-img :src="iconUrl + data.item.avatar"></v-img>
                    </v-avatar>
                    {{ data.item.name }}
                </v-chip>
            </template>
            <template v-slot:item="data">
                <template v-if="typeof data.item !== 'object'">
                    <v-list-item-content v-text="data.item"></v-list-item-content>
                </template>
                <template v-else>
                    <v-list-item-avatar>
                        <img :src="iconUrl + data.item.avatar">
                    </v-list-item-avatar>
                    <v-list-item-content>
                        <v-list-item-title v-html="data.item.name"></v-list-item-title>
                        <v-list-item-subtitle
                                v-html="data.item.group"></v-list-item-subtitle>
                    </v-list-item-content>
                </template>
            </template>
        </v-autocomplete>
    </ValidationProvider>
</template>

<script>
    import {pmNames4Select} from "../../assets/pokemonNames";
    import {required} from 'vee-validate/dist/rules'
    import {extend, ValidationProvider, setInteractionMode} from 'vee-validate'

    setInteractionMode('eager')
    extend('required', {
        ...required,
        message: '{_field_} can not be empty',
    });

    extend('minmax', {
        validate(value, {min, max}) {
            return value.length >= min && value.length <= max;
        },
        params: ['min', 'max'],
        message: 'The {_field_} field must have at least {min} items and {max} items at most',
    });

    export default {
        name: "PokemonSelector",
        components: {
            ValidationProvider
        },
        props: {
            // v-model
            value: {
                type: Array,
                required: true,
            },
            hint: {
                type: String,
                default: 'Please select up to 6 Pokemon. You can type words here to filter Pokemon.',
            },
            required: {
                type: Boolean,
                default: false,
            }
        },
        data() {
            return {
                // Sprites paths
                iconUrl: process.env.VUE_APP_STATIC_ASSET_URL + '2d/',
            }
        },
        methods: {
            removePokemon(item) {
                const index = this.valueModel.indexOf(item.name)
                if (index >= 0) this.valueModel.splice(index, 1)
            },
        },
        computed: {
            valueModel: {
                get() {
                    return this.value
                },
                set(v) {
                    this.$emit('update:value', v)
                },
            },
            pokemon() {
                return pmNames4Select
            },
        }
    }
</script>

<style scoped>

</style>