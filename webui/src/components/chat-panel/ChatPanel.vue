<script setup>
import { watch, computed, onMounted, onUnmounted } from "vue";
import { conversationStore } from "@/stores/conversationStore";
import UnselectedConversation from "./unselected-conversation/UnselectedConversation.vue";
import Conversation from "./conversation/Conversation.vue";

const currentConversation = computed(() => conversationStore.currentConversation.value);

let pollInterval = null;

onMounted(() => {
	startPolling();
});

onUnmounted(() => {
	stopPolling();
});

function startPolling() {
	stopPolling(); // clear any previous intervals

	pollInterval = setInterval(async () => {
		if (currentConversation.value) {
			const res = await conversationStore.pollConversation();
			if (res?.shouldClose) conversationStore.closeConversation(); // make sure the user gets the conversation closed if he is notified
		} else {
		}
	}, 1000);
}

function stopPolling() {
	if (pollInterval) {
		clearInterval(pollInterval);
		pollInterval = null;
	}
}
</script>

<template>
	<div class="chatPanel-container">
		<UnselectedConversation v-if="!currentConversation" />
		<Conversation v-else :conversation="currentConversation" />
	</div>
</template>
<style scoped>
.chatPanel-container {
	width: 100%;
	min-width: 500px;
	height: 90vh;
	background-color: var(--color-black);
	border: 1px solid var(--color-white-20);

	border-radius: 8px;
}
</style>
