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

async function pollConversation() {
	if (!currentConversation.value?.conversationId) return;

	try {
		const res = await getConversation(
			currentConversation.value.conversationId
		);
		if (res.status === 200) {
			const newData = res.data;

			// update other fields (title, participants, etc.)
			currentConversation.value = {
				...currentConversation.value,
				...newData,
			};
		}
	} catch (err) {}
}

function pushMessage(message) {
	if (!currentConversation.value) return;

	const oldMessages = currentConversation.value.messages || [];
	currentConversation.value.messages = [...oldMessages, message];
}

export const conversationStore = {
	currentConversation,
	isLoading,
	error,
	openConversation,
	closeConversation,
	pollConversation,
	pushMessage,
};
