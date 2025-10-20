import { endpoints } from "../config/endpoints";
import { del, post, put } from "./axios";

export function forwardMessage() {
	const { messageId } = boh;
	const url = endpoints.forwardMessage.replace(":id", messageId);

	return post(url);
}

export function commentMessage() {
	const { messageId } = boh;
	const url = endpoints.commentMessage.replace(":id", messageId);

	return put(url);
}

export function uncommentMessage() {
	const { messageId } = boh;
	const url = endpoints.uncommentMessage.replace(":id", messageId);

	return del(url);
}

export function deleteMessage() {
	const { messageId } = boh;
	const url = endpoints.deleteMessage.replace(":id", messageId);
	return del(url);
}
