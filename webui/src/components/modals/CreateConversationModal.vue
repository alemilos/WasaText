<script setup>
import { ref } from "vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { useModal } from "../../hooks/useModal";
import { createConversation } from "../../services/conversations";
import { useToast } from "vue-toast-notification";
import { computed } from "vue";
import { usersStore } from "../../stores/usersStore";

const users = computed(() => usersStore.users);

const { closeModal } = useModal();
const $toast = useToast();

const selectedUserId = ref(null);

function handleSelectUser(userId) {
	if (selectedUserId.value == userId) {
		selectedUserId.value = null;
	} else {
		selectedUserId.value = userId;
	}
}

async function handleConfirm() {
	if (!selectedUserId.value) return;
	const res = await createConversation(selectedUserId.value);
	if (res.status === 201) {
		$toast.success("La conversazione è stata creata");
		closeModal();
	} else {
		console.log(res);
	}
}
</script>
<template>
	<Modal>
		<div class="createConversationModal-container">
			<ModalTitle text="Inizia a chattare con un altro utente!" />
			<p>Seleziona un utente dalla lista ed inizia subito a chattare!</p>
			<UsersSelector
				:users="users"
				:selected-user-id="selectedUserId"
				:on-select="handleSelectUser"
			/>
			<button @click="handleConfirm" :disabled="!selectedUserId">
				Conferma
			</button>
		</div>
	</Modal>
</template>

<style scoped>
.createConversationModal-container {
	display: flex;
	flex-direction: column;
}
.createConversationModal-container p {
	color: var(--color-white);
	font-size: 18px;
}

.createConversationModal-container button {
	background-color: var(--color-green-primary);
	margin-left: auto;
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	margin-top: 16px;
	margin-left: auto;
}
.createConversationModal-container button:hover {
	opacity: 0.7;
}
.createConversationModal-container button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
