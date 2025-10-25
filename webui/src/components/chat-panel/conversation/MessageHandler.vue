<script setup>
import { ref, computed } from "vue";
import GalleryIcon from "@/assets/icons/gallery-icon.svg";
import sendMessageIcon from "@/assets/icons/send-message.svg";
import { useToast } from "vue-toast-notification";
import { sendImageMessage, sendTextMessage } from "../../../services/conversations";
import { conversationStore } from "../../../stores/conversationStore";
import { auth } from "../../../stores/authStore";
import { getError } from "../../../utils/getError";

const $toast = useToast();
const messageText = ref("");
const selectedFile = ref(null);

const currentConversation = computed(() => conversationStore.currentConversation.value);

const conversationId = computed(() => currentConversation.value?.conversationId ?? null);

const userId = computed(() => auth.userId);

const fileInputRef = ref(null);

function onFileButtonClick() {
	if (fileInputRef.value) {
		fileInputRef.value.click(); // trigger file picker
	}
}

function onFileChange(event) {
	const file = event.target.files[0];
	if (!file) return;

	// validate file size and type
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
}

async function sendMessage(conversationId) {
	try {
		let lastMessage = { authorId: userId };
		let res;
		if (selectedFile.value) {
			// send image && text
			res = await sendImageMessage(conversationId, selectedFile.value, messageText.value);
			selectedFile.value = null; // reset file
			messageText.value = ""; // reset text
		} else {
			// send text
			const text = messageText.value.trim();
			if (!text) return;
			res = await sendTextMessage(conversationId, text);
			messageText.value = ""; // reset text
		}

		lastMessage = { ...lastMessage, ...res.data };
		conversationStore.pushMessage(lastMessage);

		if (!res.status === 200) {
			const err = getError(res, "C'è stato un errore nel mandare il messaggio");
			$toast.error(err);
		}
	} catch (err) {
		$toast.error("C'è stato un errore nel mandare il messaggio");
	}
}

function onKeyPress(e) {
	if (e.key === "Enter") sendMessage(conversationId.value);
}
</script>

<template>
	<div class="messageHandler-container">
		<button class="messageHandler-image-button" @click="onFileButtonClick">
			<img :src="GalleryIcon" />
		</button>
		<!-- Hidden input to handle file selection -->
		<input
			type="file"
			ref="fileInputRef"
			accept="image/png, image/jpeg, image/webp"
			@change="onFileChange"
			style="display: none"
		/>

		<div class="messageHandler-input-container">
			<div class="messageHandler-input-wrapper">
				<input
					type="text"
					class="messageHandler-input"
					placeholder="Scrivi un messaggio..."
					v-model="messageText"
					@keypress="onKeyPress"
				/>

				<!-- Container to show that image was selecgted -->
				<div v-if="selectedFile" class="selected-file-overlay">
					<span>Hai allegato un'immagine al messaggio, cancellala se cambi idea.</span>
					<button @click="selectedFile = null" class="remove-file-button">×</button>
				</div>
			</div>
			<img
				:src="sendMessageIcon"
				class="messageHandler-submit-button"
				@click="() => sendMessage(conversationId)"
			/>
		</div>
	</div>
</template>
