<script>
    export default {
        props: {
            events: {
                required: true,
                type: Array,
            },
        },
        data() {
            return {};
        },
        filters: {
            pretty(str) {
                return JSON.stringify(str, null, '  ');
            },
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
                    <th class="text-left">#</th>
                    <th class="text-left">Subject</th>
                    <th class="text-left">Payload</th>
                    <th class="text-left">When</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item, index) in events" :key="item.name">
                    <td class="text-left">
                        {{ index + 1 }}
                    </td>
                    <td class="text-left">
                        <code>{{ item.subject }}</code>
                    </td>
                    <td class="text-left">
                        <v-tooltip left>
                            <template v-slot:activator="{ on, attrs }">
                                <span class="text-caption" v-bind="attrs" v-on="on">{{ item.payload }}</span>
                            </template>
                            <pre>{{ item.payload | pretty }}</pre>
                        </v-tooltip>
                    </td>
                    <td class="text-left">
                        {{ item.when | moment("from") }}
                    </td>
                </tr>
                </tbody>
            </template>
        </v-simple-table>
    </div>
</template>
