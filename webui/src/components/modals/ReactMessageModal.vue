<script setup>
import { ref, computed } from "vue";
import Modal from "../ui/Modal.vue";
import ModalTitle from "../ui/ModalTitle.vue";
import { commentMessage, uncommentMessage } from "../../services/messages";
import { getError } from "../../utils/getError";
import { useToast } from "vue-toast-notification";
import { auth } from "../../stores/authStore";
import { conversationStore } from "../../stores/conversationStore";

const props = defineProps({
	message: { type: Object, required: true },
});

const userId = computed(() => auth.userId);
const emojiAlreadyReacted = computed(() => props.message.comments?.find((comment) => comment.authorId == userId.value));
const conversationId = computed(() => conversationStore.currentConversation?.conversationId);

const $toast = useToast();

const emit = defineEmits(["select"]);

const EMOJIS = [
	"😀",
	"😁",
	"😂",
	"🤣",
	"😃",
	"😄",
	"😅",
	"😆",
	"😉",
	"😊",
	"🙂",
	"🙃",
	"😇",
	"😍",
	"🥰",
	"😘",
	"😗",
	"😚",
	"😙",
	"😋",
	"😛",
	"😜",
	"🤪",
	"😝",
	"🤗",
	"🤭",
	"🤫",
	"🤔",
	"🤨",
	"😐",
	"😑",
	"😶",
	"🙄",
	"😏",
	"😣",
	"😥",
	"😮",
	"😯",
	"😲",
	"😳",
	"🥺",
	"😢",
	"😭",
	"😤",
	"😠",
	"😡",
	"🤬",
	"🤯",
	"😱",
	"😨",
	"😰",
	"😓",
	"🥵",
	"🥶",
	"😴",
	"😪",
	"😵",
	"🤐",
	"🥴",
	"🤢",
	"🤮",
	"🤧",
	"😷",
	"🤒",
	"🤕",
	"🤑",
	"🤠",
	"😎",
	"🤡",
	"🥳",
	"🤩",
	"😇",
	"👍",
	"👎",
	"👊",
	"🤛",
	"🤜",
	"🤞",
	"🤟",
	"🤘",
	"👋",
	"👏",
	"🙌",
	"🙏",
	"💪",
	"🫶",
	"🤝",
	"👌",
	"❤️",
	"💔",
	"💖",
	"💗",
	"💓",
	"💕",
	"💞",
	"💘",
	"💝",
	"💛",
	"💚",
	"💙",
	"💜",
	"🖤",
	"🤍",
	"🤎",
	"🔥",
	"✨",
	"🎉",
	"🎊",
	"🎈",
	"🌈",
	"🌸",
	"🌹",
	"🌺",
	"🌻",
	"🌷",
	"☀️",
	"🌤️",
	"🌧️",
	"❄️",
	"🌙",
	"⭐️",
	"⚡️",
	"☁️",
	"🌊",
	"🍀",
	"🍕",
	"🍔",
	"🍟",
];

const selectedEmoji = ref(emojiAlreadyReacted.value?.emoji || null);

function handleSelect(emoji) {
	if (selectedEmoji.value == emoji) {
		selectedEmoji.value = null;
	} else {
		selectedEmoji.value = emoji;
	}
}

async function handleConfirm() {
	if (!selectedEmoji.value) return;

	const res = await commentMessage(props.message.messageId, selectedEmoji.value);
	if (res.status !== 200) {
		const err = getError(res, "Impossibile reagire al messaggio");
		$toast.error(err);
	} else {
		// Close on successs
		$toast.success("Reazione aggiunta");
		const comment = { emoji: selectedEmoji.value, authorId: userId.value };
		conversationStore.commentMessage(conversationId.value, props.message.messageId, comment);
		hasReacted.value = true;
	}
}

async function handleRemoveReaction() {
	const res = await uncommentMessage(props.message.messageId);
	if (res.status !== 200) {
		const err = getError(res, "Impossibile rimuovere la reazione");
		$toast.error(err);
	} else {
		$toast.success("Reazione rimossa");
		selectedEmoji.value = null;
		hasReacted.value = false;
	}
}

const hasReacted = ref(!!emojiAlreadyReacted.value);
</script>

<template>
	<Modal>
		<div class="emojiSelectorModal-container">
			<ModalTitle text="Reagisci al messaggio" />
			<p>Seleziona un'emoji dalla griglia qui sotto</p>

			<div class="emoji-grid">
				<button
					v-for="(emoji, idx) in EMOJIS"
					:key="idx"
					class="emoji-btn"
					@click="handleSelect(emoji)"
					:class="{ selected: selectedEmoji === emoji }"
				>
					{{ emoji }}
				</button>
			</div>

			<div class="actions-row">
				<button @click="handleConfirm" :disabled="!selectedEmoji" class="confirm-btn">Conferma</button>

				<button v-if="hasReacted" @click="handleRemoveReaction" class="remove-btn">Rimuovi reazione</button>
			</div>
		</div>
	</Modal>
</template>

<style scoped>
.emojiSelectorModal-container {
	display: flex;
	flex-direction: column;
}

.emojiSelectorModal-container p {
	color: var(--color-white);
	font-size: 18px;
	margin-bottom: 12px;
}

.emoji-grid {
	display: grid;
	grid-template-columns: repeat(8, 1fr);
	gap: 8px;
	max-height: 400px;
	overflow-y: auto;
	padding: 8px;
	margin-bottom: 16px;
}

.emoji-grid::-webkit-scrollbar {
	width: 8px;
}
.emoji-grid::-webkit-scrollbar-thumb {
	background: rgba(255, 255, 255, 0.1);
	border-radius: 6px;
}

.emoji-btn {
	background: transparent;
	border: none;
	font-size: 24px;
	cursor: pointer;
	padding: 8px;
	border-radius: 8px;
	transition: background 0.15s ease;
}
.emoji-btn:hover {
	background: rgba(255, 255, 255, 0.1);
}
.emoji-btn.selected {
	background: rgba(255, 255, 255, 0.2);
}

.actions-row {
	display: flex;
	gap: 12px;
	margin-top: 12px;
}

.confirm-btn {
	background-color: var(--color-green-primary);
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	cursor: pointer;
}
.confirm-btn:hover {
	opacity: 0.7;
}
.confirm-btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.remove-btn {
	background-color: #ff3b3b;
	color: white;
	border: none;
	border-radius: 8px;
	padding: 6px 16px;
	cursor: pointer;
}
.remove-btn:hover {
	opacity: 0.85;
}
</style>
