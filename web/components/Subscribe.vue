<script>
    export default {
        props: ['events'],
        data() {
            return {};
        },
    }
</script>

<template>
    <div>
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
                    <td class="text-left">{{ item.payload }}</td>
                    <td class="text-left">{{ item.when }}</td>
                </tr>
                </tbody>
            </template>
        </v-simple-table>
    </div>
</template>
