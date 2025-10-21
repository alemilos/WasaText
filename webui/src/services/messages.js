import { Endpoints } from "../config/endpoints";
import { del, post, put } from "./axios";

export function forwardMessage() {
	const { messageId } = boh;
	const url = Endpoints.forwardMessage.replace(":id", messageId);

	return post(url);
}

export function commentMessage() {
	const { messageId } = boh;
	const url = Endpoints.commentMessage.replace(":id", messageId);

	return put(url);
}

export function uncommentMessage() {
	const { messageId } = boh;
	const url = Endpoints.uncommentMessage.replace(":id", messageId);

	return del(url);
}

export function deleteMessage() {
	const { messageId } = boh;
	const url = Endpoints.deleteMessage.replace(":id", messageId);
	return del(url);
}
