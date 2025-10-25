<script setup>
import { ref, computed } from "vue";
import UsersSelector from "../side-panel/user-panel/UsersSelector.vue"; // reused for now
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { useModal } from "../../hooks/useModal";
import { useToast } from "vue-toast-notification";
import { conversationsStore } from "../../stores/conversationsStore";
import { usersStore } from "../../stores/usersStore";
import { forwardMessage } from "../../services/messages";
import { getError } from "../../utils/getError";

const props = defineProps({
	message: { type: Object, required: true },
});

const { closeModal } = useModal();
const $toast = useToast();

const conversations = computed(() => Object.values(conversationsStore.conversations));

// Map conversations to userinfo (to use UsersSelector easily)
const mappedConversations = computed(() =>
	conversations.value.map((conv) => {
		const otherUser =
			conv.type === "private"
				? usersStore.getUser(conv.otherParticipantId)
				: { username: conv.name, userId: conv.conversationId, photoPath: conv.photoPath };

		return {
			userId: conv.conversationId,
			username: otherUser?.username,
			photoPath: otherUser?.photoPath,
		};
	})
);

const selectedConversationId = ref(null);

function handleSelectConversation(conversationId) {
	selectedConversationId.value = selectedConversationId.value === conversationId ? null : conversationId;
}

async function handleConfirm() {
	if (!selectedConversationId.value) {
		$toast.error("Seleziona una conversazione prima di continuare");
		return;
	}

	const selectedConv = conversationsStore.conversations[selectedConversationId.value];
	if (!selectedConv) {
		$toast.error("Conversazione non trovata");
		return;
	}

	const res = await forwardMessage(props.message.messageId, selectedConversationId.value);
	if (res.status === 200) {
		$toast.success("Messaggio inoltrato");
	} else {
		const err = getError(res, "Impossibile inoltrare il messaggio");
		$toast.error(err);
	}

	closeModal();
}
</script>

<template>
	<Modal>
		<div class="selectConversationModal-container">
			<ModalTitle text="Inoltra il messaggio" />
			<p v-if="mappedConversations.length">Scegli una conversazione in cui inoltrare il messaggio</p>

			<UsersSelector
				:users="mappedConversations"
				:selected-user-id="selectedConversationId"
				:on-select="handleSelectConversation"
			/>

			<button v-if="mappedConversations.length" @click="handleConfirm" :disabled="!selectedConversationId">
				Conferma
			</button>

			<p v-else class="empty-text">Non hai ancora conversazioni disponibili.</p>
		</div>
	</Modal>
</template>

<style scoped>
.selectConversationModal-container {
	display: flex;
	flex-direction: column;
}

.selectConversationModal-container p {
	color: var(--color-white);
	font-size: 18px;
}

.selectConversationModal-container button {
	background-color: var(--color-green-primary);
	margin-left: auto;
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	margin-top: 16px;
}

.selectConversationModal-container button:hover {
	opacity: 0.7;
}

.selectConversationModal-container button:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.empty-text {
	color: var(--color-white-50);
	font-size: 16px;
	margin-top: 12px;
	text-align: center;
}
</style>
