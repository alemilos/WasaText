import { Endpoints } from "../config/endpoints";
import { get, post, postForm } from "./axios";

export function createConversation() {
	return post(Endpoints.createConversation);
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
export function sendImageMessage(conversationId, imageFile) {
	const formData = new FormData();
	formData.append("type", "image");
	formData.append("file", imageFile);

	const url = Endpoints.sendMessage.replace(":id", conversationId);
	return postForm(url, formData);
}

export function readMessages() {
	const { conversationId } = boh;
	const url = Endpoints.readMessages.replace(":id", conversationId);

	return post(url);
}
