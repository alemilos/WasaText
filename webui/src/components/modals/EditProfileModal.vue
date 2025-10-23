<script setup>
import { computed, ref } from "vue";
import UserPhoto from "../reusables/UserPhoto.vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { auth } from "../../stores/authStore";
import { usersStore } from "../../stores/usersStore";
import LabeledInput from "../ui/LabeledInput.vue";
import { setMyPhoto, setMyUsername } from "../../services/users";
import { useToast } from "vue-toast-notification";
import { getError } from "../../utils/getError";
import { clearUserId } from "../../services/auth";
import { useRouter } from "vue-router";
import { useModal } from "../../hooks/useModal";

const { closeModal } = useModal();
const $toast = useToast();
const userId = computed(() => auth.userId);
const currentUsername = ref(usersStore.getUsername(userId.value) || "");
// Username handling
const username = ref(currentUsername.value);
const profilePhotoPath = computed(() => usersStore.getPhotoUrl(userId.value));

// Photo handling
const selectedFile = ref(null);
const previewUrl = ref(null);
const fileInputRef = ref(null);

function onFileButtonClick() {
	if (fileInputRef.value) {
		fileInputRef.value.click();
	}
}

function onFileChange(event) {
	const file = event.target.files[0];
	if (!file) return;

	// Validate file type and size
	const allowedTypes = ["image/jpeg", "image/png", "image/webp"];
	if (!allowedTypes.includes(file.type)) {
		$toast.error("Il file deve essere jpg, png o webp.");
		return;
	}
	if (file.size > 2_000_000) {
		$toast.error("Il file può essere al massimo di 2Mb.");
		return;
	}

	selectedFile.value = file;
	previewUrl.value = URL.createObjectURL(file);
}

function removeSelectedFile() {
	selectedFile.value = null;
	if (previewUrl.value) {
		URL.revokeObjectURL(previewUrl.value);
		previewUrl.value = null;
	}
}

async function handleUpdatePhoto() {
	if (!selectedFile.value) return;
	try {
		const res = await setMyPhoto(selectedFile.value);
		if (res.status === 200) {
			$toast.success("Foto profilo aggiornata");
			usersStore.updatePhotoPath(userId.value, res.data.photoPath);
			removeSelectedFile();
		} else {
			const err = getError(
				res,
				"Errore durante l'aggiornamento della foto"
			);
			$toast.error(err);
		}
	} catch (err) {
		$toast.error("Errore durante il caricamento dell'immagine");
	}
}

async function handleUpdateUsername() {
	const res = await setMyUsername(username.value);
	if (res.status === 200) {
		currentUsername.value = username.value;
		usersStore.updateUsername(userId.value, username.value);
		$toast.success("Username modificato");
	} else {
		const err = getError(res, "Impossibile modificare l'username");
		$toast.error(err);
	}
}

const router = useRouter();
async function handleLogout() {
	clearUserId();
	router.push("/login");
	closeModal();
}

const submitUsernameEnabled = computed(() => {
	return (
		username.value.trim() !== "" && username.value !== currentUsername.value
	);
});

const displayedPhoto = computed(() => {
	return previewUrl.value || profilePhotoPath.value;
});
</script>
<template>
	<Modal>
		<div class="editProfileModal-container">
			<ModalTitle text="Modifica profilo" />
			<p>
				Queste sono le tue informazioni profilo, puoi anche modificarle.
			</p>

			<div class="editProfileModal-photoEdit">
				<UserPhoto :url="displayedPhoto" :size="100" />
				<div class="editProfileModal-photoEdit-buttons">
					<button v-if="!selectedFile" @click="onFileButtonClick">
						Modifica
					</button>
					<input
						type="file"
						ref="fileInputRef"
						accept="image/png, image/jpeg, image/webp"
						@change="onFileChange"
						style="display: none"
					/>
					<button
						v-if="selectedFile"
						class="editProfileModal-photoEdit-undo-button"
						@click="removeSelectedFile"
					>
						Annulla
					</button>
					<button v-if="selectedFile" @click="handleUpdatePhoto">
						Conferma
					</button>
				</div>
			</div>

			<div class="editProfileModal-usernameEdit">
				<div class="editProfileModal-inputWrapper">
					<LabeledInput
						label="Username"
						placeholder="Modifica il tuo username..."
						v-model:value="username"
					/>
				</div>
				<button
					:disabled="!submitUsernameEnabled"
					@click="handleUpdateUsername"
				>
					Conferma cambio
				</button>
			</div>

			<button @click="handleLogout" class="editProfileModal-logout">
				Logout
			</button>
		</div>
	</Modal>
</template>

<style scoped>
.editProfileModal-container {
	display: flex;
	flex-direction: column;
	height: 100%;
}
.editProfileModal-container p {
	color: var(--color-white);
	font-size: 18px;
}
.editProfileModal-photoEdit {
	margin-top: 16px;
	width: 100%;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;

	gap: 16px;
}
.editProfileModal-photoEdit-buttons {
	display: flex;
	gap: 8px;
}

.editProfileModal-logout,
.editProfileModal-photoEdit-undo-button {
	background-color: var(--color-red) !important;
}

.editProfileModal-logout,
.editProfileModal-usernameEdit button,
.editProfileModal-photoEdit button {
	background-color: var(--color-green-primary);
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
}

.editProfileModal-logout:hover,
.editProfileModal-photoEdit-undo-button:hover,
.editProfileModal-usernameEdit button:hover,
.editProfileModal-photoEdit button:hover {
	opacity: 0.7;
}

.editProfileModal-logout:disabled,
.editProfileModal-photoEdit-undo-button:disabled,
.editProfileModal-usernameEdit button:disabled,
.editProfileModal-photoEdit button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
.editProfileModal-usernameEdit {
	display: flex;
	margin-top: 30px;
	align-items: center;
	gap: 12px;
}

.editProfileModal-inputWrapper {
	flex: 1;
}

.editProfileModal-logout {
	margin-left: auto;
	padding: 8px 30px;
	margin-top: auto;
}
</style>
