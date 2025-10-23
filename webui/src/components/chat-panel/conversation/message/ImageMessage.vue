<script setup>
import BaseMessage from "./BaseMessage.vue";

const props = defineProps({
	message: Object,
	isSent: Boolean,
	conversationType: String,
});

const API_URL = __API_URL__;
</script>

<template>
	<BaseMessage
		:message="message"
		:is-sent="isSent"
		:conversation-type="conversationType"
	>
		<div class="image-content">
			<img
				:src="`${API_URL}/uploads${message.content}`"
				:alt="'Foto'"
				class="message-image"
			/>
			<div v-if="message.secondaryContent" class="image-caption">
				{{ message.secondaryContent }}
			</div>
		</div>
	</BaseMessage>
</template>

<style scoped>
.image-content {
	display: flex;
	flex-direction: column;
	gap: 8px;
	align-items: flex-start;
}

.message-image {
	max-width: 100%;
	max-height: 40vh;
	height: auto;
	width: auto;
	border-radius: 8px;
	object-fit: contain;
}

.image-caption {
	color: var(--color-green-primary);
	padding: 4px 0;
	text-align: left;
}
</style>
