<script>
    import Publish from './components/Publish';
    import Subscribe from './components/Subscribe';
    import UploadSchemas from './components/UploadSchemas';
    import Version from './components/Version';

    export default {
        components: {
            Publish,
            Subscribe,
            UploadSchemas,
            Version,
        },
        data() {
            return {
                version: null,
                tab: null,
                schemas: [],
                messages: [],
                readMessagesCount: 0,
            };
        },
        computed: {
            lastMessages() {
                return this.messages;
            },
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
            this.$backend.checkVersion().then(res => {
                this.version = res;
            });
            this.$backend.connectMessagesStream(message => this.messages.push(message));
        },
        methods: {
            async loadSchemas() {
                try {
                    this.$backend.fetchSchemas().then(res => this.schemas = res);
                } catch (e) {
                    this.$notify('error', e.message);
                }
            },
            onSuccess(message) {
                this.$notify('success', message);
            },
            onError(message) {
                this.$notify('error', message);
            },
            onUploadedSchemasSuccess(message) {
                this.$notify('success', message);
                this.loadSchemas();
            },
            clearMessages() {
                this.messages = [];
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
            <v-toolbar-title>NATS AVRO TESTER <Version :data="version"></Version></v-toolbar-title>
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
                                <Subscribe :events="lastMessages">
                                    <v-btn @click="clearMessages" color="info" class="mr-4">Clear</v-btn>
                                </Subscribe>
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
                <v-snackbar v-model="notifierPlugin.model" :color="notifierPlugin.color" :timeout="notifierPlugin.timeout">
                    {{ notifierPlugin.text }}
                    <template v-slot:action="{ attrs }">
                        <v-btn dark text v-bind="attrs" @click="notifierPlugin.model = false">Close</v-btn>
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
