import { endpoints } from "../config/endpoints";
import { get, post } from "./axios";

export function createConversation() {
	return post(endpoints.createConversation);
}

export function getMyConversations() {
	return get(endpoints.getMyConversations);
}

export function getConversation() {
	const { conversationId } = boh;
	const url = endpoints.getConversation.replace(":id", conversationId);

	return get(url);
}

export function sendMessage() {
	const { conversationId } = boh;
	const url = endpoints.sendMessage.replace(":id", conversationId);

	return post(url);
}

export function readMessages() {
	const { conversationId } = boh;
	const url = endpoints.readMessages.replace(":id", conversationId);

	return post(url);
}
