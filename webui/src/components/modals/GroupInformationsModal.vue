<script setup>
import { computed, ref, watch } from "vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue";
import Modal from "../ui/Modal.vue";
import { usersStore } from "../../stores/usersStore";
import { conversationsStore } from "../../stores/conversationsStore";
import ModalTitle from "../ui/ModalTitle.vue";
import { auth } from "../../stores/authStore";
import PlusIcon from "@/assets/images/plusicon.png";
import { leaveGroup, setGroupName, setGroupPhoto } from "../../services/groups";
import { useToast } from "vue-toast-notification";
import { getError } from "../../utils/getError";
import { useModal } from "../../hooks/useModal";
import { conversationStore } from "../../stores/conversationStore";
import AddGroupMemberModal from "./AddGroupMemberModal.vue";
import UserPhoto from "../reusables/UserPhoto.vue";
import LabeledInput from "../ui/LabeledInput.vue";

const props = defineProps({
	conversation: Object,
});

const localConversation = ref({ ...props.conversation });

const $toast = useToast();
const { openModal, closeModal } = useModal();

const userId = computed(() => auth.userId);
const conversationId = computed(() => props.conversation?.conversationId);

// members need to be reactive
const members = ref([...(props.conversation?.members || [])]);

// Keep in sync when store changes
watch(
	() => conversationStore.currentConversation.value?.members,
	(newMembers) => {
		if (newMembers) {
			members.value = [...newMembers];
			// props.conversation.members = [...newMembers]; // keep props in sync too
		}
	},
	{ immediate: true }
);
// const conversationMembers = computed(() => props.conversation?.members);
const conversationMembers = computed(() => members.value);

const selectedUserIds = ref([]);

