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
		const res = await getConversation(currentConversation.value.conversationId);
		if (res.status === 200) {
			const newData = res.data;

			// update other fields (title, participants, etc.)
			currentConversation.value = {
				...currentConversation.value,
				...newData,
			};
		} else {
			return { shouldClose: true }; // notify caller that he must close the conversation
		}
	} catch (err) {}
}

function pushMessage(message) {
	if (!currentConversation.value) return;

	const oldMessages = currentConversation.value.messages || [];
	currentConversation.value.messages = [...oldMessages, message];
}

function removeMessage(messageId) {
	if (!currentConversation.value) return;

	currentConversation.value.messages = currentConversation.value.messages.filter(
		(message) => message.messageId !== messageId
	);
}

function updatePhotoPath(conversationId, photoPath) {
	if (!currentConversation.value) return;
	if (!currentConversation.value.conversationId == conversationId) return;
	// Add a version query string to bypass caching
	currentConversation.value.photoPath = `${photoPath}?v=${Date.now()}`;
}

function updateConversationName(conversationId, name) {
	if (!currentConversation.value) return;
	if (!currentConversation.value.conversationId == conversationId) return;

	currentConversation.value.photoPath = name;
}

function removeMembers(conversationId, memberIds) {
	if (!currentConversation.value || currentConversation.value.conversationId !== conversationId) return;

	if (!Array.isArray(currentConversation.value.members)) return;

	currentConversation.value.members = currentConversation.value.members.filter(
		(member) => !memberIds.includes(member.userId)
	);
}

function addMembers(conversationId, memberIds) {
	if (!currentConversation.value || currentConversation.value.conversationId !== conversationId) return;

	if (!Array.isArray(memberIds) || !memberIds.length) return;

	if (!Array.isArray(currentConversation.value.members)) {
		currentConversation.value.members = [];
	}

	const existingIds = currentConversation.value.members.map((m) => m.userId);
	const newMembers = memberIds.filter((id) => !existingIds.includes(id)).map((id) => ({ userId: id }));

	if (newMembers.length) {
		currentConversation.value.members.push(...newMembers);
	}
}

function commentMessage(conversationId, messageId, comment) {
	if (!currentConversation.value || currentConversation.value.conversationId !== conversationId) return;

	const message = currentConversation.value.messages?.find((m) => m.messageId === messageId);
	if (!message) return;

	if (!Array.isArray(message.comments)) message.comments = [];
	message.comments.push(comment);
}

export const conversationStore = {
	currentConversation,
	isLoading,
	error,
	openConversation,
	closeConversation,
	pollConversation,
	pushMessage,
	removeMessage,
	updatePhotoPath,
	updateConversationName,
	removeMembers,
	addMembers,
	commentMessage,
};
