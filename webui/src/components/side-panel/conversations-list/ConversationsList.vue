<script setup>
import { getMyConversations } from "../../../services/conversations";
import { ref, onMounted } from "vue";
import GroupPreview from "./GroupPreview.vue";
import PrivateConversationPreview from "./PrivateConversationPreview.vue";
import { conversationStore } from "../../../stores/conversationStore";

const conversations = ref([]);

onMounted(async () => {
	try {
		const res = await getMyConversations();
		if (res.status === 200) {
			conversations.value = res.data.conversations;
		}
	} catch (err) {}
});

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
