import { Endpoints } from "../config/endpoints";
import { get, post, postForm } from "./axios";

export function createConversation(recipientId) {
	return post(Endpoints.createConversation, { recipientId });
}

export function getMyConversations() {
	return get(Endpoints.getMyConversations);
}

export function getConversation(conversationId) {
	const url = Endpoints.getConversation.replace(":id", conversationId);
	return get(url);
}

export function sendTextMessage(conversationId, content) {
	const url = Endpoints.sendMessage.replace(":id", conversationId);
	return post(url, { type: "text", content });
}
export function sendImageMessage(conversationId, imageFile, text) {
	const formData = new FormData();
	formData.append("type", "image");
	formData.append("file", imageFile);
	formData.append("secondaryContent", text);

	const url = Endpoints.sendMessage.replace(":id", conversationId);
	return postForm(url, formData);
}

export function readMessages(conversationId, messages) {
	if (!messages || !messages.length) return;
	const url = Endpoints.readMessages.replace(":id", conversationId);

	return post(url, { messages });
}
