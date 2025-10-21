<script setup>
import BaseMessage from "./BaseMessage.vue";

const props = defineProps({
	message: Object,
	isSent: Boolean,
	conversationType: String,
});
</script>

<template>
	<BaseMessage
		:message="message"
		:is-sent="isSent"
		:conversation-type="conversationType"
	>
		<div class="text-content">{{ message.content }}</div>
	</BaseMessage>
</template>
