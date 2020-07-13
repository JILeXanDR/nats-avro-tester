<script>
    export default {
        data() {
            return {
                schemasFile: null,
            };
        },
        methods: {
            async uploadSchemas() {
                try {
                    const res = await this.$backend.uploadSchema(this.schemasFile);
                    this.schemasFile = null;
                    this.$emit('success', res.message);
                } catch (e) {
                    this.$emit('error', e.message);
                }
            },
        },
    }
</script>

<template>
    <div>
        <v-form @submit.prevent="uploadSchemas">
                <v-file-input v-model="schemasFile" label="Schemas zip file" outlined dense></v-file-input>
                <v-btn type=submit :disabled="schemasFile==null" clearable color="success" class="mr-4">
                    Upload
                </v-btn>
        </v-form>
    </div>
</template>

<style scoped>

</style>
