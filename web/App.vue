<script>
    export default {
        data() {
            return {
                schemas: [],
                messages: [],
                alerts: [],
                form: {
                    type: null,
                    payload: '',
                    subject: '',
                },
                schemasFile: null,
            };
        },
        computed: {
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
            payloadExample() {
                let example = {};
                if (this.schemas.length > 0 && this.form.type) {
                    example = this.schemas.find(v => v.name === this.form.type).example;
                }
                return JSON.stringify(example, null, 2);
            },
            // last 5
            lastMessages() {
                return this.messages.reverse();
                // return this.messages.slice(this.messages.length - 5, this.messages.length).reverse();
            },
            lastAlerts() {
                const max = 3;
                const alerts = this.alerts.slice();
                const len = alerts.length;
                if (len <= 3) {
                    return alerts;
                }
                return alerts.slice(len - max, len).reverse();
            },

        },
        created() {
            this.loadSchemas();
            this.$backend.connectMessagesStream((message) => {
                this.messages.push(message);
            });
        },
        methods: {
            loadSchemas() {
                this.$backend.fetchSchemas().then(res => {
                    this.schemas = res;
                });
            },
            async processForm() {
                let payload;
                try {
                    payload = JSON.parse(this.form.payload);
                } catch (e) {
                    this.error = 'invalid JSON';
                    return;
                }
                try {
                    const res = await this.$backend.publishMessage({
                        subject: this.form.subject,
                        payload: payload,
                    });
                    this.showNotification(res.message);
                } catch (e) {
                    this.showError(e.message);
                }
            },
            showError(text) {
                this.alerts.push({type: 'error', text});
            },
            showNotification(text) {
                this.alerts.push({type: 'success', text});
            },
            async uploadSchemas() {
                try {
                    const res = await this.$backend.uploadSchema(this.schemasFile);
                    this.schemasFile = null;
                    this.loadSchemas();
                    this.showNotification(res.message);
                } catch (e) {
                    this.showError(e.message);
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
    <v-app>
        <v-app-bar app color="indigo" dark>
            <v-toolbar-title>NATS AVRO TESTER</v-toolbar-title>
        </v-app-bar>
        <v-main>
            <v-container fluid>
                <v-row>
                    <v-col>
                        <h1>Publish</h1>
                        <v-form @submit.prevent="processForm">
                            <v-autocomplete v-model="form.type" :items="schemas" item-text="namespace" item-value="name" dense filled label="Subject" no-data-text="No schemas found. Upload them first..."></v-autocomplete>
                            <v-textarea outlined label="Payload" :value="form.payload" v-model="form.payload" :auto-grow="true"></v-textarea>
                            <v-btn type=submit :disabled="!formValid" color="success" class="mr-4">Publish message
                            </v-btn>
                        </v-form>
                    </v-col>
                    <v-col>
                        <h1>Subscribe</h1>
                        <v-simple-table dense>
                            <template v-slot:default>
                                <thead>
                                <tr>
                                    <th class="text-left">ID</th>
                                    <th class="text-left">Subject</th>
                                    <th class="text-left">Payload</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lastMessages" :key="item.name">
                                    <td class="text-left">{{ item.id }}</td>
                                    <td class="text-left">{{ item.subject }}</td>
                                    <td class="text-left">{{ item.payload }}</td>
                                </tr>
                                </tbody>
                            </template>
                        </v-simple-table>
                    </v-col>
                </v-row>
                <v-row>
                    <v-col>
                        <h1>Upload schemas</h1>
                        <v-form @submit.prevent="uploadSchemas">
                            <v-file-input v-model="schemasFile" label="Schemas zip file" outlined dense></v-file-input>
                            <v-btn type=submit :disabled="schemasFile==null" clearable color="success" class="mr-4">Upload</v-btn>
                        </v-form>
                    </v-col>
                </v-row>
                <v-alert v-for="item in lastAlerts" dense outlined dismissible :type="item.type">{{ item.text }}</v-alert>
            </v-container>
        </v-main>
    </v-app>
</template>

<style>
    [v-cloak] {
        display: none;
    }
</style>
