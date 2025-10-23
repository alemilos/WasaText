<script setup>
import "./Conversation.css";
import { computed, ref, watch, onMounted, nextTick } from "vue";
import { conversationStore } from "../../../stores/conversationStore";
import { auth } from "../../../stores/authStore";
import TextMessage from "./message/TextMessage.vue";
import ImageMessage from "./message/ImageMessage.vue";
import isEqual from "lodash.isequal";
import { readMessages } from "../../../services/conversations";

const messages = computed(
	() => conversationStore.currentConversation.value?.messages || []
);
const conversationId = computed(
	() => conversationStore.currentConversation.value?.conversationId
);
const conversationType = computed(
	() => conversationStore.currentConversation.value?.type
);

const userId = computed(() => auth.userId);

const messagesContainer = ref(null);

async function handleReadMessages() {
	const unreadMessages = messages.value
		.filter(
			(message) =>
				!message?.messageStatus?.members?.includes(Number(userId.value))
		)
		.map((message) => message.messageId);

	await readMessages(conversationId.value, unreadMessages);
}

// Scroll to bottom when messages change
watch(
	messages,
	(newMessages, oldMessages) => {
		// avoid scrolling if identical
		if (!isEqual(newMessages, oldMessages)) {
			nextTick(() => {
				handleReadMessages();
				scrollToBottom();
			});
		}
	},
	{ deep: false } // avoid unnecessary nested reactivity
);

// Also scroll when component mounts
onMounted(() => {
	nextTick(() => {
		handleReadMessages();
		scrollToBottom();
	});
});

function scrollToBottom() {
	if (messagesContainer.value) {
		messagesContainer.value.scrollTop =
			messagesContainer.value.scrollHeight;
	}
}
</script>

<template>
	<div class="messages-container" ref="messagesContainer">
		<div v-for="message in messages" :key="message.id">
			<TextMessage
				v-if="message.type === 'text'"
				:message="message"
				:is-sent="message.authorId == userId"
				:conversationType="conversationType"
			/>
			<ImageMessage
				v-else-if="message.type === 'image'"
				:message="message"
				:is-sent="message.authorId == userId"
				:conversationType="conversationType"
			/>
		</div>
	</div>
</template>
