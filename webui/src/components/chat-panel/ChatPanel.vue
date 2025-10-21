<script setup>
import "./ChatPanel.css";

import { computed } from "vue";
import { conversationStore } from "@/stores/conversationStore";
import UnselectedConversation from "./unselected-conversation/UnselectedConversation.vue";
import Conversation from "./conversation/Conversation.vue";

// reactive reference to current conversation
const currentConversation = computed(
	() => conversationStore.currentConversation.value
);
</script>

<template>
	<div class="chatPanel-container">
		<UnselectedConversation v-if="!currentConversation" />
		<Conversation v-else :conversation="currentConversation" />
	</div>
</template>
