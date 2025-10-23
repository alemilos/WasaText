<script setup>
import "./style/Conversation.css";
import { computed, watch } from "vue";
import { usersStore } from "../../../stores/usersStore";
import UserPhoto from "../../reusables/UserPhoto.vue";
import LastMessagePreview from "./LastMessagePreview.vue";

const props = defineProps({
	conversation: { type: Object },
	onOpenConversation: { type: Function, required: true },
});

const otherId = computed(() => props.conversation.otherParticipantId);
const username = computed(() => usersStore.getUsername(otherId.value));
const photoUrl = computed(() => usersStore.getPhotoUrl(otherId.value));

function handleClick() {
	props.onOpenConversation(props.conversation.conversationId);
}

watch(
	() => usersStore.users.value,
	(newVal, oldVal) => {
		console.log("Users changed:", newVal);
	},
	{ deep: true } // necessary because Map is nested/reactive
);
</script>

<template>
	<div class="conversationPreview-container" @click="handleClick">
		<UserPhoto :url="photoUrl" />
		<LastMessagePreview
			:conversationName="username"
			:type="conversation.lastMessage?.type"
			:content="conversation.lastMessage?.content"
			:createdAt="conversation.lastMessage?.createdAt"
			:createdBy="conversation.lastMessage?.createdBy"
		/>
	</div>
</template>
