<script setup>
import { computed } from "vue";
import { conversationStore } from "../../../stores/conversationStore";
import { usersStore } from "../../../stores/usersStore";
import UserPhoto from "../../reusables/UserPhoto.vue";
import defaultGroupImage from "@/assets/images/default-group.png";
import { useModal } from "../../../hooks/useModal";
import GroupInformationsModal from "../../modals/GroupInformationsModal.vue";
import PrivateChatInformationsModal from "../../modals/PrivateChatInformationsModal.vue";
const { openModal } = useModal();

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

function handleOpenConversationInfo() {
	const conv = currentConversation.value;
	if (!conv) return;

	if (conv.type === "group") {
		openModal(GroupInformationsModal, { conversation: conv });
	} else if (conv.type === "private") {
		openModal(PrivateChatInformationsModal, { conversation: conv });
	}
}
</script>

<template>
	<div
		class="conversationHeading-container"
		@click="handleOpenConversationInfo"
	>
		<UserPhoto :url="photoPath" :type="conversationType" :size="90" />
		<span>{{ conversationName }}</span>
	</div>
</template>
