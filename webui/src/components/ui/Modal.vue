<script setup></script>

<template>
	<div class="modal-container">
		<slot />
	</div>
</template>

<style scoped>
.modal-container {
	width: 600px;
	height: 600px;
	background-color: var(--color-primary);
	padding: 8px;
}
</style>
