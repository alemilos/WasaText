<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from "vue";
import { usersStore } from "../../../../stores/usersStore";
import ChevronWhiteIcon from "@/assets/icons/chevron-white.svg";
import UserPhoto from "../../../reusables/UserPhoto.vue";
import DoubleTicks from "@/assets/icons/doubleticks.svg";
import SingleTick from "@/assets/icons/singletick.svg";
import PendingMessage from "@/assets/icons/sendingmessage.svg";
import ActionsPopup from "./ActionsPopup.vue";

const props = defineProps({
	message: { type: Object, required: true },
	isSent: { type: Boolean, required: true },
	conversationType: { type: String, default: "private" }, // "private" or "group"
});

const alignment = computed(() => (props.isSent ? "sent" : "received"));
const bubbleColor = computed(() => (props.isSent ? "sent-bubble" : "received-bubble"));
const messageReactions = computed(() => props.message.comments);

const groupedReactions = computed(() => {
	const counts = {};
	messageReactions.value?.forEach((reaction) => {
		if (counts[reaction.emoji]) counts[reaction.emoji]++;
		else counts[reaction.emoji] = 1;
	});

	return counts;
});

const showPopup = ref(false);
const popupRef = ref(null);

function handleClickOutside(event) {
	if (showPopup.value && popupRef.value && !popupRef.value.contains(event.target)) {
		showPopup.value = false;
	}
}

onMounted(() => {
	document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
	document.removeEventListener("click", handleClickOutside);
});

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

const senderPhotoPath = computed(() => {
	const isGroup = props.conversationType === "group";
	if (isGroup && !props.isSent) {
		return usersStore.getPhotoUrl(props.message.authorId);
	}
	return null;
});

const statusIcon = computed(() => {
	const status = props.message.messageStatus?.status;

	if (status === "read") return DoubleTicks;
	if (status === "received") return SingleTick;
	return PendingMessage; // default: pending
});

const togglePopup = () => {
	showPopup.value = !showPopup.value;
};
</script>

<template>
	<div :class="['base-message', alignment]">
		<div :class="['message-bubble', bubbleColor]">
			<img :src="ChevronWhiteIcon" class="chevron-icon" @click.stop="togglePopup" />

			<!-- Actions Popup -->
			<div v-if="showPopup" ref="popupRef">
				<ActionsPopup :isSent="isSent" :message="message" @close="showPopup = false" />
			</div>

			<!-- Sender name & photofor group messages (only received) -->
			<div v-if="senderName" class="sender-name">
				<UserPhoto :url="senderPhotoPath" :size="24" />
				<span> ~ {{ senderName }} </span>
			</div>

			<slot></slot>

			<div class="status-bar">
				<span class="time">{{ formatTime(message.createdAt) }}</span>
				<img v-if="isSent" :src="statusIcon" class="status-icon" />
			</div>

			<div v-if="Object.keys(groupedReactions).length" class="reactions-container">
				<div class="reactions-scroll">
					<div v-for="(count, emoji) in groupedReactions" :key="emoji" class="reaction-item">
						<span class="emoji">{{ emoji }}</span>
						<span class="count">{{ count }}</span>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div v-if="Object.keys(groupedReactions).length" class="reactionsSeparator"></div>
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

.chevron-icon {
	position: absolute;
	top: 3px;
	right: 6px;
	width: 16px;
	height: 16px;
	cursor: pointer;
	opacity: 0.7;
	transition: opacity 0.2s;
}

.chevron-icon:hover {
	opacity: 1;
}

.message-bubble {
	max-width: 70%;
	padding: 12px;
	padding-top: 16px;
	border-radius: 12px;
	position: relative;
	min-width: 70px;

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

	display: flex;
	gap: 8px;
	align-items: center;
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

.status-icon {
	width: 12px;
	height: 12px;
	margin-left: 4px;
	vertical-align: middle;
}

.time {
	margin-right: 8px;
}

.reactionsSeparator {
	margin-bottom: 40px;
}

.reactions-container {
	position: absolute;
	bottom: -28px;
	left: 8px;
	background-color: var(--color-green-dark);
	border: 1px solid var(--color-white-20);
	border-radius: 12px;
	padding: 4px 8px;
	max-width: 90%;
	overflow-x: auto;
	white-space: nowrap;
	display: flex;
	align-items: center;
	gap: 4px;
	scrollbar-width: thin;
}

.reactions-container::-webkit-scrollbar {
	height: 4px;
}
.reactions-container::-webkit-scrollbar-thumb {
	background: rgba(255, 255, 255, 0.2);
	border-radius: 4px;
}

.reactions-scroll {
	display: flex;
	align-items: center;
	gap: 6px;
}

.reaction-item {
	display: inline-flex;
	align-items: center;
	gap: 4px;
	background: rgba(255, 255, 255, 0.1);
	padding: 2px 6px;
	border-radius: 8px;
	font-size: 14px;
	color: white;
	flex-shrink: 0;
}

.reaction-item .emoji {
	font-size: 16px;
	line-height: 1;
}
.reaction-item .count {
	font-size: 12px;
	opacity: 0.8;
}
</style>
