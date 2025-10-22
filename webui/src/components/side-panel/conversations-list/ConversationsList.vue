<script setup>
import { computed } from "vue";
import GroupPreview from "./GroupPreview.vue";
import PrivateConversationPreview from "./PrivateConversationPreview.vue";
import { conversationStore } from "../../../stores/conversationStore";
import { conversationsStore } from "../../../stores/conversationsStore";

const conversations = computed(() =>
	Object.values(conversationsStore.conversations)
);

const conversationMap = {
	group: GroupPreview,
	private: PrivateConversationPreview,
};

function handleOpenConversation(conversationId) {
	conversationStore.openConversation(conversationId);
}
</script>

<template>
	<div class="conversationsList-container">
		<component
			v-for="conversation in conversations"
			:key="conversation.conversationId"
			:is="conversationMap[conversation.type]"
			:conversation="conversation"
			:onOpenConversation="handleOpenConversation"
		/>
	</div>
</template>

<style scoped>
.conversationsList-container {
	display: flex;
	flex-direction: column;
	gap: 8px;

	overflow-y: scroll;
	padding: 40px 25px 20px 25px;
}
</style>
