import { ref } from "vue";
import { getConversation } from "../services/conversations";

const currentConversation = ref(null); // no conversation opened by default
const isLoading = ref(false);
const error = ref(null);

async function openConversation(conversationId) {
	if (!conversationId) return;

	isLoading.value = true;
	error.value = null;

	try {
		const res = await getConversation(conversationId); // your API call
		if (res.status === 200) {
			currentConversation.value = res.data;
		} else {
			error.value = `Failed to fetch conversation ${conversationId}`;
			currentConversation.value = null;
		}
	} catch (err) {
		error.value = err.message || "Unknown error";
		currentConversation.value = null;
	} finally {
		isLoading.value = false;
	}
}

function closeConversation() {
	currentConversation.value = null;
	error.value = null;
}

function pushMessage(message) {
	if (!currentConversation.value) return;

	// Ensure messages array exists
	if (!currentConversation.value.messages) {
		currentConversation.value.messages = [];
	}

	// Push new message and trigger reactivity
	currentConversation.value.messages.push(message);
}

export const conversationStore = {
	currentConversation,
	isLoading,
	error,
	openConversation,
	closeConversation,
	pushMessage,
};
