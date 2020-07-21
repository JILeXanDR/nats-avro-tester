<script>
    import VJsoneditor from 'v-jsoneditor/src/index'

    export default {
        components: {
            VJsoneditor,
        },
        props: ['schemas'],
        data() {
            return {
                isPayloadValid: false,
                form: {
                    type: null,
                    payload: '',
                    subject: '',
                },
                options: {
                    mode: 'code',
                    mainMenuBar: false,
                    onValidationError: (errs) => {
                        this.isPayloadValid = errs.length === 0;
                    },
                },
            };
        },
        computed: {
            payloadExample() {
                let example = {};
                if (this.form.type) {
                    example = this.schemas.find(v => v.name === this.form.type).example;
                }
                return example;
            },
            formValid() {
                return this.isPayloadValid && this.form.subject.length > 0;
            },
            editorHeight() {
                // TODO: optimize
                const json = JSON.stringify(this.form.payload, null, '  ');
                let lines = json.split('\n').length;
                lines = lines < 8 ? 8 : lines;
                return `${lines * 19}px`;
            },
        },
        methods: {
            async processForm() {
                let payload;
                try {
                    payload = this.form.payload;
                } catch (e) {
                    return;
                }
                try {
                    const res = await this.$backend.publishMessage({
                        subject: this.form.subject,
                        payload: payload,
                    });
                    this.$emit('success', 'Message is sent.');
                } catch (e) {
                    this.$emit('error', e.message);
                }
            },
        },
        watch: {
            payloadExample: {
                immediate: true,
                handler(val) {
                    this.form.payload = val;
                },
            },
            'form.type': {
                handler(val) {
                    const schema = this.schemas.find(v => v.name === val);
                    this.form.subject = schema ? schema.namespace : '';
                },
            },
        },
    }
</script>

<template>
    <div>
        <v-form @submit.prevent="processForm">
            <v-autocomplete v-model="form.type" :items="schemas" item-text="namespace" item-value="name" dense filled label="Subject" no-data-text="No schemas found"></v-autocomplete>
            <v-row>
                <v-col>
                    <v-jsoneditor ref="editor" v-model="form.payload" :options="options" :height="editorHeight" :plus="false"></v-jsoneditor>
                </v-col>
            </v-row>
            <v-btn type=submit :disabled="!formValid" color="success" class="mr-4">Publish message</v-btn>
        </v-form>
    </div>
</template>
