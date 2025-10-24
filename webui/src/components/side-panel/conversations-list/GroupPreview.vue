<script setup>
import UserPhoto from "../../reusables/UserPhoto.vue";
import LastMessagePreview from "./LastMessagePreview.vue";
import "./style/Conversation.css";

const props = defineProps({
	conversation: { type: Object },
	onOpenConversation: { type: Function, required: true },
});

function handleClick() {
	props.onOpenConversation(props.conversation.conversationId);
}
</script>

<template>
	<div class="conversationPreview-container" @click="handleClick">
		<UserPhoto :url="props.conversation?.photoPath" type="group" />
		<LastMessagePreview
			:conversationName="conversation?.name"
			:type="conversation?.lastMessage?.type"
			:content="conversation?.lastMessage?.content"
			:createdAt="conversation?.lastMessage?.createdAt"
			:createdBy="conversation?.lastMessage?.createdBy"
		/>
	</div>
</template>
