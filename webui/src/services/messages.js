import { Endpoints } from "../config/endpoints";
import { del, post, put } from "./axios";

export function forwardMessage(messageId, conversationId) {
	const url = Endpoints.forwardMessage.replace(":id", messageId);

	return post(url, { conversationId });
}

export function commentMessage(messageId, emoji) {
	const url = Endpoints.commentMessage.replace(":id", messageId);

	return put(url, { emoji });
}

export function uncommentMessage(messageId) {
	const url = Endpoints.uncommentMessage.replace(":id", messageId);

	return del(url);
}

export function deleteMessage(messageId) {
	const url = Endpoints.deleteMessage.replace(":id", messageId);
	return del(url);
}
