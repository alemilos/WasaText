<script setup>
import FwIcon from "@/assets/icons/forward.svg?component";
import EmojiIcon from "@/assets/icons/emoji.svg";
import DeleteIcon from "@/assets/icons/delete.svg";
import { useModal } from "../../../../hooks/useModal";
import ForwardMessageModal from "../../../modals/ForwardMessageModal.vue";
import { computed, ref } from "vue";
import { deleteMessage, uncommentMessage } from "../../../../services/messages";
import { useToast } from "vue-toast-notification";
import { getError } from "../../../../utils/getError";
import { conversationStore } from "../../../../stores/conversationStore";
import ReactMessageModal from "../../../modals/ReactMessageModal.vue";
import { auth } from "../../../../stores/authStore";

const { openModal } = useModal();
const $toast = useToast();

const props = defineProps({
	message: { type: Object, required: true },
	isSent: { type: Boolean, required: true },
});

function handleOpenFowardModal() {
	openModal(ForwardMessageModal, { message: props.message });
}

function handleOpenReactModal() {
	openModal(ReactMessageModal, { message: props.message });
}

const userId = computed(() => auth.userId);
const emojiAlreadyReacted = computed(() => props.message.comments?.find((comment) => comment.authorId == userId.value));

async function handleDeleteMessage() {
	const res = await deleteMessage(props.message.messageId);
	if (res.status === 200) {
		$toast.success("Messaggio cancellato");
		conversationStore.removeMessage(props.message.messageId);
	} else {
		const err = getError(res, "Impossibile cancellare il messaggio");
		$toast.error(err);
	}
}

async function handleRemoveReaction() {
	const res = await uncommentMessage(props.message.messageId); // adjust params as needed
	if (res.status !== 200) {
		const err = getError(res, "Impossibile rimuovere la reazione");
		$toast.error(err);
	} else {
		$toast.success("Reazione rimossa");
		hasReacted.value = false;
	}
}

const hasReacted = ref(!!emojiAlreadyReacted.value);
</script>

<template>
	<div class="actionsPopup-container" :class="{ sent: isSent, received: !isSent }">
		<div class="actionsPopup-row" @click="handleOpenFowardModal">
			<img :src="FwIcon" />
			<span>Inoltra</span>
		</div>
		<div class="actionsPopup-row" @click="handleOpenReactModal">
			<img :src="EmojiIcon" />
			<span>Reagisci</span>
		</div>
		<div v-if="hasReacted" class="actionsPopup-row deleteIcon" @click="handleRemoveReaction">
			<img :src="EmojiIcon" />
			<span>Rimuovi Reazione</span>
		</div>
		<div v-if="isSent" class="actionsPopup-row deleteIcon" @click="handleDeleteMessage">
			<img :src="DeleteIcon" />
			<span>Cancella</span>
		</div>
	</div>
</template>

<style scoped>
.actionsPopup-container {
	position: absolute;
	width: 200px;
	height: fit-content;
	border-radius: 8px;
	top: 24px;
	z-index: 10;

	background-color: var(--color-green-dark);
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4), 0 0 8px rgba(0, 0, 0, 0.25);
	display: flex;
	flex-direction: column;
	gap: 4px;
	padding: 8px;
}

/* Sent → align to right */
.actionsPopup-container.sent {
	right: 10px;
}

/* Received → align to left */
.actionsPopup-container.received {
	left: 10px;
}

.actionsPopup-row {
	display: flex;
	gap: 8px;
	cursor: pointer;
	padding: 8px;
	border-radius: 8px;
	align-items: center;
}

.actionsPopup-row:hover {
	background-color: var(--color-white-20);
}

.actionsPopup-row.deleteIcon:hover {
	background-color: var(--color-red-40);
}

.actionsPopup-row img {
	width: 16px;
	height: 16px;
}
</style>
