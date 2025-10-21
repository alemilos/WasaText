<script setup>
import "./Conversation.css";
import { computed, ref, watch, onMounted, nextTick } from "vue";
import { conversationStore } from "../../../stores/conversationStore";
import { auth } from "../../../stores/authStore";
import TextMessage from "./message/TextMessage.vue";
import ImageMessage from "./message/ImageMessage.vue";

const messages = computed(
	() => conversationStore.currentConversation.value?.messages || []
);
const conversationType = computed(
	() => conversationStore.currentConversation.value?.type
);
const userId = computed(() => auth.userId);

const messagesContainer = ref(null);

// Scroll to bottom when messages change
watch(
	messages,
	() => {
		nextTick(() => {
			scrollToBottom();
		});
	},
	{ deep: true }
);

// Also scroll when component mounts
onMounted(() => {
	nextTick(() => {
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
			/>
		</div>
	</div>
</template>
