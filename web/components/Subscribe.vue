<script>
    export default {
        props: ['events'],
        data() {
            return {
                list: [],
            };
        },
        watch: {
            events(val) {
                this.list = val;
            },
        },
    }
</script>

<template>
    <div>
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
                <tr v-for="item in list" :key="item.name">
                    <td class="text-left">{{ item.id }}</td>
                    <td class="text-left">{{ item.subject }}</td>
                    <td class="text-left">{{ item.payload }}</td>
                </tr>
                </tbody>
            </template>
        </v-simple-table>
    </div>
</template>
