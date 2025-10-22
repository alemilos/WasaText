<script setup>
import ConversationsList from "../side-panel/conversations-list/ConversationsList.vue";
import UserPanel from "./user-panel/UserPanel.vue";
</script>

<template>
	<div class="sidePanel-container"><ConversationsList /> <UserPanel /></div>
</template>

<style scoped>
.sidePanel-container {
	width: 350px;
	min-width: 350px;
	height: 90vh;
	background-color: var(--color-black);
	border: 1px solid var(--color-white-20);

	border-radius: 8px;

	display: flex;
	flex-direction: column;
	justify-content: space-between;
}
</style>
