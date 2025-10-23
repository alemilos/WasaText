<script setup>
import { computed } from "vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { usersStore } from "../../stores/usersStore";
import UserPhoto from "../reusables/UserPhoto.vue";
import { formatDay } from "@/utils/time";

const props = defineProps({
	conversation: Object,
});

const otherParticipantId = computed(
	() => props.conversation?.otherParticipantId
);

const userPhotoPath = computed(() =>
	usersStore.getPhotoUrl(otherParticipantId.value)
);

const username = computed(() =>
	usersStore.getUsername(otherParticipantId.value)
);

const userCreationTime = computed(() =>
	usersStore.getUserCreationTime(otherParticipantId.value)
);
</script>
<template>
	<Modal>
		<div class="privateChatInformationModal-container">
			<ModalTitle :text="`Informazioni profilo`" />

			<div class="privateChatInformationModal-infos">
				<UserPhoto :url="userPhotoPath" :size="200" />
				<span>{{ username }}</span>
				<p>
					{{ `Su WasaText dal `
					}}<span>{{ formatDay(userCreationTime) }}</span>
				</p>
			</div>
		</div>
	</Modal>
</template>

<style scoped>
.privateChatInformationModal-container {
	display: flex;
	flex-direction: column;
	height: 100%;
}
.privateChatInformationModal-infos {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 20px;
	margin-top: 30px;
}

.privateChatInformationModal-infos span {
	font-size: 32px;
	color: var(--color-white);
}
.privateChatInformationModal-infos p {
	font-size: 22px;
	color: var(--color-white);
}
.privateChatInformationModal-infos p span {
	font-size: 22px;
	color: var(--color-green-primary);
}
</style>
