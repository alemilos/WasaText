<script setup>
import { computed, onMounted, onUnmounted } from "vue";
import { conversationStore } from "@/stores/conversationStore";
import UnselectedConversation from "./unselected-conversation/UnselectedConversation.vue";
import Conversation from "./conversation/Conversation.vue";

// reactive reference to current conversation
const currentConversation = computed(
	() => conversationStore.currentConversation.value
);

let pollInterval = null;

onMounted(() => {
	pollInterval = setInterval(() => {
		conversationStore.pollConversation();
	}, 3000);
});

onUnmounted(() => {
	clearInterval(pollInterval);
});
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
