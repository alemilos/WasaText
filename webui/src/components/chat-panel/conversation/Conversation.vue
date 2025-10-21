<script setup>
import "./Conversation.css";
import ConversationHeading from "./ConversationHeading.vue";
import MessageHandler from "./MessageHandler.vue";
import Messages from "./Messages.vue";
</script>

<template>
	<div class="conversation-container">
		<ConversationHeading />
		<Messages />
		<MessageHandler />
	</div>
</template>
