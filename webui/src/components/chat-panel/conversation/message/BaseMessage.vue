<script setup>
import { computed } from "vue";
import { conversationStore } from "../../../../stores/conversationStore";
import { auth } from "../../../../stores/authStore";
import { usersStore } from "../../../../stores/usersStore";

const props = defineProps({
	message: { type: Object, required: true },
	isSent: { type: Boolean, required: true },
	conversationType: { type: String, default: "private" }, // "private" or "group"
});

const alignment = computed(() => (props.isSent ? "sent" : "received"));
const bubbleColor = computed(() =>
	props.isSent ? "sent-bubble" : "received-bubble"
);

const formatTime = (dateString) => {
	return new Date(dateString).toLocaleTimeString("it-IT", {
		hour: "2-digit",
		minute: "2-digit",
	});
};

const senderName = computed(() => {
	const isGroup = props.conversationType === "group";
	if (isGroup && !props.isSent) {
		return usersStore.getUsername(props.message.authorId);
	}
	return null;
});

const getStatusText = computed(() => {
	const status = props.message.messageStatus?.status;
	const members = props.message.messageStatus?.members || [];

	if (props.conversationType === "group") {
		const totalMembers =
			conversationStore.currentConversation.value?.participants?.length ||
			0;
		return `Letto da ${members.length}/${totalMembers}`;
	} else {
		// private conversation
		return status === "received" ? "Letto" : "Inviato";
	}
});
</script>

<template>
	<div :class="['base-message', alignment]">
		<div :class="['message-bubble', bubbleColor]">
			<!-- Sender name for group messages (only received) -->
			<div v-if="senderName" class="sender-name">
				~ {{ senderName }}
				<!-- You might want to map this to username -->
			</div>

			<!-- Message content slot -->
			<slot></slot>

			<!-- Status bar -->
			<div class="status-bar">
				<span class="time">{{ formatTime(message.createdAt) }}</span>
				<span class="status" v-if="isSent">{{ getStatusText }}</span>
			</div>
		</div>
	</div>
</template>

<style scoped>
.base-message {
	display: flex;
	margin: 8px 0;
}

.base-message.sent {
	justify-content: flex-end;
}

.base-message.received {
	justify-content: flex-start;
}

.base-message.sent .message-bubble {
	border-top-right-radius: 0px;
}

.base-message.received .message-bubble {
	border-top-left-radius: 0px;
}

.message-bubble {
	max-width: 70%;
	padding: 12px;
	border-radius: 12px;
	position: relative;

	color: var(--color-green-primary);
	border: 1px solid var(--color-white-20);
}

.sent-bubble {
	background-color: #07581c;
}

.received-bubble {
	background-color: #0b0c0e;
}

.sender-name {
	font-size: 14px;
	font-weight: bold;
	margin-bottom: 4px;
	color: var(--color-white);
}

.status-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: 4px;
	font-size: 11px;
	max-width: fit-content;
	color: var(--color-green-primary);
	justify-self: end;
}

.time {
	margin-right: 8px;
}
</style>
