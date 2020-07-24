<script>
    export default {
        props: ['data'],
    }
</script>

<template>
    <span v-if="data">{{ data.current }}<span v-if="data.should_update"> (last version {{ data.last }})</span></span>
</template>
