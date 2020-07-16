<script>
    import Publish from './components/Publish';
    import Subscribe from './components/Subscribe';
    import UploadSchemas from './components/UploadSchemas';

    export default {
        components: {Publish, Subscribe, UploadSchemas},
        data() {
            return {
                tab: null,
                schemas: [],
                messages: [],
                readMessagesCount: 0,
                snackbar: {
                    model: false,
                    timeout: 5000,
                    color: 'info',
                    text: '',
                },
            };
        },
        computed: {
            messagesCount() {
                return this.messages.length;
            },
            schemasCount() {
                return this.schemas.length;
            },
            unreadMessagesCount() {
                return this.messagesCount - this.readMessagesCount;
            },
        },
        created() {
            this.loadSchemas();
            this.$backend.connectMessagesStream(message => this.messages.push(message));
        },
        methods: {
            async loadSchemas() {
                try {
                    this.$backend.fetchSchemas().then(res => this.schemas = res);
                } catch (e) {
                    this.showError(e.message);
                }
            },
            showNotification(type, text) {
                this.snackbar.model = true;
                this.snackbar.color = type;
                this.snackbar.text = text;
            },
            showError(text) {
                this.showNotification('error', text);
            },
            showSuccess(text) {
                this.showNotification('success', text);
            },
            onSuccess(message) {
                this.showSuccess(message);
            },
            onError(message) {
                this.showError(message);
            },
            onUploadedSchemasSuccess(message) {
                this.showNotification(message);
                this.loadSchemas();
            },
        },
        watch: {
            messages(v) {
                // tab with messages is active now
                if (this.tab === 1) {
                    this.readMessagesCount = v.length;
                }
            },
            tab(v) {
                if (v === 1) {
                    this.readMessagesCount = this.messages.length;
                }
            }
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
                <v-tabs v-model="tab" grow background-color="primary">
                    <v-tab>Publish</v-tab>
                    <v-tab>
                        <v-badge :color="unreadMessagesCount === 0 ? 'grey' : 'green'" :content="unreadMessagesCount.toString()">Subscribe</v-badge>
                    </v-tab>
                    <v-tab>
                        <v-badge :color="schemasCount === 0 ? 'grey' : 'green'" :content="schemasCount.toString()">Manage schemas</v-badge>
                    </v-tab>
                </v-tabs>
                <v-tabs-items v-model="tab">
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
                                <v-alert dense outlined dismissible type="info">Publish message using JSON examples with default values generated from Avro schemas. Subject is filled using "namespace" field.</v-alert>
                                <v-alert v-if="!schemas.length" dense outlined dismissible type="warning">You need upload zip file with Avro schemas.</v-alert>
                                <Publish :schemas="schemas" @success="onSuccess" @error="onError"></Publish>
                            </v-card-text>
                        </v-card>
                    </v-tab-item>
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
                                <v-alert dense outlined dismissible type="info">You'll see messages from all subjects. Even they are not encoded using Avro.</v-alert>
                                <Subscribe :events="messages"></Subscribe>
                            </v-card-text>
                        </v-card>
                    </v-tab-item>
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
                                <UploadSchemas @success="onUploadedSchemasSuccess" @error="onError"></UploadSchemas>
                            </v-card-text>
                        </v-card>
                    </v-tab-item>
                </v-tabs-items>
                <v-snackbar v-model="snackbar.model" :color="snackbar.color" :timeout="snackbar.timeout">
                    {{ snackbar.text }}
                    <template v-slot:action="{ attrs }">
                        <v-btn dark text v-bind="attrs" @click="snackbar.model = false">Close</v-btn>
                    </template>
                </v-snackbar>
            </v-container>
        </v-main>
    </v-app>
</template>

<style>
    [v-cloak] {
        display: none;
    }
</style>
