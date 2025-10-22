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

export const conversationsStore = {
	conversations,
	loadConversations,
	addConversation,
	removeConversation,
};