const isAdmin = computed(() => {
	return conversationMembers.value.find((user) => user.userId == userId.value)?.role === "admin";
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
			$toast.success(`Hai rimosso ${successfullyRemovedIds.length} utenti`);

			conversationsStore.removeMembers(conversationId.value, successfullyRemovedIds);
			conversationStore.removeMembers(conversationId.value, successfullyRemovedIds);
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

// Photo handling
const selectedFile = ref(null);
const previewUrl = ref(null);
const fileInputRef = ref(null);
const groupName = ref(localConversation.value.name || "");

const groupPhotoPath = computed(() => conversationsStore.getPhotoPath(conversationId.value));
const displayedGroupPhoto = computed(() => {
	return previewUrl.value || groupPhotoPath.value;
});

function onFileButtonClick() {
	if (fileInputRef.value) fileInputRef.value.click();
}

function onFileChange(event) {
	const file = event.target.files[0];
	if (!file) return;

	const allowedTypes = ["image/jpeg", "image/png", "image/webp"];
	if (!allowedTypes.includes(file.type)) {
		$toast.error("Il file deve essere jpg, png o webp.");
		return;
	}
	if (file.size > 2000000) {
		$toast.error("Il file può essere al massimo di 2Mb.");
		return;
	}

	selectedFile.value = file;
	previewUrl.value = URL.createObjectURL(file);
}

function removeSelectedFile() {
	if (previewUrl.value) {
		URL.revokeObjectURL(previewUrl.value);
	}
	selectedFile.value = null;
	previewUrl.value = null;
}

async function handleSetGroupPhoto() {
	if (!selectedFile.value) return;
	try {
		const res = await setGroupPhoto(conversationId.value, selectedFile.value);
		if (res.status === 200) {
			$toast.success("Foto del gruppo aggiornata");
			conversationsStore.updatePhotoPath(conversationId.value, res.data.photoPath); // conversationS (store)
			conversationStore.updatePhotoPath(conversationId.value, res.data.photoPath); // conversatioN
			removeSelectedFile();
		} else {
			const err = getError(res, "Errore durante l'aggiornamento della foto");
			$toast.error(err);
		}
	} catch (err) {
		$toast.error("Errore durante il caricamento dell'immagine");
	}
}

async function handleSetGroupName() {
	if (!groupName.value.trim()) return;
	try {
		const res = await setGroupName(conversationId.value, groupName.value);
		if (res.status === 200) {
			localConversation.value.name = groupName.value;
			conversationsStore.updateConversationName(conversationId.value, groupName.value); // conversationS (store)
			conversationStore.updateConversationName(conversationId.value, groupName.value); // conversatioN
			$toast.success("Nome del gruppo aggiornato");
		} else {
			const err = getError(res, "Impossibile modificare il nome del gruppo");
			$toast.error(err);
		}
	} catch (err) {
		$toast.error("Errore durante l'aggiornamento del nome");
	}
}

function handleOpenAddMemberModal() {
	openModal(AddGroupMemberModal, { conversation: props.conversation });
}

const submitGroupNameEnabled = computed(() => {
	return groupName.value.trim() !== "" && groupName.value !== localConversation.value.name;
});
</script>

<template>
	<Modal>
		<div class="groupInformationsModal-container">
			<ModalTitle text="Informazioni del gruppo" />
			<div class="groupInformationsModal-edit">
				<UserPhoto :url="displayedGroupPhoto" :size="100" />
				<div class="groupInformationsModal-editPhoto-buttons">
					<button
						v-if="!selectedFile"
						class="groupInformationsModal-modify-button"
						@click="onFileButtonClick"
					>
						Modifica
					</button>
					<input
						type="file"
						ref="fileInputRef"
						accept="image/png, image/jpeg, image/webp"
						@change="onFileChange"
						style="display: none"
					/>
					<button v-if="selectedFile" class="groupInformationsModal-undo-button" @click="removeSelectedFile">
						Annulla
					</button>
					<button
						v-if="selectedFile"
						class="groupInformationsModal-update-button"
						@click="handleSetGroupPhoto"
					>
						Conferma
					</button>
				</div>

				<div class="groupInformationsModal-editName">
					<span>Nome del gruppo</span>
					<LabeledInput placeholder="Cambia nome del gruppo..." v-model:value="groupName" />
					<button
						class="groupInformationsModal-update-button"
						:disabled="!submitGroupNameEnabled"
						@click="handleSetGroupName"
					>
						Aggiorna
					</button>
				</div>
			</div>

			<div class="groupInformationsModal-currentUsers">
				<p>Utenti nel gruppo</p>

				<div v-if="isAdmin" class="addMembers-container" @click="handleOpenAddMemberModal">
					<img :src="PlusIcon" />
					<span>Aggiungi un membro</span>
				</div>

				<div class="usersSelector-wrapper">
					<UsersSelector
						:users="conversationMembers"
						:selectedUserIds="selectedUserIds"
						:onSelect="onSelect"
						:maxHeight="isAdmin ? 200 : 300"
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

.groupInformationsModal-edit {
	margin-bottom: 20px;
	display: flex;
	gap: 8px;
	align-items: center;
	width: 100%;
}
.groupInformationsModal-editPhoto-buttons {
	display: flex;
	flex-direction: column;
	gap: 8px;
	align-items: center;
}

.groupInformationsModal-editName {
	margin-left: auto;
	display: flex;
	flex-direction: column;
	gap: 8px;
	color: white;

	width: 300px;
}

.groupInformationsModal-editName span {
	margin-left: 10px;
}
.groupInformationsModal-editName button {
	margin-right: auto;
	margin-left: 10px;
	margin-top: 4px;
}

.groupInformationsModal-modify-button,
.groupInformationsModal-undo-button,
.groupInformationsModal-update-button {
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	width: 90px;
}

.groupInformationsModal-undo-button {
	background-color: var(--color-red);
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

.groupInformationsModal-modify-button,
.groupInformationsModal-update-button,
.removeUsers-button {
	background-color: var(--color-green-primary);
}

.leaveGroup-button {
	margin-top: auto;
	margin-left: auto;
	background-color: var(--color-red);
}

.groupInformationsModal-modify-button:hover,
.groupInformationsModal-undo-button:hover,
.groupInformationsModal-update-button:hover,
.addMembers-container:hover,
.removeUsers-button:hover,
.leaveGroup-button:hover {
	opacity: 0.7;
}
.groupInformationsModal-modify-button:disabled,
.groupInformationsModal-undo-button:disabled,
.groupInformationsModal-update-button:disabled,
.removeUsers-button:disabled,
.leaveGroup-button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
