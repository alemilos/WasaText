import { Endpoints } from "../config/endpoints";
import { get, post } from "./axios";

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

export function sendMessage() {
	const { conversationId } = boh;
	const url = Endpoints.sendMessage.replace(":id", conversationId);

	return post(url);
}

export function readMessages() {
	const { conversationId } = boh;
	const url = Endpoints.readMessages.replace(":id", conversationId);

	return post(url);
}
