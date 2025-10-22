<script setup>
const props = defineProps({
	text: { type: String, required: true },
});
</script>

<template>
	<h1 class="modalTitle">{{ text.toUpperCase() }}</h1>
</template>
<style scoped>
.modalTitle {
	color: var(--color-green-primary);
	font-size: 24px;
}
</style>
