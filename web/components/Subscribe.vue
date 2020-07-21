<script>
    export default {
        props: ['events'],
        data() {
            return {
                preview: null,
            };
        },
        methods: {
            showPreview(item) {
                this.preview = item;
            },
            mouseout(item) {
                this.preview = null;
            }
        },
    }
</script>

<template>
    <div>
        <slot></slot>
        <v-simple-table dense>
            <template v-slot:default>
                <thead>
                <tr>
                    <th class="text-left">Subject</th>
                    <th class="text-left">Payload</th>
                    <th class="text-left">When</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="item in events" :key="item.name">
                    <td class="text-left">{{ item.subject }}</td>
                    <td class="text-left">
                        <pre v-if="preview === item">{{ preview.payload }}</pre>
                        <span v-else>{{ item.payload }}</span>
                    </td>
                    <td class="text-left">{{ item.when }}</td>
                    <td>
                        <v-btn v-if="preview !== item" color="success" @click="showPreview(item)">Preview</v-btn>
                        <v-btn v-else color="success" @click="preview = null">Close</v-btn>
                    </td>
                </tr>
                </tbody>
            </template>
        </v-simple-table>
    </div>
</template>
