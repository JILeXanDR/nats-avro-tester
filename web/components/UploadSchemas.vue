<script>
    export default {
        data() {
            return {
                file: null,
            };
        },
        computed: {
            fileEmpty() {
                return this.file === null;
            },
        },
        methods: {
            async uploadSchemas() {
                try {
                    const res = await this.$backend.uploadSchema(this.file);
                    this.clearFile();
                    this.$emit('success', 'Schemas was synchronized.');
                } catch (e) {
                    this.$emit('error', e.message);
                }
            },
            clearFile() {
                this.file = null;
            }
        },
    }
</script>

<template>
    <div>
        <v-form @submit.prevent="uploadSchemas">
            <v-file-input v-model="file" label="Schemas zip file" outlined dense></v-file-input>
            <v-btn type=submit :disabled="fileEmpty" clearable color="success" class="mr-4">Upload</v-btn>
        </v-form>
    </div>
</template>
