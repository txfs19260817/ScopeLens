<template>
    <ValidationProvider v-slot="{ errors }" name="Format" :rules="`${required ? 'required' : ''}`">
        <v-autocomplete
                id="format"
                v-model="valueModel"
                :label="required? `*`+$t('formatSelector.format'):$t('formatSelector.format')"
                :items="formats"
                persistent-hint
                :hint="hint"
                menu-props="auto"
                outlined
                :error-messages="errors"
        ></v-autocomplete>
    </ValidationProvider>
</template>

<script>
    import {required} from 'vee-validate/dist/rules'
    import {extend, ValidationProvider, setInteractionMode} from 'vee-validate'
    import {formats} from "../../assets/formats";

    setInteractionMode('eager')
    extend('required', {
        ...required,
        message: '{_field_} can not be empty',
    });

    export default {
        name: "FormatSelector",
        components: {
            ValidationProvider
        },
        props: {
            // v-model
            value: {
                required: true,
            },
            hint: {
                type: String,
                default: 'You can type words here to search for desired format. ',
            },
            required: {
                type: Boolean,
                default: false,
            }
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
            formats() {
                return formats
            },
        }
    }
</script>

<style scoped>

</style>