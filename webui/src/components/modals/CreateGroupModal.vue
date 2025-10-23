<script setup>
import { computed, ref } from "vue";
import { usersStore } from "../../stores/usersStore";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue";
import LabeledInput from "../ui/LabeledInput.vue";
import { auth } from "../../stores/authStore";
import { createGroup } from "../../services/groups";
import { useToast } from "vue-toast-notification";
import { getError } from "../../utils/getError";
import { useModal } from "../../hooks/useModal";

const { closeModal } = useModal();
const $toast = useToast();
const groupName = ref("");
const selectedUsersIds = ref([]);

const selectableUsers = computed(() => {
	return Object.values(usersStore.users).filter(
		(user) => user.userId != auth.userId // filter out the user himself
	);
});

function handleSelectUser(userId) {
	if (selectedUsersIds.value.includes(userId)) {
		selectedUsersIds.value = selectedUsersIds.value.filter(
			(uid) => uid !== userId
		);
	} else {
		selectedUsersIds.value = [...selectedUsersIds.value, userId];
	}
}

async function handleConfirm() {
	if (!selectedUsersIds.value?.length) return;

	const res = await createGroup(groupName.value, selectedUsersIds.value);
	console.log(res);
	if (res.status === 201) {
		$toast.success("Il gruppo è stato creato");
		console.log(res.data);
		// conversationsStore.addConversation(res.data);
		closeModal();
	} else {
		const err = getError(res, "Impossibile creare il gruppo");
		$toast.error(err);
	}
}

const submitEnabled = computed(
	() => groupName.value.trim() && selectedUsersIds.value.length > 0
);
</script>
<template>
	<Modal>
		<div class="createGroupModal-container">
			<ModalTitle text="Crea un gruppo!" />
			<p>Seleziona un utente dalla lista ed inizia subito a chattare!</p>

			<LabeledInput
				label="Nome"
				placeholder="Inserisci il nome del gruppo..."
				v-model:value="groupName"
			/>

			<div class="createGroupModal-separator"></div>
			<UsersSelector
				:users="selectableUsers"
				:selectedUserIds="selectedUsersIds"
				:onSelect="handleSelectUser"
			/>
			<button
				v-if="selectableUsers.length"
				@click="handleConfirm"
				:disabled="!submitEnabled"
			>
				Crea gruppo
			</button>
		</div>
	</Modal>
</template>

<style scoped>
.createGroupModal-container {
	display: flex;
	flex-direction: column;
}
.createGroupModal-container p {
	color: var(--color-white);
	font-size: 18px;
}
.createGroupModal-separator {
	margin-top: 16px;
}

.createGroupModal-container button {
	background-color: var(--color-green-primary);
	margin-left: auto;
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	margin-top: 16px;
	margin-left: auto;
}
.createGroupModal-container button:hover {
	opacity: 0.7;
}
.createGroupModal-container button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
