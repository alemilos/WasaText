<script setup>
import { computed, ref } from "vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue";
import Modal from "../ui/Modal.vue";
import { usersStore } from "../../stores/usersStore";
import { conversationsStore } from "../../stores/conversationsStore";
import ModalTitle from "../ui/ModalTitle.vue";
import { auth } from "../../stores/authStore";
import PlusIcon from "@/assets/images/plusicon.png";
import { leaveGroup } from "../../services/groups";
import { useToast } from "vue-toast-notification";
import { getError } from "../../utils/getError";
import { useModal } from "../../hooks/useModal";
import { conversationStore } from "../../stores/conversationStore";

const props = defineProps({
	conversation: Object,
});

const $toast = useToast();
const { closeModal } = useModal();

const userId = computed(() => auth.userId);
const conversationId = computed(() => props.conversation?.conversationId);
const conversationMembers = computed(() => props.conversation?.members);

const selectedUserIds = ref([]);

const isAdmin = computed(() => {
	return conversationMembers.value.find((user) => user.userId == userId.value).role === "admin";
});

function onSelect(userId) {
	if (!isAdmin.value) return;

	if (selectedUserIds.value.includes(userId)) {
		selectedUserIds.value = selectedUserIds.value.filter((uid) => uid !== userId);
	} else {
		selectedUserIds.value = [...selectedUserIds.value, userId];
	}
}

async function removeUsers() {
	if (!selectedUserIds.value.length) return;

	try {
		const shouldDeleteHimself = selectedUserIds.value.some((uid) => uid == userId.value);
		const otherUsersIds = selectedUserIds.value.filter((uid) => uid != userId.value);

		// First delete other users
		const results = await Promise.all(otherUsersIds.map((id) => leaveGroup(conversationId.value, Number(id))));

		let errorThrown = false;
		const successfullyRemovedIds = [];

		results.forEach((res, index) => {
			const id = selectedUserIds.value[index];
			if (res.status === 204) {
				successfullyRemovedIds.push(id);
			} else {
				const err = getError(res);
				if (!errorThrown) {
					errorThrown = true;
					$toast.error(err);
				}
			}
		});

		if (successfullyRemovedIds.length) {
			// conversationStore.removeMembers(successfullyRemovedIds);
			console.log("removed ", successfullyRemovedIds);
			$toast.success(`Hai rimosso ${successfullyRemovedIds.length} utenti`);
		}

		// Remove himself if chosen
		if (shouldDeleteHimself) {
			const res = await leaveGroup(conversationId.value, Number(userId.value));
			if (res.status === 204) {
				$toast.success("Hai abbandonato il gruppo");
				conversationStore.closeConversation();
				conversationsStore.removeConversation(conversationId.value);
				closeModal();
			}
		}

		selectedUserIds.value = [];
	} catch (err) {
		$toast.error("Errore durante la rimozione degli utenti");
	}
}

async function leaveGroupHandler() {
	const res = await leaveGroup(conversationId.value, Number(userId.value));
	if (res.status === 204) {
		$toast.success("Hai abbandonato il gruppo");
		conversationStore.closeConversation();
		conversationsStore.removeConversation(conversationId.value);
		closeModal();
	} else {
		const err = getError(res);
		$toast.error(err);
	}
}
</script>

<template>
	<Modal>
		<div class="groupInformationsModal-container">
			<ModalTitle text="Informazioni del gruppo" />
			<div class="groupInformationsModal-currentUsers">
				<p>Utenti nel gruppo</p>

				<div v-if="isAdmin" class="addMembers-container">
					<img :src="PlusIcon" />
					<span>Aggiungi un membro</span>
				</div>

				<div class="usersSelector-wrapper">
					<UsersSelector
						:users="conversationMembers"
						:selectedUserIds="selectedUserIds"
						:onSelect="onSelect"
						:maxHeight="300"
					/>
				</div>

				<button
					v-if="isAdmin"
					class="removeUsers-button"
					:disabled="!selectedUserIds.length"
					@click="removeUsers"
				>
					Rimuovi utenti
				</button>
			</div>
			<button class="leaveGroup-button" @click="leaveGroupHandler">Esci dal gruppo</button>
		</div>
	</Modal>
</template>

<style scoped>
.groupInformationsModal-container {
	height: 100%;
	display: flex;
	flex-direction: column;
}

.groupInformationsModal-currentUsers {
	font-size: 16px;
	color: var(--color-white);
}

.addMembers-container {
	display: flex;
	align-items: center;
	gap: 12px;
	background-color: var(--color-green-primary);
	color: var(--color-black);
	border-radius: 8px;
	margin-bottom: 12px;
	padding: 8px 16px;
	cursor: pointer;
}

.addMembers-container img {
	width: 20px;
	height: 20px;
}

.usersSelector-wrapper {
	margin-bottom: 20px;
}

.removeUsers-button,
.leaveGroup-button {
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
}

.removeUsers-button {
	background-color: var(--color-green-primary);
}
.leaveGroup-button {
	margin-top: auto;
	margin-left: auto;
	background-color: var(--color-red);
}

.addMembers-container:hover,
.removeUsers-button:hover,
.leaveGroup-button:hover {
	opacity: 0.7;
}
.removeUsers-button:disabled,
.leaveGroup-button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
