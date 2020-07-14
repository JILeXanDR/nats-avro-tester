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
                alerts: [],
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
            lastMessages() {
                return this.messages;
            },
            messagesCount() {
                return this.lastMessages.length;
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
            showError(text) {
                this.snackbar.model = true;
                this.snackbar.color = 'error';
                this.snackbar.text = text;
            },
            showNotification(text) {
                this.snackbar.model = true;
                this.snackbar.color = 'success';
                this.snackbar.text = text;
            },
            onSuccess(message) {
                this.showNotification(message);
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
                                <v-alert dense outlined dismissible type="info">Publish message using JSON examples generated from Avro schemas. Subject is received from the "namespace" field.</v-alert>
                                <v-alert v-if="!schemas.length" dense outlined dismissible type="warning">You need upload zip file with Avro schemas.</v-alert>
                                <Publish :schemas="schemas" @success="onSuccess" @error="onError"></Publish>
                            </v-card-text>
                        </v-card>
                    </v-tab-item>
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
                                <v-alert dense outlined dismissible type="info">You'll see messages from all subjects</v-alert>
                                <Subscribe :events="lastMessages"></Subscribe>
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
