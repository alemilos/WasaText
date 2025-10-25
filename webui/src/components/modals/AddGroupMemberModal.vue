<script setup>
import { ref, watch } from "vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { useModal } from "../../hooks/useModal";
import { useToast } from "vue-toast-notification";
import { computed } from "vue";
import { usersStore } from "../../stores/usersStore";
import { auth } from "../../stores/authStore";
import GroupInformationsModal from "./GroupInformationsModal.vue";
import { addToGroup } from "../../services/groups";
import { getError } from "../../utils/getError";
import { conversationStore } from "../../stores/conversationStore";

const props = defineProps({
	conversation: Object,
});

const userId = computed(() => auth.userId);

const conversationId = computed(() => props.conversation?.conversationId);

// members need to be reactive
const members = ref([...(props.conversation?.members || [])]);
watch(
	() => conversationStore.currentConversation.value?.members,
	(newMembers) => {
		if (newMembers) {
			members.value = [...newMembers];
		}
	},
	{ immediate: true }
);
const conversationMembers = computed(() => members.value);

const isAdmin = computed(() => {
	return conversationMembers.value.find((user) => user.userId == userId.value).role === "admin";
});

const selectedUserIds = ref([]);

// users that  haven't started a conversation with the authenticated user
const notMembers = computed(() => {
	return Object.values(usersStore.users).filter(
		(user) => !conversationMembers.value.find((u) => u.userId === user.userId)
	);
});

const { openModal, closeModal } = useModal();
const $toast = useToast();

function handleSelectUser(userId) {
	if (!isAdmin.value) return;

	if (selectedUserIds.value.includes(userId)) {
		selectedUserIds.value = selectedUserIds.value.filter((uid) => uid !== userId);
	} else {
		selectedUserIds.value = [...selectedUserIds.value, userId];
	}
}

async function handleConfirm() {
	if (!selectedUserIds.value.length) return;

	const results = await Promise.all(selectedUserIds.value.map((id) => addToGroup(conversationId.value, Number(id))));

	let errorThrown = false;
	const successfullyAddedIds = [];

	results.forEach((res, index) => {
		const id = selectedUserIds.value[index];
		if (res.status === 200) {
			successfullyAddedIds.push(id);
		} else {
			const err = getError(res);
			if (!errorThrown) {
				errorThrown = true;
				$toast.error(err);
			}
		}
	});

	if (successfullyAddedIds.length) {
		$toast.success(`Hai aggiunto ${successfullyAddedIds.length} utenti`);
	}

	selectedUserIds.value = [];
}

function handleGoBack() {
	openModal(GroupInformationsModal, { conversation: props.conversation });
}
</script>
<template>
	<Modal>
		<div class="addGroupMemberModal-container">
			<ModalTitle text="Aggiungi un membro al gruppo" />
			<p v-if="notMembers.length">Seleziona gli utenti che vuoi aggiungere al gruppo!</p>
			<p v-if="!notMembers.length" style="color: #ff4444">
				Non ci sono nuovi utenti da poter aggiungere al gruppo!
			</p>
			<UsersSelector
				v-if="notMembers.length"
				:users="notMembers"
				:selected-user-ids="selectedUserIds"
				:on-select="handleSelectUser"
			/>
			<div class="addGroupMemberModal-buttons">
				<button style="background-color: white" @click="handleGoBack">Torna indietro</button>
				<button @click="handleConfirm" :disabled="!selectedUserIds.length">Aggiungi</button>
			</div>
		</div>
	</Modal>
</template>

<style scoped>
.addGroupMemberModal-container {
	display: flex;
	flex-direction: column;
	height: 100%;
}
.addGroupMemberModal-container p {
	color: var(--color-white);
	font-size: 18px;
}

.addGroupMemberModal-container button {
	background-color: var(--color-green-primary);
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	margin-top: 16px;
}
.addGroupMemberModal-container button:hover {
	opacity: 0.7;
}
.addGroupMemberModal-container button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.addGroupMemberModal-buttons {
	margin-top: auto;
	display: flex;
	align-items: center;
	gap: 8px;
	width: 100%;
	justify-content: space-between;
}
</style>
