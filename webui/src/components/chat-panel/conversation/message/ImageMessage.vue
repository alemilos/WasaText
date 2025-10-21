<script setup>
import BaseMessage from "./BaseMessage.vue";

const props = defineProps({
	message: Object,
	isSent: Boolean,
	conversationType: String,
});

console.log("IMAGE:", props.message);
</script>

<template>
	<BaseMessage
		:message="message"
		:is-sent="isSent"
		:conversation-type="conversationType"
	>
		<div class="image-content">
			<img
				:src="`${__API_URL__}/uploads/${message.content}`"
				:alt="'Foto'"
				class="message-image"
			/>
			<div v-if="message.caption" class="image-caption">
				{{ message.caption }}
			</div>
		</div>
	</BaseMessage>
</template>

<style scoped>
.image-content {
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.message-image {
	max-width: 100%;
	max-height: 300px;
	border-radius: 8px;
	object-fit: cover;
}

.image-caption {
	font-size: 14px;
	color: #333;
	padding: 4px 0;
}
</style>
