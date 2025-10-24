import { reactive } from "vue";
import { getMyConversations } from "../services/conversations";

const conversations = reactive({});

async function loadConversations() {
	try {
		const res = await getMyConversations();
		if (res.status === 200) {
			// clear old data
			for (const key in conversations) delete conversations[key];

			// add new conversations
			res.data.conversations.forEach((conv) => {
				conversations[conv.conversationId] = conv;
			});
		}
	} catch (err) {
		console.error("Failed to load conversations:", err);
	}
}

function addConversation(conversation) {
	if (!conversation || !conversation.conversationId) {
		return;
	}
	conversations[conversation.conversationId] = conversation;
}

function removeConversation(conversationId) {
	if (!conversationId) return;
	delete conversations[conversationId];
}

function clearConversations() {
	for (const key in conversations) delete conversations[key];
}

function getPhotoPath(conversationId) {
	const path = conversations[conversationId]?.photoPath ?? null;
	if (!path) return null;
	// Add a version query string to bypass caching
	return `${path}?v=${Date.now()}`;
}

function updatePhotoPath(conversationId, photoPath) {
	if (!conversations[conversationId]) {
		return;
	}
	// Add a version query string to bypass caching
	conversations[conversationId].photoPath = `${photoPath}?v=${Date.now()}`;
}

function updateConversationName(conversationId, name) {
	if (!conversations[conversationId]) {
		return;
	}

	conversations[conversationId].name = name;
}

function removeMembers(conversationId, memberIds) {
	if (!conversationId || !Array.isArray(memberIds) || !memberIds.length) return;
	const conversation = conversations[conversationId];
	if (!conversation || !Array.isArray(conversation.members)) return;

	conversation.members = conversation.members.filter((member) => !memberIds.includes(member.userId));
}

function addMembers(conversationId, memberIds) {
	if (!conversationId || !Array.isArray(memberIds) || !memberIds.length) return;
	const conversation = conversations[conversationId];
	if (!conversation) return;

	if (!Array.isArray(conversation.members)) {
		conversation.members = [];
	}

	const existingIds = conversation.members.map((m) => m.userId);
	const newMembers = memberIds.filter((id) => !existingIds.includes(id)).map((id) => ({ userId: id }));

	if (newMembers.length) {
		conversation.members.push(...newMembers);
	}
}

export const conversationsStore = {
	conversations,
	loadConversations,
	addConversation,
	removeConversation,
	clearConversations,
	getPhotoPath,
	updatePhotoPath,
	updateConversationName,
	removeMembers,
	addMembers,
};
