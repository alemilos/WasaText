<script setup>
defineProps({
	onClick: {
		type: Function,
		required: false,
	},
});
</script>

<template>
	<div class="modalOverlay-container" @click="onClick">
		<div @click.stop>
			<slot />
		</div>
	</div>
</template>

<style scoped>
.modalOverlay-container {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	backdrop-filter: blur(8px);
	z-index: 100;
}
</style>
