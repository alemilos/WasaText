<script setup>
import { computed } from "vue";
import { conversationStore } from "../../../stores/conversationStore";
import { usersStore } from "../../../stores/usersStore";
import UserPhoto from "../../reusables/UserPhoto.vue";

const currentConversation = computed(
	() => conversationStore.currentConversation.value
);

const otherId = computed(() => currentConversation.value?.otherParticipantId);
console.log({ currentConversation, otherId });

const photoPath = computed(() => {
	const conv = currentConversation.value;

	console.log("conv: ", conv);
	if (!conv) return null;

	if (conv.type === "group") {
		return conv.photoPath ?? null;
	}

	if (conv.type === "private") {
		console.log(usersStore.getPhotoUrl(otherId.value));
		return usersStore.getPhotoUrl(otherId.value);
	}

	return null;
});
</script>

<template>
	<div class="conversationHeading-container">
		<UserPhoto :url="photoPath" :size="100" />
	</div>
</template>
