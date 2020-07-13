<script>
    import Publish from './components/Publish';
    import Subscribe from './components/Subscribe';
    import UploadSchemas from './components/UploadSchemas';

    export default {
        components: {Publish, Subscribe, UploadSchemas},
        data() {
            return {
                tab: null,
                items: [
                    'Appetizers', 'Entrees', 'Deserts', 'Cocktails',
                ],
                text: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
                schemas: [],
                messages: [],
                alerts: [],
            };
        },
        computed: {
            lastMessages() {
                return this.messages.reverse();
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
            messagesCount() {
                return this.lastMessages.length.toString();
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
                this.alerts.push({type: 'error', text});
            },
            showNotification(text) {
                this.alerts.push({type: 'success', text});
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
    }
</script>

<template>
    <v-app>
        <v-app-bar app color="indigo" dark>
            <v-toolbar-title>NATS AVRO TESTER</v-toolbar-title>
        </v-app-bar>
        <v-main>
            <v-container fluid>

                <v-tabs v-model="tab" color="basil" grow dark background-color="primary">
                    <v-tab>Publish</v-tab>
                    <v-tab>
                        <v-badge color="green" :content="messagesCount">Subscribe</v-badge>
                    </v-tab>
                    <v-tab>Manage schemas</v-tab>
                </v-tabs>

                <v-tabs-items v-model="tab">
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
                                <Publish :schemas="schemas" @success="onSuccess" @error="onError"></Publish>
                            </v-card-text>
                        </v-card>
                    </v-tab-item>
                    <v-tab-item>
                        <v-card flat>
                            <v-card-text>
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

                <v-alert v-for="item in lastAlerts" dense outlined dismissible :type="item.type">{{ item.text }}
                </v-alert>

                <v-row>
                    <v-col>
                    </v-col>
                </v-row>
                <v-row>
                    <v-col>
                    </v-col>
                </v-row>
                <v-row>
                    <v-col>
                    </v-col>
                </v-row>
            </v-container>
        </v-main>
    </v-app>
</template>

<style>
    [v-cloak] {
        display: none;
    }

    /* Helper classes */
    .basil {
        background-color: #FFFBE6 !important;
    }

    .basil--text {
        color: #356859 !important;
    }
</style>
