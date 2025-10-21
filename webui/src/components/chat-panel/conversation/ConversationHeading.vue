<script setup>
import { computed } from "vue";
import { conversationStore } from "../../../stores/conversationStore";
import { usersStore } from "../../../stores/usersStore";
import UserPhoto from "../../reusables/UserPhoto.vue";
import defaultGroupImage from "@/assets/images/default-group.png";

const currentConversation = computed(
	() => conversationStore.currentConversation.value
);

const conversationType = computed(
	() => currentConversation.value?.type ?? "private"
);

const conversationName = computed(() => {
	const conv = currentConversation.value;
	if (!conv) return null;
	if (conv.type === "group") return conv.name;
	if (conv.type === "private") return usersStore.getUsername(otherId.value);
	return null;
});

const otherId = computed(() => currentConversation.value?.otherParticipantId);

const photoPath = computed(() => {
	const conv = currentConversation.value;
	if (!conv) return null;
	if (conv.type === "group") return conv.photoPath ?? defaultGroupImage;
	if (conv.type === "private") return usersStore.getPhotoUrl(otherId.value);
	return null;
});
</script>

<template>
	<div class="conversationHeading-container">
		<UserPhoto :url="photoPath" :type="conversationType" :size="90" />
		<p>{{ conversationName }}</p>
	</div>
</template>
