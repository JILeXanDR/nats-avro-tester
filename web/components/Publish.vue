<script>
    export default {
        props: ['schemas'],
        data() {
            return {
                form: {
                    type: null,
                    payload: '',
                    subject: '',
                },
            };
        },
        computed: {
            payloadExample() {
                let example = {};
                if (this.form.type) {
                    example = this.schemas.find(v => v.name === this.form.type).example;
                }
                return JSON.stringify(example, null, 2);
            },
            isPayloadValid() {
                try {
                    JSON.parse(this.form.payload);
                    return true;
                } catch (e) {
                    return false;
                }
            },
            formValid() {
                return this.isPayloadValid && this.form.subject.length > 0;
            },
        },
        methods: {
            async processForm() {
                let payload;
                try {
                    payload = JSON.parse(this.form.payload);
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
            <v-textarea outlined label="Payload" :value="form.payload" v-model="form.payload" :auto-grow="true"></v-textarea>
            <v-btn type=submit :disabled="!formValid" color="success" class="mr-4">Publish message</v-btn>
        </v-form>
    </div>
</template>
